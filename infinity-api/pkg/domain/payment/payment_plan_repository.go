package payment

import (
	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type planRepository struct {
	app *infinity.Application
	log logrus.FieldLogger
}

func NewPaymentPlanRepository(app *infinity.Application) infinity.PaymentPlanRepository {
	return &planRepository{
		app: app,
		log: app.Logger.WithField("component", "payment_plan"),
	}
}

func (repository *planRepository) GetGuestPlan() (*infinity.PaymentPlan, error) {
	value := repository.app.Consul.GetString(infinity.ConsulKeyGuestPlanUuid, "")

	if value == "" {
		return nil, errors.Errorf("Key %s is missing in consul", infinity.ConsulKeyGuestPlanUuid)
	}

	id, err := uuid.FromString(value)
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid uuid '%s' provided from consul", value)
	}

	return repository.GetByUuid(id)
}

func (repository *planRepository) GetBasicPlan() (*infinity.PaymentPlan, error) {
	value := repository.app.Consul.GetString(infinity.ConsulKeyBasicPlanUuid, "")

	if value == "" {
		return nil, errors.Errorf("Key %s is missing in consul", infinity.ConsulKeyBasicPlanUuid)
	}

	id, err := uuid.FromString(value)
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid uuid '%s' provided from consul", value)
	}

	return repository.GetByUuid(id)
}

func (repository *planRepository) GetByUuid(id uuid.UUID) (*infinity.PaymentPlan, error) {
	plan := &infinity.PaymentPlan{}

	err := repository.app.DB.Model(plan).Where("uuid = ?", id.String()).Select()
	if err == pg.ErrNoRows {
		return plan, infinity.PaymentPlayNotFoundErr
	}

	return plan, err
}

func (repository *planRepository) GetByUserUuid(id uuid.UUID) (*infinity.PaymentPlan, error) {

	user := &infinity.User{}

	err := repository.app.DB.Model(user).
		Column("PaymentPlan").
		Where("uuid = ?", id.String()).
		Select()

	if err == pg.ErrNoRows {
		return nil, infinity.PaymentPlayNotFoundErr
	}

	return user.PaymentPlan, err
}

func (repository *planRepository) GetAll() ([]*infinity.PaymentPlan, error) {

	plans := []*infinity.PaymentPlan{}

	err := repository.app.DB.Model(&plans).Order("price ASC").Select()
	if err != nil && err != pg.ErrNoRows {
		return plans, errors.Wrap(err, "Failed to retrieve the payment plans")
	}

	return plans, err
}
