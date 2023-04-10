package cleanup

import (
	"context"
	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	"github.com/go-pg/pg"
	"github.com/pkg/errors"
)

type cleanupRepository struct {
	app *infinity.Application
}

func NewCleanupRepository(app *infinity.Application) infinity.CleanupRepository {
	return &cleanupRepository{app}
}

func (r *cleanupRepository) UpdateNonDeletableRecordings(ctx context.Context, recordings []infinity.NonDeletableRecording) error {
	return r.app.DB.WithContext(ctx).RunInTransaction(func(tx *pg.Tx) error {
		if _, err := tx.Exec("DELETE FROM non_deletable_recordings"); err != nil {
			return errors.Wrapf(err, "Failed to delete existing non deletable recordings")
		}

		for _, recording := range recordings {
			if _, err := tx.Model(&recording).OnConflict("DO NOTHING").Insert(); err != nil {
				return errors.Wrapf(err, "Failed to insert non deletable recording")
			}
		}

		return nil
	})
}

func (r *cleanupRepository) IsUpstoreHashDeletable(ctx context.Context, hash string) (bool, error) {
	if hash == "" {
		return true, nil
	}

	recording := &infinity.NonDeletableRecording{}
	err := r.app.DB.WithContext(ctx).
		Model(recording).
		Where("upstore_hash = ?", hash).
		Select()

	switch err {
	case nil, pg.ErrNoRows:
		return err == pg.ErrNoRows, nil
	default:
		return false, err
	}
}
