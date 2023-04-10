package performer

import (
	"time"

	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

type repository struct {
	app *infinity.Application
}

func NewRepository(app *infinity.Application) infinity.PerformerRepository {
	return &repository{app: app}
}

func (repo *repository) GetByUuid(id uuid.UUID) (*infinity.Performer, error) {
	performer := new(infinity.Performer)

	err := repo.app.DB.Model(performer).Where("uuid = ?", id.String()).Select()
	if err == pg.ErrNoRows {
		return performer, infinity.PerformerNotFoundErr
	}

	return performer, err
}

func (repo *repository) GetByExternalId(id uint64) (*infinity.Performer, error) {
	performer := new(infinity.Performer)

	err := repo.app.DB.Model(performer).Where("external_id = ?", id).Select()
	if err == pg.ErrNoRows {
		return performer, infinity.PerformerNotFoundErr
	}

	return performer, err
}

func (repo *repository) GetBySlug(slug string) (*infinity.Performer, error) {
	performer := &infinity.Performer{}

	err := repo.app.DB.Model(performer).Where("slug = ?", slug).Select()
	switch err {
	case nil:
		return performer, nil

	case pg.ErrNoRows:
		return performer, infinity.PerformerNotFoundErr

	default:
		repo.app.Logger.
			WithError(err).
			Error("Database error in RecordingRepository.GetBySlug")

		return nil, errors.Wrap(err, "Unexpected database error")
	}
}

func (repo *repository) GetAllMissingSlug() ([]infinity.Performer, error) {
	var performers []infinity.Performer

	err := repo.app.DB.Model(&performers).Where("slug IS NULL").Select()
	switch err {
	case nil:
		return performers, nil

	case pg.ErrNoRows:
		return performers, nil

	default:
		repo.app.Logger.
			WithError(err).
			Error("Database error in PerformerRepository.GetAllMissingSlug")

		return nil, errors.Wrap(err, "Unexpected database error")
	}
}

func (repo *repository) Matching(criteria *infinity.PerformerRepositoryCriteria) ([]*infinity.Performer, int, error) {

	var performers []*infinity.Performer

	query := repo.app.DB.Model(&performers)
	query.Apply(repo.injectConstraints(criteria))

	numRows, err := query.SelectAndCount()

	return performers, numRows, err
}

func (repo *repository) Create(performer *infinity.Performer) error {

	performer.Uuid = uuid.NewV4()
	performer.CreatedAt = time.Now()
	performer.UpdatedAt = time.Now()

	return repo.app.DB.Insert(performer)
}

func (repo *repository) Update(performer *infinity.Performer) error {
	performer.UpdatedAt = time.Now()

	return repo.app.DB.Update(performer)
}

func (repo *repository) injectConstraints(criteria *infinity.PerformerRepositoryCriteria) func(*orm.Query) (*orm.Query, error) {
	return func(q *orm.Query) (*orm.Query, error) {

		if criteria.Limit != 0 {
			q.Limit(criteria.Limit)
		}

		if criteria.Offset != 0 {
			q.Offset(criteria.Offset)
		}

		if len(criteria.PerformerIds) > 0 {
			q.Where("performer.uuid IN (?)", pg.In(criteria.PerformerIds))
		}

		if criteria.IncludeLatestRecording {
			q.Column("performer.*", "Recordings").Relation("Recordings", func(q *orm.Query) (*orm.Query, error) {
				return q.Order("created_at DESC"), nil
			})
		}

		return q, nil
	}
}
