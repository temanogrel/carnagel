package payment

import (
	"github.com/blockcypher/gobcy"
	"github.com/sasha-s/go-deadlock"
	"math"
	"time"

	"context"

	"fmt"

	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/infrastructure/http/rpc"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type service struct {
	app    *infinity.Application
	logger logrus.FieldLogger
	mtx    deadlock.RWMutex
}

func NewService(app *infinity.Application) infinity.PaymentService {
	service := &service{
		app:    app,
		logger: app.Logger.WithField("service", "PaymentService"),
		mtx:    deadlock.RWMutex{},
	}

	return service
}

func (service *service) Run(ctx context.Context) {
	logger := service.logger.WithField("operation", "Run")

	// Midnight in UTC we send the payment plan expiry reminder
	midnight := time.Now().
		Truncate(time.Hour * 24).
		Add(time.Hour * 24).
		UTC()

	downgradeExpiredPlans := time.NewTimer(0)
	sendExpiryReminder := time.NewTimer(time.Until(midnight))
	checkPendingPayments := time.NewTimer(time.Hour * 6)
	markTransactionsAsExpired := time.NewTimer(time.Hour)

	for {
		select {
		case <-ctx.Done():
			downgradeExpiredPlans.Stop()
			sendExpiryReminder.Stop()
			checkPendingPayments.Stop()

			return

		case <-markTransactionsAsExpired.C:
			service.markTransactionsAsExpired()
			markTransactionsAsExpired.Reset(time.Hour * 24)

		case <-checkPendingPayments.C:
			// Disable this for now since payments seems to be pending
			service.checkPendingPayments()
			checkPendingPayments.Reset(time.Hour * 6)

		case <-downgradeExpiredPlans.C:
			if err := service.downgradeUsersWithExpiredPaymentPlan(); err != nil {
				logger.WithError(err).Error("Failed to downgrade expired payment plans")
			}

			downgradeExpiredPlans.Reset(time.Hour)

		case <-sendExpiryReminder.C:
			if err := service.sendEmailsToUsersWithExpiringPaymentPlan(); err != nil {
				logger.WithError(err).Errorf("Failed to send payment plan expiry email")
			}

			sendExpiryReminder.Reset(time.Hour * 24)
		}
	}
}

func (service *service) markTransactionsAsExpired() {
	logger := service.logger.WithField("operation", "markTransactionsAsExpired")

	criteria := infinity.NewPaymentTransactionRepositoryCriteria(0, 0)
	// Expire the transaction after a day
	criteria.CreatedBefore = time.Now().Add(-time.Hour * 24 * 1)
	criteria.States = []infinity.PaymentTransactionState{
		infinity.PaymentTransactionStatePending,
		infinity.PaymentTransactionStatePartiallyPaid,
	}

	transactions, _, err := service.app.PaymentTransactionRepository.Matching(criteria)
	if err != nil {
		logger.WithError(err).Error("Failed to retrieve transactions to mark as expired")
		return
	}

	hooks, err := service.app.BlockcypherClient.ListHooks()
	if err != nil {
		logger.WithError(err).Error("Failed to retrieve webhooks")
		return
	}

	forwards, err := service.app.BlockcypherClient.ListPayFwds()
	if err != nil {
		logger.WithError(err).Error("Failed to retrieve address forwards")
		return
	}

	addrToHookIdMap := map[string]string{}
	for _, hook := range hooks {
		addrToHookIdMap[hook.Address] = hook.ID
	}

	addrToForwardIdMap := map[string]string{}
	for _, forward := range forwards {
		addrToForwardIdMap[forward.InputAddr] = forward.ID
	}

	// Go through each transaction
	for _, transaction := range transactions {
		logger = logger.WithField("transactionUuid", transaction.Uuid)

		if webhookId, ok := addrToHookIdMap[transaction.PaymentAddress]; ok {
			// To avoid 429 errors
			time.Sleep(time.Second * 2)

			if err := service.app.BlockcypherClient.DeleteHook(webhookId); err != nil {
				logger.WithError(err).Error("Failed to delete blockcypher webhook")
				continue
			}
		}

		if forwardingId, ok := addrToForwardIdMap[transaction.PaymentAddress]; ok {
			// To avoid 429 errors
			time.Sleep(time.Second * 2)

			if err := service.app.BlockcypherClient.DeletePayFwd(forwardingId); err != nil {
				logger.WithError(err).Error("Failed to delete blockcypher addfress forwarding")
				continue
			}
		}

		transaction.State = uint8(infinity.PaymentTransactionStateExpired)
		if err := service.app.PaymentTransactionRepository.Update(transaction); err != nil {
			logger.WithError(err).Error("Failed to update transaction state to expired")
		}
	}
}

func (service *service) checkPendingPayments() {
	logger := service.logger.WithField("operation", "checkPendingPayments")

	criteria := infinity.NewPaymentTransactionRepositoryCriteria(0, 0)
	// Don't use a too large diff since all viewing of payments methods generate unpaid addresses
	criteria.CreatedAfter = time.Now().Add(-time.Hour * 18)
	criteria.States = []infinity.PaymentTransactionState{infinity.PaymentTransactionStatePending}

	transactions, _, err := service.app.PaymentTransactionRepository.Matching(criteria)
	if err != nil {
		logger.WithError(err).Error("Error retrieving pending payments")

		return
	}

	// Go through each transaction page
	for _, transaction := range transactions {
		logger = logger.WithField("transactionUuid", transaction.Uuid)

		addr, err := service.app.BlockcypherClient.GetAddrBal(transaction.PaymentAddress, nil)
		if err != nil {
			logger.WithError(err).Errorf("Failed to retrieve balance of payment address")
			continue
		}

		// Since we don't reuse addresses and we use address forwarding the balance might have been moved
		// So recalculate the balance using the amount received and sent
		if addr.Balance == 0 && addr.UnconfirmedBalance == 0 {
			addr.Balance = addr.TotalReceived
			addr.UnconfirmedBalance = 0
		}

		logger = logger.WithField("bitcoinAddressData", addr)

		if err := service.processUpdatedPayment(transaction.Uuid, addr); err != nil {
			logger.WithError(err).Errorf("Failed to process pending payment")
		}

		// Sleep some between each request since there's also a rate limit per second
		time.Sleep(time.Second * 2)
	}
}

func (service *service) InitiatePurchase(user *infinity.User, plan *infinity.PaymentPlan) (*infinity.PaymentTransaction, error) {
	webhookUuid := uuid.NewV4()

	forwardingAddr := gobcy.PayFwd{
		Destination: infinity.BlockcypherForwardingDestination,
	}

	addr, err := service.app.BlockcypherClient.CreatePayFwd(forwardingAddr)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to set up address forwarding")
	}

	amountInBtc, rate, err := service.app.CryptoExchangeRateService.ConvertUSDToBtc(float64(plan.Price))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get current btc/usd exchange rate")
	}

	hook := gobcy.Hook{
		Event:   "tx-confirmation",
		Address: addr.InputAddr,
		URL:     fmt.Sprintf("https://api.camtube.co/blockcypher/webhook/%s", webhookUuid.String()),
	}

	if hook, err = service.app.BlockcypherClient.CreateHook(hook); err != nil {
		return nil, errors.Wrap(err, "Failed to register blockcypher webhook")
	}

	transaction := &infinity.PaymentTransaction{
		WebhookUuid:                    webhookUuid,
		User:                           user,
		UserUuid:                       user.Uuid,
		PaymentPlan:                    plan,
		PaymentPlanUuid:                plan.Uuid,
		PaymentAddress:                 addr.InputAddr,
		ExpectedAmountInSatoshis:       uint64(math.Ceil(amountInBtc * infinity.SatoshisPerBitCoin)),
		ConversionRate:                 rate,
		BlockcypherWebhookId:           hook.ID,
		BlockcypherAddressForwardingId: addr.ID,
	}

	if err = service.app.PaymentTransactionRepository.Create(transaction); err != nil {
		return nil, errors.Wrap(err, "Failed to create payment transaction")
	}

	return transaction, nil
}

func (service *service) UpdatePurchase(transaction *infinity.PaymentTransaction) error {
	addr, err := service.app.BlockcypherClient.GetAddrBal(transaction.PaymentAddress, nil)
	if err != nil {
		return errors.Wrap(err, "Failed to retrieve balance of payment address")
	}

	// Since we don't reuse addresses and we use address forwarding the balance might have been moved
	// So recalculate the balance using the amount received and sent
	if addr.Balance == 0 && addr.UnconfirmedBalance == 0 {
		addr.Balance = addr.TotalReceived
		addr.UnconfirmedBalance = 0
	}

	if err := service.processUpdatedPayment(transaction.Uuid, addr); err != nil {
		return err
	}

	return nil
}

func (service *service) AddressForwardingCallback(transaction *infinity.PaymentTransaction, addrForward *gobcy.Payback) error {
	// Just map it straight to the correct type
	//return service.processUpdatedPayment(transaction.Uuid, gobcy.Addr{UnconfirmedBalance: addrForward.Value})
	return nil
}

func (service *service) processUpdatedPayment(transactionId uuid.UUID, addr gobcy.Addr) error {
	logger := service.app.Logger.WithField("operation", "processUpdatedPayment")

	// Make sure we can only process one at a time
	service.mtx.Lock()
	defer service.mtx.Unlock()

	// Load updated because hooks can fire multiple times in an extremely short amount of time
	transaction, err := service.app.PaymentTransactionRepository.GetByUuid(transactionId)
	if err != nil {
		return errors.Wrapf(err, "Failed to retrieve updated transaction")
	}

	if transaction.ShouldSkipUpdate(addr) {
		return nil
	}

	currentAmountReceived := transaction.GetTotalAmountReceivedInSatoshis()
	isCurrentlyFullyPaid := transaction.FullyPaid()

	// Make sure we have a balance change
	if addr.Balance == 0 && addr.UnconfirmedBalance == 0 {
		return infinity.NoPaymentReceivedErr
	}

	logger.WithFields(logrus.Fields{
		"currentBalance":            transaction.ReceivedAmountInSatoshis,
		"currentUnconfirmedBalance": transaction.UnconfirmedReceivedAmountInSatoshis,
		"currentState":              transaction.State,
		"expectedAmount":            transaction.ExpectedAmountInSatoshis,
		"updatedBalance":            uint64(addr.Balance),
		"updatedUnconfirmedBalance": uint64(addr.UnconfirmedBalance),
	}).Debug("Updating payment transaction")

	transaction.ReceivedAmountInSatoshis = uint64(addr.Balance)
	transaction.UnconfirmedReceivedAmountInSatoshis = uint64(addr.UnconfirmedBalance)
	transaction.UpdateStateFromAmountReceived()

	if transaction.FullyPaid() {
		if err := service.upgradeUserPaymentPlan(transaction, !isCurrentlyFullyPaid); err != nil {
			return errors.Wrap(err, "Failed to upgrade user payment plan")
		}
	}

	if err := service.app.PaymentTransactionRepository.Update(transaction); err != nil {
		return errors.Wrap(err, "Failed to update transaction")
	}

	// Tell frontend to poll new transaction
	if currentAmountReceived != transaction.GetTotalAmountReceivedInSatoshis() {
		payload := struct {
			TransactionUuid uuid.UUID `json:"transactionUuid"`
		}{TransactionUuid: transaction.Uuid}

		service.app.ActiveSessionCollection.SendToUser(transaction.UserUuid, rpc.CreateBroadcastResponse("purchase:update", payload))
	}

	// Clean up webhook if transaction is fully paid
	if transaction.CanDeleteBlockcypherWebhook() {
		service.app.BlockcypherClient.DeleteHook(transaction.BlockcypherWebhookId)
	}

	if transaction.CanDeleteBlockCypherAddressForwarding() {
		service.app.BlockcypherClient.DeletePayFwd(transaction.BlockcypherAddressForwardingId)
	}

	return nil
}

func (service *service) downgradeUsersWithExpiredPaymentPlan() error {
	users, err := service.app.UserRepository.GetUsersWithExpiredPaymentPlan()
	if err != nil {
		return err
	}

	basicPlan, err := service.app.PaymentPlanRepository.GetBasicPlan()
	if err != nil {
		return errors.Wrap(err, "Failed to retrieve basic payment plan")
	}

	for _, user := range users {
		now := time.Now()

		user.PaymentPlan = basicPlan
		user.PaymentPlanUuid = basicPlan.Uuid
		user.PaymentPlanEndsAt = nil
		user.PaymentPlanSubscribedAt = &now

		if err = service.app.UserRepository.Update(user); err != nil {
			return errors.Wrap(err, "Failed to set guest plan for user")
		}

		service.app.EmailService.SendPaymentPlanExpiredEmail(user)
	}

	return nil
}

func (service *service) sendEmailsToUsersWithExpiringPaymentPlan() error {
	// Send email to users with one week before expiration
	users, err := service.app.UserRepository.GetUsersWithExpiringPaymentPlan(7)
	if err != nil {
		return err
	}

	for _, user := range users {
		service.app.EmailService.SendPaymentPlanExpiryReminderEmail(user, 7)
	}

	// Send email to users with one day before expiration
	users, err = service.app.UserRepository.GetUsersWithExpiringPaymentPlan(1)
	if err != nil {
		return err
	}

	for _, user := range users {
		service.app.EmailService.SendPaymentPlanExpiryReminderEmail(user, 7)
	}

	return nil
}

func (service *service) upgradeUserPaymentPlan(transaction *infinity.PaymentTransaction, sendEmail bool) error {
	user, err := service.app.UserRepository.GetByUuid(transaction.UserUuid)
	if err != nil {
		return errors.Wrap(err, "Failed to retrieve transaction owner")
	}

	plan, err := service.app.PaymentPlanRepository.GetByUuid(transaction.PaymentPlanUuid)
	if err != nil {
		return errors.Wrap(err, "Failed to retrieve payment plan associated with transaction")
	}

	start := time.Now()
	if user.PaymentPlanEndsAt != nil && user.PaymentPlanEndsAt.After(time.Now()) {
		start = *user.PaymentPlanEndsAt
	}

	var endsAt time.Time

	if transaction.ConfirmedFullyPaid() {
		endsAt = start.Add(time.Hour * 24 * time.Duration(plan.Duration))
	} else {
		endsAt = start.Add(infinity.UnconfirmedTransactionPaymentPlanDuration)
	}

	now := time.Now()

	user.PaymentPlan = plan
	user.PaymentPlanUuid = plan.Uuid
	user.PaymentPlanEndsAt = &endsAt
	user.PaymentPlanSubscribedAt = &now

	if err = service.app.UserRepository.Update(user); err != nil {
		return errors.Wrap(err, "Failed to persist the transactions owner's new payment plan")
	}

	if sendEmail {
		service.app.EmailService.SendPurchasedRegularPremiumEmail(user, plan, transaction)
	}

	return nil
}
