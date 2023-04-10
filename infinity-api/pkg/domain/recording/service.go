package recording

import (
	"context"

	"fmt"

	"time"

	"encoding/json"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	"github.com/gosimple/slug"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type service struct {
	app *infinity.Application
	log logrus.FieldLogger
}

func NewService(app *infinity.Application) infinity.RecordingService {
	service := &service{
		app: app,
		log: app.Logger.WithField("service", "file"),
	}

	go service.updateRecordingsWithoutSlug()

	return service
}

func (service *service) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		default:
			validateRecordings := service.app.Consul.GetBool("infinity/validate-recordings", false)
			if !validateRecordings {
				// Sleep for one minute to avoid spamming consul with requests
				time.Sleep(time.Minute)

				continue
			}

			service.validateRecordings()
		}
	}
}

func (service *service) queryAndDeleteRecordingsNotFoundInMinerva(recordings []infinity.Recording) {
	logger := service.log.WithField("method", "queryAndDeleteRecordingsNotFoundInMinerva")

	amqpCh, err := ecosystem.GetAmqpChannel(service.app.Amqp, 0)
	if err != nil {
		logger.WithError(err).Error("Failed to create amqp channel")
		return
	}

	defer amqpCh.Close()

	var recordingsToMarkAsValidated []uuid.UUID

	ids := make([]ecosystem.ExternalId, len(recordings))
	for i, recording := range recordings {
		ids[i] = ecosystem.ExternalId(recording.ExternalId)
	}

	result, err := service.app.MinervaClient.HasBulk(ids)
	if err != nil {
		logger.WithError(err).Error("Failed to connect to minerva")
		return
	}

	for _, recording := range recordings {
		logger := service.log.WithField("recordingUuid", recording.Uuid)

		exists, ok := result[ecosystem.ExternalId(recording.ExternalId)]
		if !ok {
			continue
		}

		if exists {
			recordingsToMarkAsValidated = append(recordingsToMarkAsValidated, recording.Uuid)

			continue
		}

		logger.Debug("Deleting recording")

		// Delete it in our repository first to avoid in premium user collection / view errors
		// when calling delete from coordinator
		_, err := service.app.RecordingRepository.Remove(recording.Uuid)
		if err != nil {
			logger.WithError(err).Error("Failed to delete recording from database")

			continue
		}

		body, err := json.Marshal(&ecosystem.SystemScopedRecordingPayload{
			RecordingId: uint64(recording.ExternalId),
			System:      ecosystem.SystemInfinity,
		})

		if err != nil {
			logger.WithError(err).Error("Failed to encode amqp payload")

			continue
		}

		payload := amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		}

		if err := amqpCh.Publish("", ecosystem.AmqpQueueRemoveRecording, false, false, payload); err != nil {
			logger.WithError(err).Errorf("Failed to publish to queue %s", ecosystem.AmqpQueueRemoveRecording)
		}
	}

	if err := service.app.RecordingRepository.MarkAsValidated(recordingsToMarkAsValidated); err != nil {
		logger.WithError(err).Error("Failed to mark recordings as validated")
	}
}

func (service *service) validateRecordings() {
	logger := service.log.WithField("method", "validateRecordings")
	logger.Info("Validating recordings")

	var recordingsProcessed uint64

	criteria := &infinity.RecordingRepositoryCriteria{
		SortMode:      infinity.SortModeExternalIdAsc,
		ExternalIdMin: 0,
		Limit:         100,
		// Only retrieve recordings that haven't been validated the last 3 weeks
		LastValidationBefore: time.Now().Add(-time.Hour * 24 * 21),
	}

	for {
		recordings, _, err := service.app.RecordingRepository.Matching(criteria)
		if err != nil {
			logger.WithError(err).Error("Failed to retrieve recordings by criteria")
			return
		}

		service.queryAndDeleteRecordingsNotFoundInMinerva(recordings)

		// Add recordings processed
		recordingsProcessed += uint64(len(recordings))

		// Process recordings in bulks of maximum 25k
		if recordingsProcessed == 25000 {
			logger.Info("Successfully validated 10000 recordings against minerva")
			break
		}

		if len(recordings) < 100 {
			logger.Infof("Successfully validated %d recordings against minerva", recordingsProcessed)
			break
		}

		criteria.ExternalIdMin = recordings[len(recordings)-1].ExternalId + 1
		// Add slight sleep to reduce stress of server
		time.Sleep(100)
	}
}

func (service *service) DeleteRecording(ctx context.Context, id uint64) error {
	log := service.log.WithField("RecordingId", id)
	log.Debug("Deleting recording")

	recording, err := service.app.RecordingRepository.GetByExternalId(id)
	if err != nil {
		if err == infinity.RecordingNotFoundErr {
			// just return success code if the recording doesn't exist
			return nil
		}

		return err
	}

	if err := service.app.RecordingRepository.CanRemove(recording.Uuid); err != nil {
		return err
	}

	_, err = service.app.RecordingRepository.Remove(recording.Uuid)

	return err
}

func (service *service) ImportRecording(ctx context.Context, id uint64) (*infinity.Recording, *infinity.Performer, error) {
	log := service.log.WithField("RecordingId", id)
	log.Debug("Upserting recording")

	aphroditeRecording, err := service.app.AphroditeClient.GetRecording(id)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to retrieve recording from aphrodite")
	}

	aphroditePerformer, err := service.app.AphroditeClient.GetPerformer(aphroditeRecording.PerformerId)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to retrieve performer from ")
	}

	performer, err := service.app.PerformerService.ImportFromAphrodite(aphroditePerformer)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to import performer from aphrodite")
	}

	recording, err := service.ImportFromAphrodite(aphroditeRecording, performer)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to import recording")
	}

	return recording, performer, nil
}

func (service *service) ImportFromAphrodite(raw *ecosystem.Recording, performer *infinity.Performer) (*infinity.Recording, error) {

	recording, err := service.app.RecordingRepository.GetByExternalId(raw.Id)
	switch err {
	case nil:
		if err := service.updateRecordingWithSlug(recording, performer); err != nil {
			service.app.Logger.WithError(err).Error("Failed to update recording with slug")

			return nil, err
		}

		return service.updateFromAphrodite(raw, recording)

	case infinity.RecordingNotFoundErr:
		recording := &infinity.Recording{
			ExternalId:    raw.Id,
			StageName:     raw.StageName,
			VideoUuid:     raw.VideoHlsUuid,
			VideoManifest: raw.VideoManifest,
			CollageUuid:   raw.InfinityCollageUuid,
			Sprites:       raw.Sprites,
			Images:        raw.Images,
			Duration:      raw.Duration,
			CreatedAt:     time.Now(),
		}

		if err := service.updateRecordingWithSlug(recording, performer); err != nil {
			service.app.Logger.WithError(err).Error("Failed to update recording with slug")

			return nil, err
		}

		return recording, service.app.RecordingRepository.Create(recording, performer)

	default:
		return nil, errors.Wrap(err, "Failed to import from aphrodite due to database error")
	}
}

func (service *service) updateFromAphrodite(raw *ecosystem.Recording, recording *infinity.Recording) (*infinity.Recording, error) {
	recording.CollageUuid = raw.InfinityCollageUuid
	recording.Sprites = raw.Sprites
	recording.Images = raw.Images
	recording.StageName = raw.StageName

	return recording, service.app.RecordingRepository.Update(recording)
}

func (service *service) updateRecordingsWithoutSlug() {
	log := service.app.Logger.WithField("operation", "updateRecordingsWithoutSlug")
	log.Info("Updating recordings missing slug")

	recordings, err := service.app.RecordingRepository.GetAllMissingSlug()
	if err != nil {
		log.WithError(err).Error("Failed to retrieve recordings without slug")
		return
	}

	for _, recording := range recordings {
		performer, err := service.app.PerformerRepository.GetByUuid(recording.PerformerUuid)
		if err != nil {
			log.WithError(err).Error("Failed to retrieve performer of recording")
			return
		}

		if err := service.updateRecordingWithSlug(&recording, performer); err != nil {
			log.WithError(err).Error("Failed to generate slug for recording")
		}

		if err := service.app.RecordingRepository.Update(&recording); err != nil {
			log.WithError(err).Error("Failed to update recording with slug")
		}
	}

	log.Info("Finished updating recordings missing slug")
}

// Code is basically ported straight from ultron
func (service *service) updateRecordingWithSlug(recording *infinity.Recording, performer *infinity.Performer) error {
	serviceName, err := performer.GetFullOriginServiceName()
	if err != nil {
		return err
	}

	dateString := recording.CreatedAt.Format("020106 1504")

	var formatString string

	if performer.OriginService == infinity.OriginServiceChaturbate {
		formatString = fmt.Sprintf("%s %s %s %s", recording.StageName, dateString, serviceName, performer.OriginSection.String)
	} else {
		formatString = fmt.Sprintf("%s %s %s", recording.StageName, dateString, serviceName)
	}

	generatedSlug := slug.Make(formatString)

	if recording.Slug == generatedSlug {
		return nil
	}

	for i := 1; true; i++ {
		if _, err := service.app.RecordingRepository.GetBySlug(generatedSlug); err != nil {
			if err == infinity.RecordingNotFoundErr {
				recording.Slug = generatedSlug
				return nil
			} else {
				return err
			}
		}

		generatedSlug = slug.Make(fmt.Sprintf("%s %d", formatString, i))
	}

	return nil
}
