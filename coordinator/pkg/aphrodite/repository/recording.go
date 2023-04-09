package repository

import (
	"context"

	"time"

	"git.misc.vee.bz/carnagel/coordinator/pkg"
	"github.com/pkg/errors"
)

type recordingRepository struct {
	app *coordinator.Application
}

func NewRecordingRepository(app *coordinator.Application) coordinator.RecordingRepository {
	return &recordingRepository{app: app}
}

func (repo *recordingRepository) GetPublishedRecordingIds(ctx context.Context, lastSeenId uint64) (<-chan uint64, error) {

	startedAt := time.Now()

	rows, err := repo.app.AphroditeDb.QueryContext(ctx, "SELECT id FROM recordings WHERE state = 'published' and bit_rate IS NOT NULL AND id > ?", lastSeenId)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to retrieve recordings")
	}

	repo.app.Logger.WithField("duration", time.Since(startedAt).String()).Debug("Executed query")

	// probably a bit excessive but i think it's a read timeout causing our problems
	ch := make(chan uint64, 12000000)

	go func() {
		defer rows.Close()
		defer close(ch)

		for {
			select {
			case <-ctx.Done():
				return

			default:
				if !rows.Next() {
					if rows.Err() != nil {
						repo.app.Logger.WithError(rows.Err()).Error("An unexpected error occurred in the result set loop")
					}

					return
				}

				var recordingId uint64
				if err := rows.Scan(&recordingId); err != nil {
					repo.app.Logger.WithError(err).Error("Failed to scan row")
					return
				}

				ch <- recordingId
			}
		}
	}()

	return ch, nil
}
