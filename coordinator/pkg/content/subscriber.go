package content

import (
	"context"
	"encoding/json"

	"time"

	"git.misc.vee.bz/carnagel/coordinator/pkg"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func NewContentSubscriber(application *coordinator.Application) coordinator.ContentSubscriber {
	return &contentSubscriber{
		app: application,
		log: application.Logger.WithField("component", "content"),
	}
}

type contentSubscriber struct {
	app *coordinator.Application
	log logrus.FieldLogger
}

func (subscriber *contentSubscriber) getChannel() (<-chan amqp.Delivery, *amqp.Channel, error) {
	ch, err := ecosystem.GetAmqpChannel(subscriber.app.Amqp, 60)
	if err != nil {
		return nil, nil, err
	}

	deliveryCh, err := ch.Consume(ecosystem.AmqpQueueUploaded, ecosystem.AmqpConsumerCoordinator, false, false, false, false, nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to consume from queue")
	}

	return deliveryCh, ch, nil
}

func (subscriber *contentSubscriber) Run(ctx context.Context) error {

	messages, channel, err := subscriber.getChannel()
	if err != nil {
		return errors.Wrap(err, "Failed to connect to amqp queue")
	}

	defer channel.Close()

	for {
		select {
		case <-ctx.Done():
			return nil

		case msg := <-messages:
			go subscriber.processMessage(msg)
		}
	}
}

func (subscriber *contentSubscriber) processMessage(message amqp.Delivery) {

	file := &ecosystem.FileUploadedPayload{}

	if err := json.Unmarshal(message.Body, file); err != nil {
		message.Nack(false, true)

		subscriber.log.
			WithError(err).
			WithFields(logrus.Fields{"body": string(message.Body), "messageId": message.MessageId}).
			Error("Failed to decode json")

		return
	}

	recording, err := subscriber.app.AphroditeClient.GetRecording(file.ExternalId)
	if err != nil {
		subscriber.log.WithError(err).Error("Failed to retrieve recording")

		// Assume it's temporary and requeue
		message.Nack(false, true)
		return
	}

	logger := subscriber.log.WithField("recording", recording)

	recording.State = "publishing"
	if err := subscriber.app.AphroditeClient.UpdateRecording(recording); err != nil {
		logger.WithError(err).Error("Failed to update state")

		// Assume it's temporary and requeue
		message.Nack(false, true)
		return
	}

	if err := subscriber.processRecording(recording); err == nil {
		recording.State = "published"
	} else if err == ecosystem.FailedToCallApiErr {
		// Assume temporary issue with api and requeue
		message.Nack(false, true)
		return
	} else {
		recording.State = "publification_failed"
	}

	if err := subscriber.app.AphroditeClient.UpdateRecording(recording); err != nil {
		logger.WithError(err).Error("Failed to update state")
	}

	message.Ack(false)
}

func (subscriber *contentSubscriber) processRecording(recording *ecosystem.Recording) error {

	logger := subscriber.log.WithField("recording", recording)
	logger.Debug("Started processing")

	startedAt := time.Now()

	if err := subscriber.app.ContentService.Publish(recording); err != nil {
		logger.WithError(err).Error("Failed to handle publification of recording")
		return err
	}

	logger.
		WithField("duration", time.Since(startedAt).Seconds()).
		Info("Publish file")

	return nil
}
