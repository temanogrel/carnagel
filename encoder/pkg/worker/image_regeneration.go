package worker

import (
	"context"
	"encoding/json"
	"git.misc.vee.bz/carnagel/encoder/pkg"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"time"
)

type imageRegenerationWorker struct {
	app    *encoder.Application
	logger logrus.FieldLogger
}

func NewImageRegenerationWorker(application *encoder.Application) encoder.Worker {
	return &imageRegenerationWorker{
		app:    application,
		logger: application.Logger.WithField("component", "worker"),
	}
}

func (worker *imageRegenerationWorker) Run(ctx context.Context, routines int) {
	ch, err := ecosystem.GetAmqpChannel(worker.app.Amqp, 1)
	if err != nil {
		worker.app.Logger.WithError(err).Fatal("Failed to create amqp channel")
	}

	defer ch.Close()

	consumingChannel, err := ch.Consume(ecosystem.AmqpQueueImageRegeneration, worker.app.Hostname, false, false, false, false, nil)
	if err != nil {
		worker.app.Logger.WithError(err).Fatal("Failed to consume from channel")
	}

	for i := 1; i <= routines; i++ {
		go worker.processMessages(ctx, consumingChannel)
	}

	// Wait for context to close
	<-ctx.Done()
}

func (worker *imageRegenerationWorker) processMessages(ctx context.Context, ch <-chan amqp.Delivery) {
	worker.app.Logger.Debug("Started new process")

	for {
		select {
		case <-ctx.Done():
			worker.app.Logger.Warn("Context closed")
			return

		default:
			ok, err := encoder.HasEnoughSpace(1024 * 1024 * 1024 * 5)
			if err != nil {
				worker.app.Logger.WithError(err).Error("Failed to check available disk space")
			}

			// Not enough space or an error curred sleep and try again
			if err != nil || !ok {
				worker.app.Logger.WithError(err).Debug("Not enough disk space, waiting 5 seconds")
				time.Sleep(time.Second * 5)
				continue
			}

			msg, ok := <-ch
			if !ok {
				worker.app.Logger.Fatal("Incoming amqp.Delivery channel closed")
				return
			}

			payload := ecosystem.SystemScopedRecordingPayload{}
			if err := json.Unmarshal(msg.Body, &payload); err != nil {
				msg.Ack(false)

				worker.app.Logger.
					WithError(err).
					Error("Failed to unmarshal json payload")

				continue
			}

			log := worker.app.Logger.WithField("recordingId", payload.RecordingId)
			log.Info("Regenerate images for recording")

			recording, err := worker.app.AphroditeClient.GetRecording(payload.RecordingId)
			if err != nil {

				// pointless to requeue recording not found errors
				requeueRecording := err != ecosystem.RecordingNotFoundErr

				msg.Nack(false, requeueRecording)

				log.
					WithError(err).
					Error("Failed to retrieve recording")

				continue
			}
			
			msg.Ack(false)

			if err := worker.app.Pipelines.RegenerateImages(ctx, recording, payload.System); err != nil {
				log.
					WithError(err).
					Error("Failed image regeneration of recording")

				continue
			}

			log.Info("Successfully regenerated images for recording")
		}
	}
}
