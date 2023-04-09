package process

import (
	"context"
	"git.misc.vee.bz/carnagel/encoder/pkg"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"github.com/pkg/errors"
	"os"
)

func UpstoreUpload(ctx context.Context, app *encoder.Application, recording *ecosystem.Recording) (context.Context, error) {
	log := log(ctx, app)
	log.Debug("Running UpstoreUpload")

	path, ok := ctx.Value(encoder.CtxMp4File).(string)
	if !ok {
		return ctx, encoder.MissingMp4FileErr
	}

	file, err := os.Open(path)
	if err != nil {
		return ctx, err
	}

	defer file.Close()

	hash, err := app.UpstoreClient.Upload(recording.GetPublishedFilename("mp4"), file)
	if err != nil {
		return ctx, errors.Wrap(err, "Failed to upload to upstore")
	}

	recording.UpstoreHash = hash

	return ctx, nil
}

func UpstoreDownload(ctx context.Context, app *encoder.Application, recordig *ecosystem.Recording) (context.Context, error) {
	return ctx, nil
}
