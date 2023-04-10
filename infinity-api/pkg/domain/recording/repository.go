package recording

import (
	"fmt"
	"github.com/sasha-s/go-deadlock"
	"sync"
	"time"

	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

type repository struct {
	app *infinity.Application
	mtx deadlock.RWMutex

	// in-memory map of external id to uuid
	externalIdMap map[uint64]uuid.UUID
}

func NewWrappedRepository(app *infinity.Application) *repository {
	return &repository{
		app: app,
		mtx: deadlock.RWMutex{},

		externalIdMap: make(map[uint64]uuid.UUID),
	}
}

func (repo *repository) MarkAsValidated(ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}

	recording := &infinity.Recording{}
	now := time.Now()
	recording.LastValidationAt = &now

	_, err := repo.app.DB.Model(recording).
		Where("uuid IN (?)", pg.In(ids)).
		Column("last_validation_at").
		Update()

	return err
}

func (repo *repository) GetByUuid(id uuid.UUID) (*infinity.Recording, error) {
	recording := &infinity.Recording{}

	err := repo.app.DB.Model(recording).
		Column("recording.*", "Performer").
		Where("recording.uuid = ?", id.String()).
		Select()

	switch err {
	case nil:
		return recording, nil

	case pg.ErrNoRows:
		return recording, infinity.RecordingNotFoundErr

	default:
		repo.app.Logger.
			WithError(err).
			Error("Database error in RecordingRepository.GetByUuid")

		return nil, errors.Wrap(err, "Unexpected database error")
	}
}

func (repo *repository) GetUuidByExternalId(externalId uint64) (uuid.UUID, error) {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	id, ok := repo.externalIdMap[externalId]
	if ok {
		return id, nil
	}

	err := repo.app.DB.
		Model(&infinity.Recording{}).
		Column("uuid").
		Where("external_id = ?", externalId).
		Select(&id)

	switch err {
	case nil:
		// cache for future reference
		repo.externalIdMap[externalId] = id

		return id, nil

	case pg.ErrNoRows:
		return id, infinity.RecordingNotFoundErr

	default:
		return id, errors.Wrapf(err, "Unexpected database error")
	}
}

func (repo *repository) GetByUuidWithUserContext(recordingId uuid.UUID, userId uuid.UUID) (*infinity.RecordingWithUserData, error) {
	recording := new(infinity.RecordingWithUserData)

	err := repo.app.DB.Model(recording).
		Column("recording.*", "Performer").
		Where("recording.uuid = ?", recordingId.String()).
		Apply(repo.injectUserContext(userId)).
		Select()

	switch err {
	case nil:
		return recording, nil

	case pg.ErrNoRows:
		return recording, infinity.RecordingNotFoundErr

	default:
		repo.app.Logger.
			WithError(err).
			Error("Database error in RecordingRepository.GetByUuid")

		return nil, errors.Wrap(err, "Unexpected database error")
	}
}

func (repo *repository) AddView(recording *infinity.Recording, userHash string) (bool, error) {
	err := repo.app.DB.Model(&infinity.RecordingView{}).
		Where(`
			recording_uuid = ? AND
			user_hash = ? AND
			created_at > NOW() - '1 hours'::INTERVAL
		`, recording.Uuid.String(), userHash).
		Select()

	if err == pg.ErrNoRows {
		err := repo.app.DB.Insert(&infinity.RecordingView{
			RecordingUuid: recording.Uuid,
			UserHash:      userHash,
			CreatedAt:     time.Now(),
		})

		_, err = repo.app.DB.Model(recording).
			Set("view_count = view_count + 1").
			Where("uuid = ?", recording.Uuid.String()).
			Update()

		return true, err
	}

	return false, err
}

func (repo *repository) Create(recording *infinity.Recording, performer *infinity.Performer) error {
	return repo.app.DB.RunInTransaction(func(tx *pg.Tx) error {
		recording.Uuid = uuid.NewV4()
		recording.PerformerUuid = performer.Uuid
		recording.CreatedAt = time.Now()
		recording.UpdatedAt = time.Now()
		recording.ViewCount = 0
		recording.LikeCount = 0

		if err := tx.Insert(recording); err != nil {
			return err
		}

		sql := "UPDATE performers SET recording_count = recording_count + 1 WHERE uuid = ?"
		if _, err := tx.Exec(sql, performer.Uuid); err != nil {
			return err
		}

		// Update after performing the SQL query
		performer.RecordingCount += 1

		return nil
	})
}

func (repo *repository) GetByExternalId(id uint64) (*infinity.Recording, error) {

	recording := &infinity.Recording{}

	err := repo.app.DB.Model(recording).Where("external_id = ?", id).Select()
	switch err {
	case nil:
		return recording, nil

	case pg.ErrNoRows:
		return recording, infinity.RecordingNotFoundErr

	default:
		repo.app.Logger.
			WithError(err).
			Error("Database error in RecordingRepository.GetByExternalId")

		return nil, errors.Wrap(err, "Unexpected database error")
	}
}

func (repo *repository) GetBySlug(slug string) (*infinity.Recording, error) {
	recording := &infinity.Recording{}

	err := repo.app.DB.Model(recording).
		Column("recording.*", "Performer").
		Where("recording.slug = ?", slug).
		Select()
	switch err {
	case nil:
		return recording, nil

	case pg.ErrNoRows:
		return recording, infinity.RecordingNotFoundErr

	default:
		repo.app.Logger.
			WithError(err).
			Error("Database error in RecordingRepository.GetBySlug")

		return nil, errors.Wrap(err, "Unexpected database error")
	}
}

func (repo *repository) GetBySlugWithUserContext(slug string, userId uuid.UUID) (*infinity.RecordingWithUserData, error) {
	recording := new(infinity.RecordingWithUserData)

	err := repo.app.DB.Model(recording).
		Column("recording.*", "Performer").
		Where("recording.slug = ?", slug).
		Apply(repo.injectUserContext(userId)).
		Select()

	switch err {
	case nil:
		return recording, nil

	case pg.ErrNoRows:
		return recording, infinity.RecordingNotFoundErr

	default:
		repo.app.Logger.
			WithError(err).
			Error("Database error in RecordingRepository.GetBySlugWithUserContext")

		return nil, errors.Wrap(err, "Unexpected database error")
	}
}

func (repo *repository) GetAllMissingSlug() ([]infinity.Recording, error) {
	var recordings []infinity.Recording

	err := repo.app.DB.Model(&recordings).Where("slug IS NULL").Select()
	switch err {
	case nil:
		return recordings, nil

	case pg.ErrNoRows:
		return recordings, nil

	default:
		repo.app.Logger.
			WithError(err).
			Error("Database error in RecordingRepository.GetAllMissingSlug")

		return nil, errors.Wrap(err, "Unexpected database error")
	}
}

func (repo *repository) CanRemove(id uuid.UUID) error {
	// Check if the recording has views
	var views []infinity.RecordingView

	if err := repo.app.DB.Model(&views).Where("recording_uuid = ?", id).Select(); err != nil && err != pg.ErrNoRows {
		return err
	}

	if len(views) > 0 {
		return infinity.RecordingHasViewsErr
	}

	basicPlan, err := repo.app.PaymentPlanRepository.GetBasicPlan()
	if err != nil {
		return err
	}

	// Check if the recording is part of a premium user favorite collection
	var favorites []infinity.UserFavorite

	err = repo.app.DB.Model(&favorites).
		Where("recording_uuid = ?", id).
		Where(`user_uuid IN (
			SELECT u.uuid FROM users u
			WHERE 
				u.payment_plan_uuid != ? AND
				u.payment_plan_ends_at > NOW()
		)`, basicPlan.Uuid).
		Select()

	if err != nil && err != pg.ErrNoRows {
		return err
	}

	if len(favorites) > 0 {
		return infinity.RecordingInPremiumUserCollectionErr
	}

	return nil
}

func (repo *repository) Remove(id uuid.UUID) (bool, error) {
	var deleted bool

	err := repo.app.DB.RunInTransaction(func(tx *pg.Tx) error {
		recording := &infinity.Recording{}

		if err := tx.Model(recording).Where("uuid = ?", id.String()).Select(); err != nil {
			if err == pg.ErrNoRows {
				return nil
			}

			return errors.Wrapf(err, "Error retrieving recording by uuid")
		}

		result, err := tx.Model(&infinity.Recording{}).Where("uuid = ?", id.String()).Delete()
		if err != nil {
			repo.app.Logger.
				WithError(err).
				Error("Database error in RecordingRepository.Remove")

			return errors.Wrap(err, "Unexpected database error")
		}

		sql := "UPDATE performers SET recording_count = recording_count - 1 WHERE uuid = ?"
		if _, err := tx.Exec(sql, recording.PerformerUuid); err != nil {
			return errors.Wrapf(err, "Error updating recording count for performer")
		}

		deleted = result.RowsAffected() >= 1

		return nil
	})

	return deleted, err
}

func (repo *repository) Update(recording *infinity.Recording) error {
	_, err := repo.app.DB.Model(recording).Update()
	if err != nil {
		repo.app.Logger.
			WithError(err).
			Error("Database error in RecordingRepository.Update")

		return errors.Wrap(err, "Unexpected database error")
	}

	return nil
}

func (repo *repository) GetRecordingsCreatedBefore(time time.Time, limit int) ([]infinity.Recording, error) {
	recordings := make([]infinity.Recording, 0)

	err := repo.app.DB.Model(&recordings).
		Column("recording.*", "Performer").
		Where("recording.created_at <= ?", time).
		Limit(limit).
		Select()

	return recordings, err
}

func (repo *repository) Matching(criteria *infinity.RecordingRepositoryCriteria, includeCount bool) ([]infinity.Recording, int, error) {

	var recordings []infinity.Recording

	dataQuery := repo.app.DB.Model(&recordings).Column("recording.*", "Performer")
	dataQuery.Apply(repo.injectCriteria(criteria, uuid.Nil))

	// If we're ignoring count
	if !includeCount {
		if e := dataQuery.Select(); e != nil {
			return recordings, 0, e
		}

		return recordings, 0, nil
	}

	countQuery := repo.app.DB.Model(&infinity.Recording{})
	countQuery.Apply(repo.injectCriteria(criteria, uuid.Nil))

	var err error
	var count int

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		if e := dataQuery.Select(); e != nil {
			err = e
		}
	}()

	go func() {
		defer wg.Done()

		var e error
		count, e = countQuery.Count()
		if e != nil {
			err = e
		}
	}()

	wg.Wait()

	return recordings, count, err
}

func (repo *repository) MatchingWithUserContext(criteria *infinity.RecordingRepositoryCriteria, user uuid.UUID, includeCount bool) ([]infinity.RecordingWithUserData, int, error) {
	var recordings []infinity.RecordingWithUserData

	dataQuery := repo.app.DB.Model(&recordings).Column("recording.*", "Performer")
	dataQuery.Apply(repo.injectUserContext(user))
	dataQuery.Apply(repo.injectCriteria(criteria, user))

	// If we're ignoring count
	if !includeCount {
		if e := dataQuery.Select(); e != nil {
			return recordings, 0, e
		}

		return recordings, 0, nil
	}

	countQuery := repo.app.DB.Model(&infinity.Recording{})
	countQuery.Apply(repo.injectCriteria(criteria, user))
	countQuery.Offset(0)

	var err error
	var count int

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		if e := dataQuery.Select(); e != nil {
			err = e
		}
	}()

	go func() {
		defer wg.Done()

		var e error
		count, e = countQuery.Count()
		if e != nil {
			err = e
		}
	}()

	wg.Wait()

	return recordings, count, err
}

func (repo *repository) injectCriteria(criteria *infinity.RecordingRepositoryCriteria, user uuid.UUID) func(*orm.Query) (*orm.Query, error) {
	return func(q *orm.Query) (*orm.Query, error) {
		if criteria.Limit != 0 {
			q.Limit(criteria.Limit)
		}

		if criteria.Offset != 0 {
			q.Offset(criteria.Offset)
		}

		if criteria.PerformerId != uuid.Nil {
			q.Where("performer_uuid = ?", criteria.PerformerId)
		}

		if !criteria.CreatedAfter.IsZero() {
			q.Where("recording.created_at >= ?", criteria.CreatedAfter)
		}

		if !criteria.CreatedBefore.IsZero() {
			q.Where("recording.created_at <= ?", criteria.CreatedBefore)
		}

		if !criteria.LastValidationBefore.IsZero() {
			q.WhereGroup(func(q *orm.Query) (*orm.Query, error) {
				return q.
					WhereOr("recording.last_validation_at IS NULL").
					WhereOr("recording.last_validation_at <= ?", criteria.LastValidationBefore), nil
			})
		}

		if criteria.FavoritesOnly {
			q.Join(`
				INNER JOIN user_recording_favorites urf
				ON
					urf.user_uuid = ? AND
					urf.recording_uuid = recording.uuid
			`, user.String())
		}

		if criteria.OriginService != "" {
			q.Where("performer.origin_service = ?", criteria.OriginService)
		}

		if criteria.OriginSection != "" {
			q.Where("performer.origin_section = ?", criteria.OriginSection)
		}

		if criteria.ExternalIdMin != 0 {
			q.Where("recording.external_id >= ?", criteria.ExternalIdMin)
		}

		if criteria.ExternalIdMax != 0 {
			q.Where("recording.external_id <= ?", criteria.ExternalIdMax)
		}

		if criteria.StageName != "" {
			// Similarity is basically useless to apply unless there's 4 characters
			// since it will return too low similarity for short strings
			if len(criteria.StageName) < 5 {
				regex := fmt.Sprintf("%s%s%s", "%", criteria.StageName, "%")
				q.Where("LOWER(recording.stage_name) LIKE LOWER(?)", regex)
			} else {
				q.Where("recording.stage_name % ?", criteria.StageName)
				q.OrderExpr("SIMILARITY(recording.stage_name, ?) DESC", criteria.StageName)
			}
		}

		switch criteria.SortMode {
		case infinity.SortModeExternalIdAsc:
			q.Order("external_id ASC")

		case infinity.SortModePopularity:
			q.Where(fmt.Sprintf("recording.updated_at >= NOW() - '%d days'::INTERVAL", criteria.Interval))
			q.Where(fmt.Sprintf(`recording.uuid IN (
					SELECT url3.recording_uuid
					FROM user_recording_likes url3
					WHERE url3.created_at >= NOW() - '%d days'::INTERVAL
					ORDER BY (
					SELECT COUNT(*) AS like_count
						FROM user_recording_likes url4
						WHERE
							url3.recording_uuid = url4.recording_uuid AND
							url4.created_at >= NOW() - '%d days'::INTERVAL
					) DESC
				)`, criteria.Interval, criteria.Interval))
			q.Order("like_count DESC")

		case infinity.SortModeViews:
			q.Where(fmt.Sprintf("recording.updated_at >= NOW() - '%d days'::INTERVAL", criteria.Interval))
			q.Where(fmt.Sprintf(`recording.uuid IN (
					SELECT rv.recording_uuid
					FROM recording_views rv
					WHERE rv.created_at >= NOW() - '%d days'::INTERVAL
					ORDER BY (
					SELECT COUNT(*) AS view_count
						FROM recording_views rv2
						WHERE
							rv.recording_uuid = rv2.recording_uuid AND
							rv2.created_at >= NOW() - '%d days'::INTERVAL
					) DESC
				)`, criteria.Interval, criteria.Interval))
			q.Order("view_count DESC")

		case infinity.SortModeOldest:
			q.Order("created_at ASC")

		default:
			q.Order("created_at DESC")
		}

		return q, nil
	}
}

func (repo *repository) injectUserContext(user uuid.UUID) func(*orm.Query) (*orm.Query, error) {
	return func(q *orm.Query) (*orm.Query, error) {
		q.
			ColumnExpr(`(
			SELECT COUNT(urf2.recording_uuid)
			FROM user_recording_favorites AS urf2
			WHERE
				urf2.user_uuid = ? AND
				urf2.recording_uuid = recording.uuid
			) AS is_favorite`, user.String()).
			ColumnExpr(`(
			SELECT COUNT(url2.recording_uuid)
			FROM user_recording_likes as url2
			WHERE
				url2.user_uuid = ? AND
				url2.recording_uuid = recording.uuid
			) as is_liked`, user.String())

		return q, nil
	}
}
