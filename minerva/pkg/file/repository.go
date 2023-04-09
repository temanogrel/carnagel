package file

import (
	"time"

	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"git.misc.vee.bz/carnagel/minerva/pkg"
	"git.misc.vee.bz/carnagel/minerva/pkg/internal/database"
	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

type repository struct {
	app *minerva.Application
}

func NewFileRepository(application *minerva.Application) minerva.FileRepository {
	return &repository{app: application}
}

func (repo *repository) GetByExternalIds(ids []ecosystem.ExternalId) ([]*minerva.File, error) {
	files := make([]*minerva.File, 0)

	err := repo.app.Database.Model(&files).Where("external_id IN(?)", pg.In(ids)).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return files, nil
		}

		return nil, err
	}

	return files, nil
}

func (repo *repository) GetWithPendingOperations(hostname minerva.Hostname, limit uint64) ([]*minerva.File, error) {
	files := []*minerva.File{}

	err := repo.app.Database.Model(&files).Where(`	hostname = ? AND (pending_upload = true OR pending_deletion = true)`, string(hostname)).Select()
	switch err {
	case nil, pg.ErrNoRows:
		return files, nil

	default:
		return files, errors.Wrap(err, database.UnexpectedDatabaseErr)
	}
}

func (repo *repository) Delete(file *minerva.File) error {
	if _, err := repo.app.Database.Model(file).Delete(); err != nil {
		return errors.Wrap(err, database.UnexpectedDatabaseErr)
	}

	return nil
}

func (repo *repository) TrackHit(id uuid.UUID, remoteAddr string) error {
	err := repo.app.Database.Insert(&minerva.FileHit{
		Uuid:          uuid.NewV4(),
		FileUuid:      id,
		RemoteAddress: remoteAddr,
		CreatedAt:     time.Now(),
	})

	if err != nil {
		return errors.Wrap(err, database.UnexpectedDatabaseErr)
	}

	return nil
}

func (repo *repository) GetByUuid(id uuid.UUID) (*minerva.File, error) {
	file := &minerva.File{}

	err := repo.app.Database.Model(file).Where("uuid = ?", id.String()).Select()
	switch err {
	case nil:
		return file, nil

	case pg.ErrNoRows:
		return nil, minerva.FileNotFoundErr

	default:
		repo.app.Logger.WithError(err).Error("Unexpected database error occurred, and both return values tend to be nil")

		return nil, errors.Wrap(err, database.UnexpectedDatabaseErr)
	}
}

func (repo *repository) GetByLocation(hostname, path string) (*minerva.File, error) {
	file := &minerva.File{}

	err := repo.app.Database.Model(file).Where("hostname = ? AND path = ?", hostname, path).Select()
	switch err {
	case nil:
		return file, nil

	case pg.ErrNoRows:
		return nil, minerva.FileNotFoundErr

	default:
		repo.app.Logger.WithError(err).Error("Unexpected database error occurred, and both return values tend to be nil")

		return nil, errors.Wrap(err, database.UnexpectedDatabaseErr)
	}
}

func (repo *repository) GetByExternalId(id ecosystem.ExternalId) ([]*minerva.File, error) {
	files := []*minerva.File{}

	err := repo.app.Database.Model(&files).Where("external_id = ?", id).Select()
	switch err {
	case nil, pg.ErrNoRows:
		return files, nil

	default:
		return nil, errors.Wrap(err, database.UnexpectedDatabaseErr)
	}
}

func (repo repository) Create(file *minerva.File) error {
	if file.Uuid == uuid.Nil {
		file.Uuid = uuid.NewV4()
	}

	file.CreatedAt = time.Now()
	file.UpdatedAt = time.Now()

	if err := repo.app.Database.Insert(file); err != nil {
		return errors.Wrap(err, database.UnexpectedDatabaseErr)
	}

	return nil
}

func (repo *repository) Update(file *minerva.File) error {
	file.UpdatedAt = time.Now()

	if err := repo.app.Database.Update(file); err != nil {
		return errors.Wrap(err, database.UnexpectedDatabaseErr)
	}

	return nil
}

func (repo *repository) GetOldestUpdatedAtByHostnameAndAccumulatedSize(hostname minerva.Hostname, amount uint64) ([]*minerva.File, error) {

	files := []*minerva.File{}

	sql := `SELECT f.uuid, f.external_id, f.type, f.hostname, f.path, f.size, f.meta, f.pending_deletion, f.pending_upload, f.created_at, f.updated_at, f.original_filename, f.upstore_hash
FROM (SELECT
        f.*,
        sum(f.size)
        OVER (
          ORDER BY f.updated_at ) AS amount
      FROM files f
      WHERE f.hostname = ? AND f.pending_upload = FALSE AND f.pending_deletion = FALSE) f
WHERE amount - size < ?
ORDER BY f.updated_at DESC`

	_, err := repo.app.Database.Query(&files, sql, string(hostname), amount)

	switch err {
	case nil, pg.ErrNoRows:
		return files, nil

	default:
		repo.app.Logger.
			WithError(err).
			Error("Unexpected database error occurred")

		return nil, errors.Wrap(err, "Unexpected database error")
	}
}
