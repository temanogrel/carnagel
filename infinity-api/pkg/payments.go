package infinity

import (
	"github.com/blockcypher/gobcy"
	"time"

	"context"

	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

const (
	ConsulKeyGuestPlanUuid = "infinity/payment-plans/guest-plan-uuid"
	ConsulKeyBasicPlanUuid = "infinity/payment-plans/basic-plan-uuid"
)

const SatoshisPerBitCoin = 100000000
const UnconfirmedTransactionPaymentPlanDuration = time.Hour * 24

type PaymentTransactionState uint8

const (
	_ = PaymentTransactionState(iota)
	PaymentTransactionStatePending
	PaymentTransactionStatePartiallyPaid
	PaymentTransactionStateFullyPaid
	PaymentTransactionStateTooMuchPaid
	PaymentTransactionStatePaidWithinAllowedDiff
	PaymentTransactionStateExpired
)

var (
	PaymentPlayNotFoundErr           = errors.New("The plan could not be found")
	UserNotOnPaymentPlanErr          = errors.New("The current user is not on a active payment plan")
	NoMatchingPaymentTransactionsErr = errors.New("There are no payment transactions matching the criteria")
	PaymentTransactionNotFoundErr    = errors.New("The payment transaction could not be found")
	NoPaymentReceivedErr             = errors.New("No payment has been received on the given address")
)

type PaymentPlan struct {
	TableName             struct{}  `sql:"payment_plans, alias:pp" json:"-"`
	Uuid                  uuid.UUID `sql:",pk" json:"uuid"`
	Name                  string    `json:"name"`
	Description           string    `json:"description"`
	Bandwidth             uint64    `json:"bandwidth"`
	PerRecordingBandwidth uint64    `json:"perRecordingBandwidth"`
	Devices               uint8     `json:"devices"`
	Price                 float32   `json:"price"`
	Duration              uint      `json:"duration"`
	UpdatedAt             time.Time `json:"updatedAt"`
	CreatedAt             time.Time `json:"createdAt"`
}

type PaymentPlanCollection map[uuid.UUID]*PaymentPlan

type PaymentPlanRepository interface {
	GetByUuid(id uuid.UUID) (*PaymentPlan, error)
	GetGuestPlan() (*PaymentPlan, error)
	GetBasicPlan() (*PaymentPlan, error)
	GetByUserUuid(id uuid.UUID) (*PaymentPlan, error)
	GetAll() ([]*PaymentPlan, error)
}

type PaymentService interface {
	Run(ctx context.Context)
	InitiatePurchase(user *User, plan *PaymentPlan) (*PaymentTransaction, error)
	UpdatePurchase(transaction *PaymentTransaction) error
	AddressForwardingCallback(transaction *PaymentTransaction, addrForward *gobcy.Payback) error
}

type PaymentTransaction struct {
	TableName                           struct{}     `sql:"payment_transactions, alias:pt" json:"-"`
	Uuid                                uuid.UUID    `sql:",pk" json:"uuid"`
	WebhookUuid                         uuid.UUID    `json:"-"`
	PaymentAddress                      string       `json:"paymentAddress"`
	State                               uint8        `json:"state"`
	ConversionRate                      float64      `json:"conversionRate"`
	UnconfirmedReceivedAmountInSatoshis uint64       `json:"unconfirmedReceivedAmountInSatoshis"`
	ReceivedAmountInSatoshis            uint64       `json:"receivedAmountInSatoshis"`
	ExpectedAmountInSatoshis            uint64       `json:"expectedAmountInSatoshis"`
	UpdatedAt                           time.Time    `json:"updatedAt"`
	CreatedAt                           time.Time    `json:"createdAt"`
	PaymentPlan                         *PaymentPlan `json:"-"`
	PaymentPlanUuid                     uuid.UUID    `json:"paymentPlanUuid"`
	User                                *User        `json:"-"`
	UserUuid                            uuid.UUID    `json:"userUuid"`

	// Blockcypher related
	BlockcypherWebhookId           string `json:"-"`
	BlockcypherAddressForwardingId string `json:"-"`
}

func (t *PaymentTransaction) UpdateStateFromAmountReceived() {
	if !t.FullyPaid() {
		t.State = uint8(PaymentTransactionStatePartiallyPaid)
		return
	}

	if t.GetTotalAmountReceivedInSatoshis() > t.ExpectedAmountInSatoshis {
		t.State = uint8(PaymentTransactionStateTooMuchPaid)
	} else if t.GetTotalAmountReceivedInSatoshis() == t.ExpectedAmountInSatoshis {
		t.State = uint8(PaymentTransactionStateFullyPaid)
	} else {
		t.State = uint8(PaymentTransactionStatePaidWithinAllowedDiff)
	}
}

func (t *PaymentTransaction) MinAllowedPayment() uint64 {
	return uint64(float64(t.ExpectedAmountInSatoshis) * 0.95)
}

func (t *PaymentTransaction) FullyPaid() bool {
	// Within the allowed diff
	return t.GetTotalAmountReceivedInSatoshis() >= t.MinAllowedPayment()
}

func (t *PaymentTransaction) CanDeleteBlockcypherWebhook() bool {
	return t.ConfirmedFullyPaid() && t.BlockcypherWebhookId != ""
}

func (t *PaymentTransaction) CanDeleteBlockCypherAddressForwarding() bool {
	return t.ConfirmedFullyPaid() && t.BlockcypherAddressForwardingId != ""
}

func (t *PaymentTransaction) ShouldSkipUpdate(addr gobcy.Addr) bool {
	// Check if transaction is expired
	if t.State == uint8(PaymentTransactionStateExpired) {
		return true
	}

	if t.ConfirmedFullyPaid() {
		return true
	}

	return t.ReceivedAmountInSatoshis >= uint64(addr.Balance) && t.UnconfirmedReceivedAmountInSatoshis >= uint64(addr.UnconfirmedBalance)
}

func (t *PaymentTransaction) ConfirmedFullyPaid() bool {
	// Within the allowed diff
	return t.ReceivedAmountInSatoshis >= t.MinAllowedPayment()
}

func (t *PaymentTransaction) GetTotalAmountReceivedInSatoshis() uint64 {
	return t.ReceivedAmountInSatoshis + t.UnconfirmedReceivedAmountInSatoshis
}

type PaymentTransactionRepositoryCriteria struct {
	Offset int
	Limit  int

	CreatedBefore time.Time
	CreatedAfter  time.Time

	User uuid.UUID

	States []PaymentTransactionState

	Sorting map[string]string
}

func NewPaymentTransactionRepositoryCriteria(offset, limit int) *PaymentTransactionRepositoryCriteria {
	return &PaymentTransactionRepositoryCriteria{
		Offset: offset,
		Limit:  limit,
	}
}

type PaymentTransactionRepository interface {
	Create(paymentTransaction *PaymentTransaction) error
	Update(paymentTransaction *PaymentTransaction) error
	GetByWebhookUuid(id uuid.UUID) (*PaymentTransaction, error)
	GetByUuid(id uuid.UUID) (*PaymentTransaction, error)
	Matching(criteria *PaymentTransactionRepositoryCriteria) ([]*PaymentTransaction, int, error)
}
