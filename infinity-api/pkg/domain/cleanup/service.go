package cleanup

import (
	"context"
	"encoding/json"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"sync"
	"time"
)

type cleanupService struct {
	app    *infinity.Application
	logger logrus.FieldLogger
}

func NewCleanupService(app *infinity.Application) infinity.CleanupService {
	return &cleanupService{
		app:    app,
		logger: app.Logger.WithField("service", "CleanupService"),
	}
}

func (s *cleanupService) Run(ctx context.Context) {
	timer := time.NewTimer(0)

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Received <-ctx.Done(), stopping...")
			timer.Stop()
			return

		case <-timer.C:
			func() {
				defer timer.Reset(time.Hour * 24)

				if !s.app.Consul.GetBool("infinity/perform-cleanup", false) {
					return
				}

				s.RunCleanup(ctx)
			}()
		}
	}
}

func (s *cleanupService) RunCleanup(ctx context.Context) {
	logger := s.logger.WithField("operation", "RunCleanup")
	logger.Debug("Deleting recordings with no views and not in premium user collection")

	amqpCh, err := ecosystem.GetAmqpChannel(s.app.Amqp, 0)
	if err != nil {
		logger.WithError(err).Error("Failed to create amqp channel for cleanup background task")
		return
	}

	defer amqpCh.Close()

	stats := struct {
		deleted int
	}{}

	const limit = 120
	createdBefore := time.Now().Add(-time.Hour * 24 * 60)

	wg := sync.WaitGroup{}

	workerFunc := func(recording infinity.Recording) {
		defer wg.Done()

		if deleted, err := s.deleteRecording(ctx, amqpCh, recording); err != nil {
			logger.WithError(err).Error("Failed to remove recording")
		} else if deleted {
			logger.WithField("recordingUuid", recording.Uuid).Debug("Deleted recording")
			stats.deleted += 1
		}
	}

	for {
		// Use this instead of matching because count queries have a lot of overhead
		recordings, err := s.app.RecordingRepository.GetRecordingsCreatedBefore(createdBefore, limit)
		if err != nil {
			logger.WithError(err).Error("Failed to retrieve recordings created before")
			continue
		}

		for _, recording := range recordings {
			wg.Add(1)

			go workerFunc(recording)
		}

		wg.Wait()

		if len(recordings) < limit {
			break
		}
	}

	logger.WithField("stats", stats).Debug("Finished deleting recordings")
	// Invalidate the entire page cache
	s.app.RecordingRepository.RebuildPageCache(1)
}

func (s *cleanupService) isRecordingDeletable(ctx context.Context, recording infinity.Recording) (bool, error) {
	if err := s.app.RecordingRepository.CanRemove(recording.Uuid); err != nil {
		if err == infinity.RecordingHasViewsErr || err == infinity.RecordingInPremiumUserCollectionErr {
			return false, nil
		}

		return false, errors.Wrapf(err, "Failed to check if recording can be removed")
	}

	aphroditeRecording, err := s.app.AphroditeClient.GetRecording(recording.ExternalId)
	if err != nil {
		// The recording is not tracked by aphrodite, safe to remove
		if err == ecosystem.RecordingNotFoundErr {
			return true, nil
		}

		return false, errors.Wrapf(err, "Failed to retrieve recording from aphrodite by external id")
	}

	return s.app.CleanupRepository.IsUpstoreHashDeletable(ctx, aphroditeRecording.UpstoreHash)
}

func (s *cleanupService) deleteRecording(ctx context.Context, amqpCh *amqp.Channel, recording infinity.Recording) (bool, error) {
	deletable, err := s.isRecordingDeletable(ctx, recording)
	if err != nil || !deletable {
		return false, errors.Wrapf(err, "Failed to check if recording is deletable")
	}

	deleted, err := s.app.RecordingRepository.Remove(recording.Uuid)
	if err != nil {
		return false, errors.Wrapf(err, "Failed to remove recording")
	}

	body, err := json.Marshal(&ecosystem.SystemScopedRecordingPayload{
		RecordingId: uint64(recording.ExternalId),
		System:      ecosystem.SystemInfinity,
	})

	if err != nil {
		return false, errors.Wrapf(err, "Failed to encode amqp payload")
	}

	payload := amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	}

	if err := amqpCh.Publish("", ecosystem.AmqpQueueRemoveRecording, false, false, payload); err != nil {
		return false, errors.Wrapf(err, "Failed to publish to queue: %s", ecosystem.AmqpQueueRemoveRecording)
	}

	return deleted, nil
}
