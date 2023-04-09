package cleanup

import (
	"context"
	"encoding/json"
	"time"

	"fmt"

	"git.misc.vee.bz/carnagel/coordinator/pkg"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"git.misc.vee.bz/carnagel/minerva/pkg"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type cleanupService struct {
	app *coordinator.Application
}

func NewCleanupService(application *coordinator.Application) coordinator.CleanupService {
	return &cleanupService{application}
}

func (service *cleanupService) Run(ctx context.Context) {
	ch, err := ecosystem.GetAmqpChannel(service.app.Amqp, 1)
	if err != nil {
		service.app.Logger.WithError(err).Error("Failed to get amqp channel")
	}

	defer ch.Close()

	consumingCh, err := ch.Consume(ecosystem.AmqpQueueRemoveRecording, ecosystem.AmqpConsumerCoordinator, false, false, false, false, nil)
	if err != nil {
		service.app.Logger.WithError(err).Fatal("Failed to consume from channel")
	}

	for {
		select {
		case <-ctx.Done():
			return

		default:
			msg, ok := <-consumingCh
			if !ok {
				service.app.Logger.Fatal("Incoming amqp.Delivery channel closed")
				return
			}

			service.processRecordingToDelete(msg)
		}
	}
}

func (service *cleanupService) processRecordingToDelete(msg amqp.Delivery) {
	logger := service.app.Logger.WithFields(logrus.Fields{
		"service": "CleanUpService",
		"method":  "processRecordingToDelete",
	})

	logger.Info("Received delete recording request")

	payload := ecosystem.SystemScopedRecordingPayload{}
	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		msg.Ack(false)

		logger.WithError(err).Error("Failed to unmarshal json payload")

		return
	}

	logger = logger.WithFields(logrus.Fields{"recordingId": payload.RecordingId, "systemScope": payload.System})
	logger.Debug("Attempting to delete recording")

	var recording *ecosystem.Recording

	if payload.System == ecosystem.SystemInfinity {
		recording = &ecosystem.Recording{Id: payload.RecordingId}
	} else {
		var err error

		recording, err = service.app.AphroditeClient.GetRecording(payload.RecordingId)
		if err == ecosystem.RecordingNotFoundErr {
			logger.Debug("Unable to delete recording because it was not found in aphrodite")
			msg.Nack(false, false)

			return
		} else if err != nil {
			msg.Nack(false, true)
			logger.WithError(err).Error("Failed to retrieve recording to delete")

			return
		}
	}

	msg.Ack(false)

	if err := service.Delete(recording, payload.System); err != nil {
		logger.WithError(err).Error("Failed to delete recording")
	} else {
		logger.Info("Successfully deleted the recording")
	}
}

func (service *cleanupService) TransferUpstoreHashToAphrodite() error {

	var timestamp time.Time

	stmt, err := service.app.AphroditeDb.Prepare("UPDATE recordings SET upstore_hash = ? WHERE id = ?")
	if err != nil {
		return err
	}

	for {
		var files []minerva.File

		err := service.app.MinervaDb.
			Model(&files).
			Where("file.upstore_hash IS NOT NULL AND file.created_at >= to_timestamp(?)", timestamp.Unix()).
			Order("created_at DESC").
			Limit(10000).
			Select()

		if err != nil {
			return err
		}

		if len(files) == 0 {
			return nil
		}

		for _, file := range files {
			timestamp = file.CreatedAt

			if _, err := stmt.Exec(file.UpstoreHash, file.ExternalId); err != nil {
				return err
			}

			service.app.Logger.WithFields(logrus.Fields{
				"recordingId": file.ExternalId,
				"upstoreHash": file.UpstoreHash,
			}).Info("Updated aphrodite recording")
		}
	}
}

func (service *cleanupService) DispatchMp42HlsConversion() error {
	ch, err := ecosystem.GetAmqpChannel(service.app.Amqp, 0)
	if err != nil {
		return err
	}

	defer ch.Close()

	if _, err := ch.QueuePurge(ecosystem.AmqpQueueMp42Hls, false); err != nil {
		return errors.Wrapf(err, "Failed to purge the queue %s", ecosystem.AmqpQueueMp42Hls)
	}

	rows, err := service.app.AphroditeDb.Query("SELECT id FROM recordings WHERE video_hls_uuid IS NULL AND state = 'published'")
	if err != nil {
		return errors.Wrap(err, "Failed to retrieve recordings")
	}

	defer rows.Close()

	for rows.Next() {
		var recordingId uint64

		err = rows.Scan(&recordingId)
		if err != nil {
			return errors.Wrap(err, "failed to scan recordingId column")
		}

		operation := &ecosystem.RecordingPayload{
			RecordingId: recordingId,
		}

		data, err := json.Marshal(operation)
		if err != nil {
			return errors.Wrap(err, "failed to create payload")
		}

		payload := amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		}

		if err := ch.Publish("", ecosystem.AmqpQueueMp42Hls, false, false, payload); err != nil {
			return errors.Wrap(err, "Failed to publish to mp42hls queue")
		}
	}

	return nil
}

func (service *cleanupService) RinsePublishingFailures(delete bool) error {
	ch, err := ecosystem.GetAmqpChannel(service.app.Amqp, 0)
	if err != nil {
		return err
	}

	defer ch.Close()

	states := []string{
		"uploaded",
		"publishing",
		"publishing_failed",
		"publification_failed",
	}

	for _, state := range states {

		after := time.Now().Add(-6 * time.Hour).Unix()

		params := ecosystem.RecordingCollectionQueryParams{
			State: state,
			After: uint64(after),
		}

		recordings, err := service.app.AphroditeClient.GetAllRecordings(context.TODO(), params)
		if err != nil {
			return errors.Wrapf(err, "Failed to get all recordings for state %s", state)
		}

		for recording := range recordings {
			if delete {
				service.Delete(recording, "")
				continue
			}

			service.app.Logger.WithField("recording", recording).Debug("Dispatching publishing")

			data := ecosystem.FileUploadedPayload{
				ExternalId: recording.Id,
				Timestamp:  time.Now(),
			}

			body, err := json.Marshal(data)
			if err != nil {
				service.app.Logger.WithError(err).Error("failed to marshal json")
				continue
			}

			payload := amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			}

			ch.Publish("", ecosystem.AmqpQueueUploaded, false, false, payload)
		}
	}

	return nil
}

func (service *cleanupService) RinseMp4Recordings() error {

	sql := `
		SELECT id, video_mp4uuid FROM recordings WHERE video_mp4uuid IS NOT NULL AND video_hls_uuid IS NOT NULL AND state = 'published'
	`

	rows, err := service.app.AphroditeDb.Query(sql)
	if err != nil {
		return errors.Wrap(err, "Failed to query the database")
	}

	for rows.Next() {
		var (
			recordingId uint64
			mp4Uuid     uuid.UUID
		)

		if err := rows.Scan(&recordingId, &mp4Uuid); err != nil {
			return errors.Wrap(err, "Failed to scan")
		}

		if err := service.app.MinervaClient.RequestDeletion(mp4Uuid); err != nil {
			return errors.Wrap(err, "Failed to request deletion")
		}

		if _, err := service.app.AphroditeDb.Exec("UPDATE recordings set video_mp4uuid = NULL WHERE id = ?", recordingId); err != nil {
			return errors.Wrap(err, "Failed to mark video mp4 uuid as null")
		}

		service.app.Logger.Info("Deleted mp4 from recording")
	}

	return nil
}

func (service *cleanupService) RinseEncodingFailures(delete bool) error {
	ch, err := ecosystem.GetAmqpChannel(service.app.Amqp, 0)
	if err != nil {
		return err
	}

	defer ch.Close()

	states := []string{
		"downloaded",
		"encoding",
		"encoding_failed",
	}

	for _, state := range states {

		after := time.Now().Add(-6 * time.Hour).Unix()

		params := ecosystem.RecordingCollectionQueryParams{
			State: state,
			After: uint64(after),
		}

		recordings, err := service.app.AphroditeClient.GetAllRecordings(context.TODO(), params)
		if err != nil {
			return errors.Wrapf(err, "Failed to get all recordings for state %s", state)
		}

		for recording := range recordings {
			if delete {
				service.Delete(recording, "")
				continue
			}

			service.app.Logger.WithField("recording", recording).Debug("Dispatching transcode")

			body := fmt.Sprintf(`[[%d], {}, {"callbacks": null, "errbacks": null, "chain": null, "chord": null}]`, recording.Id)

			payload := amqp.Publishing{
				ContentType: "application/json",
				Body:        []byte(body),
			}

			if err := ch.Publish("", ecosystem.AmqpQueueTranscode, false, false, payload); err != nil {
				service.app.Logger.WithError(err).Error("Failed to publish transcoding")
			}
		}
	}

	return nil
}

func (service *cleanupService) RinseUploadingFailures(delete bool) error {

	states := []string{
		"encoded",
		"uploading",
		"uploading_failed",
	}

	for _, state := range states {

		after := time.Now().Add(-6 * time.Hour).Unix()

		params := ecosystem.RecordingCollectionQueryParams{
			State: state,
			After: uint64(after),
		}

		recordings, err := service.app.AphroditeClient.GetAllRecordings(context.TODO(), params)
		if err != nil {
			return errors.Wrapf(err, "Failed to get all recordings for state %s", state)
		}

		for recording := range recordings {

			if delete {
				service.Delete(recording, "")
				continue
			}

			err := service.app.MinervaClient.RequestUpload(
				recording.VideoMp4Uuid,
				recording.GetPublishedFilename("mp4"),
			)

			if err != nil {
				service.app.Logger.WithError(err).Error("Failed to request upload")
			}
		}
	}

	return nil
}

func (service *cleanupService) deleteFromInfinity(id ecosystem.ExternalId) error {
	if err := service.app.MinervaClient.RequestDeletionByType(id, ecosystem.FileTypeRecordingHls); err != nil {
		return err
	}

	if err := service.app.MinervaClient.RequestDeletionByType(id, ecosystem.FileTypeRecordingMp4); err != nil {
		return err
	}

	if err := service.app.MinervaClient.RequestDeletionByType(id, ecosystem.FileTypeInfinityCollage); err != nil {
		return err
	}

	if err := service.app.MinervaClient.RequestDeletionByType(id, ecosystem.FileTypeInfinitySprite); err != nil {
		return err
	}

	if err := service.app.MinervaClient.RequestDeletionByType(id, ecosystem.FileTypeInfinityImage); err != nil {
		return err
	}

	return nil
}

func (service *cleanupService) deleteFromWordpress(recording *ecosystem.Recording) error {
	if err := service.app.MinervaClient.RequestDeletion(recording.WordpressCollageUuid); err != nil {
		return err
	}

	if err := service.app.AphroditeClient.DeleteRecording(recording.Id); err != nil {
		return err
	}

	return nil
}

func (service *cleanupService) Delete(recording *ecosystem.Recording, system string) error {
	logger := service.app.Logger.WithFields(logrus.Fields{"recording": recording, "systemScope": system})

	if err := service.app.ContentService.Remove(recording, system); err != nil {
		logger.WithError(err).Error("Failed to delete the published content")
		return err
	}

	switch system {
	case "":
		// If system is not specified, remove from both system
		if err := service.deleteFromInfinity(ecosystem.ExternalId(recording.Id)); err != nil {
			logger.WithError(err).Error("Failed to delete files associated with infinity")
			return err
		}

		if err := service.deleteFromWordpress(recording); err != nil {
			logger.WithError(err).Error("Failed to delete files associated with wordpress")
			return err
		}

	case ecosystem.SystemInfinity:
		if err := service.deleteFromInfinity(ecosystem.ExternalId(recording.Id)); err != nil {
			logger.WithError(err).Error("Failed to delete files associated with infinity")
			return err
		}

	case ecosystem.SystemWordpress:
		if err := service.deleteFromWordpress(recording); err != nil {
			logger.WithError(err).Error("Failed to delete files associated with wordpress")
			return err
		}
	}

	logger.Debug("Successfully deleted recording")

	return nil
}
