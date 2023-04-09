package process

import (
	"context"
	"encoding/json"
	"git.misc.vee.bz/carnagel/encoder/pkg"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"os"
	"time"
)

func log(ctx context.Context, app *encoder.Application) logrus.FieldLogger {
	log := app.Logger.WithFields(logrus.Fields{
		"recordingId": ctx.Value("recordingId").(uint64),
		"processUuid": ctx.Value("processUuid").(string),
		"processCtx": map[string]interface{}{
			"mp4": ctx.Value(encoder.CtxMp4File),
			"flv": ctx.Value(encoder.CtxFlvFile),
			"hls": ctx.Value(encoder.CtxHlsFile),
		},
	})

	return log
}

func DeleteLocalFiles(ctx context.Context, app *encoder.Application, rec *ecosystem.Recording) (context.Context, error) {
	keys := []string{
		encoder.CtxHlsFile,
		encoder.CtxFlvFile,
		encoder.CtxMp4File,
	}

	for _, key := range keys {
		if file, ok := ctx.Value(key).(string); ok {
			os.Remove(file)
		}
	}

	return ctx, nil
}

func DispatchPublish(ctx context.Context, app *encoder.Application, rec *ecosystem.Recording) (context.Context, error) {
	ch, err := ecosystem.GetAmqpChannel(app.Amqp, 0)
	if err != nil {
		return ctx, err
	}

	defer ch.Close()

	body, err := json.Marshal(&ecosystem.FileUploadedPayload{
		ExternalId: rec.Id,
		Timestamp:  time.Now(),
	})

	if err != nil {
		return ctx, errors.Wrap(err, "Failed to encode fileUploaded jsn")
	}

	payload := amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	}

	if err := ch.Publish("", ecosystem.AmqpQueueUploaded, false, false, payload); err != nil {
		return ctx, errors.Wrap(err, "Failed to publish to the uploaded exchange")
	}

	return ctx, nil
}
