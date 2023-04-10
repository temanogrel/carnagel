package user

import (
	"time"

	"fmt"

	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type repository struct {
	app *infinity.Application
}

func NewRepository(app *infinity.Application) infinity.UserRepository {
	return &repository{app}
}

func (r *repository) GetByUuid(id uuid.UUID) (*infinity.User, error) {

	user := &infinity.User{}

	err := r.app.DB.Model(user).Where("uuid = ?", id).Select()
	switch err {
	case nil:
		return user, nil

	case pg.ErrNoRows:
		return user, infinity.UserNotFoundErr

	default:
		return nil, errors.Wrap(err, "Unexpected database error")
	}
}

func (r *repository) GetByEmail(email string) (*infinity.User, error) {

	user := &infinity.User{}

	err := r.app.DB.Model(user).Where("lower(email) = lower(?)", email).Select()
	switch err {
	case nil:
		return user, nil

	case pg.ErrNoRows:
		return user, infinity.UserNotFoundErr

	default:
		return nil, errors.Wrap(err, "Unexpected database error")
	}
}

func (r *repository) GetByUsername(username string) (*infinity.User, error) {

	user := &infinity.User{}

	err := r.app.DB.Model(user).Where("lower(username) = lower(?)", username).Select()
	switch err {
	case nil:
		return user, nil

	case pg.ErrNoRows:
		return user, infinity.UserNotFoundErr

	default:
		return nil, errors.Wrap(err, "Unexpected database error")
	}
}

func (r *repository) GetByUsernameOrEmail(identity string) (*infinity.User, error) {

	user := &infinity.User{}

	err := r.app.DB.Model(user).
		Where("LOWER(username) = LOWER(?)", identity).
		WhereOr("LOWER(email) = LOWER(?)", identity).
		Select()

	switch err {
	case nil:
		return user, nil

	case pg.ErrNoRows:
		return user, infinity.UserNotFoundErr

	default:
		return nil, errors.Wrap(err, "Unexpected database error")
	}
}

func (r *repository) GetUsersWithExpiredPaymentPlan() ([]*infinity.User, error) {

	users := []*infinity.User{}

	err := r.app.DB.Model(&users).
		Where("payment_plan_ends_at IS NOT NULL").
		Where("payment_plan_ends_at < ?", time.Now()).
		Select()

	if err != nil && err != pg.ErrNoRows {
		return nil, errors.Wrapf(err, "Unexpected database error")
	}

	return users, nil
}

func (r *repository) GetUsersWithExpiringPaymentPlan(daysBeforeExpiration uint) ([]*infinity.User, error) {
	users := []*infinity.User{}

	err := r.app.DB.Model(&users).
		Where("payment_plan_ends_at IS NOT NULL").
		Where(fmt.Sprintf("CAST(payment_plan_ends_at - '%d days'::INTERVAL AS DATE) = CURRENT_DATE", daysBeforeExpiration)).
		Select()

	if err != nil && err != pg.ErrNoRows {
		return nil, errors.Wrapf(err, "Unexpected database error")
	}

	return users, nil
}

func (r *repository) RemoveById(id uuid.UUID) (bool, error) {

	result, err := r.app.DB.
		Model(&infinity.User{}).
		Where("uuid = ?", id.String()).
		Delete()

	return result.RowsAffected() == 1, err
}

func (r *repository) Create(user *infinity.User) error {

	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// default the role in-case it's missing
	if user.Role == "" {
		user.Role = infinity.RoleUser
	}

	if user.PaymentPlanUuid == uuid.Nil {
		plan, err := r.app.PaymentPlanRepository.GetBasicPlan()
		if err != nil {
			return errors.Wrapf(err, "User was missing a payment plan, and retrieving basic plan failed")
		}

		now := time.Now()

		user.PaymentPlanUuid = plan.Uuid
		user.PaymentPlanSubscribedAt = &now
	}

	user.Uuid = uuid.NewV4()
	user.Password = string(password)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err = r.app.DB.Model(user).Insert()

	return err
}

func (r *repository) Update(user *infinity.User) error {
	user.UpdatedAt = time.Now()

	return r.app.DB.Update(user)
}

func (r *repository) ToggleLike(user uuid.UUID, recording uuid.UUID) (bool, error) {

	result, err := r.app.DB.Model(new(infinity.UserLike)).
		Where("user_uuid = ? and recording_uuid = ?", user.String(), recording.String()).
		Delete()

	if err != nil {
		return false, err
	}

	if result.RowsAffected() == 0 {
		_, err := r.app.DB.Model(&infinity.UserLike{
			RecordingUuid: recording,
			UserUuid:      user,
			CreatedAt:     time.Now(),
		}).Insert()

		r.app.DB.Model(new(infinity.Recording)).
			Set("like_count = like_count + 1").
			Where("uuid = ?", recording.String()).
			Update()

		return true, err
	}

	_, err = r.app.DB.Model(new(infinity.Recording)).
		Set("like_count = like_count - 1").
		Where("uuid = ?", recording.String()).
		Update()

	return false, err
}

func (r *repository) ToggleFavorite(user uuid.UUID, recording uuid.UUID) (bool, error) {

	result, err := r.app.DB.Model(new(infinity.UserFavorite)).
		Where("user_uuid = ? and recording_uuid = ?", user.String(), recording.String()).
		Delete()

	if err != nil {
		return false, err
	}

	if result.RowsAffected() == 0 {
		_, err := r.app.DB.Model(&infinity.UserFavorite{
			RecordingUuid: recording,
			UserUuid:      user,
		}).Insert()

		return true, err
	}

	return false, nil
}

func (r *repository) Matching(criteria *infinity.UserRepositoryCriteria) ([]infinity.User, int, error) {
	users := make([]infinity.User, 0)

	query := r.app.DB.Model(&users)
	query.Apply(r.injectCriteria(criteria))

	numRows, err := query.SelectAndCount()

	return users, numRows, err
}

func (r *repository) injectCriteria(criteria *infinity.UserRepositoryCriteria) func(*orm.Query) (*orm.Query, error) {
	return func(q *orm.Query) (*orm.Query, error) {
		if criteria.Limit != 0 {
			q.Limit(criteria.Limit)
		}

		if criteria.Offset != 0 {
			q.Offset(criteria.Offset)
		}

		if !criteria.CreatedAfter.IsZero() {
			q.Where("u.created_at > ?", criteria.CreatedAfter)
		}

		basicPlan, err := r.app.PaymentPlanRepository.GetBasicPlan()
		if err != nil {
			return nil, err
		}

		if criteria.CurrentlyPremium {
			q.
				Where("u.payment_plan_uuid != ?", basicPlan.Uuid).
				Where("u.payment_plan_ends_at > NOW()")
		}

		if criteria.ExPremium {
			q.
				Where("u.payment_plan_uuid = ?", basicPlan.Uuid).
				Where(`(
					SELECT COUNT(*) FROM payment_transactions t
					WHERE
						t.user_uuid = u.uuid AND
						t.confirmed_fully_paid = TRUE
				) > 0`)
		}

		if criteria.NeverPremium {
			q.
				Where("u.payment_plan_uuid = ?", basicPlan.Uuid).
				Where(`(
					SELECT COUNT(*) FROM payment_transactions t
					WHERE
						t.user_uuid = u.uuid AND
						t.confirmed_fully_paid = TRUE
				) = 0`)
		}

		for field, sorting := range criteria.Sorting {
			if sorting != "asc" && sorting != "desc" {
				continue
			}

			// Not safe for sql injection, do not allow user input directly here
			q.Order(fmt.Sprintf("%s %s", field, sorting))
		}

		return q, nil
	}
}
