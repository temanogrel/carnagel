package payment

import (
	"time"

	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type paymentTransactionRepository struct {
	app *infinity.Application
	log logrus.FieldLogger
}

func NewPaymentTransactionRepository(app *infinity.Application) infinity.PaymentTransactionRepository {
	return &paymentTransactionRepository{
		app: app,
		log: app.Logger.WithField("component", "payment_transactions"),
	}
}

func (repository *paymentTransactionRepository) Create(paymentTransaction *infinity.PaymentTransaction) error {
	paymentTransaction.Uuid = uuid.NewV4()
	paymentTransaction.State = uint8(infinity.PaymentTransactionStatePending)
	paymentTransaction.CreatedAt = time.Now()
	paymentTransaction.UpdatedAt = time.Now()

	_, err := repository.app.DB.Model(paymentTransaction).Insert()

	return err
}

func (repository *paymentTransactionRepository) Update(paymentTransaction *infinity.PaymentTransaction) error {
	paymentTransaction.UpdatedAt = time.Now()

	return repository.app.DB.Update(paymentTransaction)
}

func (repository *paymentTransactionRepository) GetByUuid(id uuid.UUID) (*infinity.PaymentTransaction, error) {
	transaction := &infinity.PaymentTransaction{}

	err := repository.app.DB.Model(transaction).Where("uuid = ?", id.String()).Select()
	if err == pg.ErrNoRows {
		return transaction, infinity.PaymentTransactionNotFoundErr
	}

	return transaction, err
}

func (repository *paymentTransactionRepository) GetByWebhookUuid(id uuid.UUID) (*infinity.PaymentTransaction, error) {
	transaction := &infinity.PaymentTransaction{}

	err := repository.app.DB.Model(transaction).Where("webhook_uuid = ?", id.String()).Select()
	if err == pg.ErrNoRows {
		return transaction, infinity.PaymentTransactionNotFoundErr
	}

	return transaction, err
}

func (repository *paymentTransactionRepository) Matching(criteria *infinity.PaymentTransactionRepositoryCriteria) ([]*infinity.PaymentTransaction, int, error) {
	var transactions []*infinity.PaymentTransaction

	count, err := repository.app.DB.Model(&transactions).Apply(repository.injectCriteria(criteria)).SelectAndCount()
	if err != nil {
		return nil, count, err
	}

	return transactions, count, err
}

func (repository *paymentTransactionRepository) injectCriteria(criteria *infinity.PaymentTransactionRepositoryCriteria) func(*orm.Query) (*orm.Query, error) {
	return func(q *orm.Query) (*orm.Query, error) {
		if criteria.Limit != 0 {
			q.Limit(criteria.Limit)
		}

		if criteria.Offset != 0 {
			q.Offset(criteria.Offset)
		}

		if !criteria.CreatedBefore.IsZero() {
			q.Where("created_at <= ?", criteria.CreatedBefore)
		}

		if !criteria.CreatedAfter.IsZero() {
			q.Where("created_at >= ?", criteria.CreatedAfter)
		}

		if criteria.User != uuid.Nil {
			q.Where("user_uuid = ?", criteria.User)
		}

		if len(criteria.States) > 0 {
			q.Where("state IN(?)", pg.In(criteria.States))
		}

		q.Order("created_at ASC")

		return q, nil
	}
}
