package process

import (
	"context"
	"git.misc.vee.bz/carnagel/encoder/pkg"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"github.com/pkg/errors"
	"os"
)

func WordpressGenerateCollage(ctx context.Context, app *encoder.Application, rec *ecosystem.Recording) (context.Context, error) {
	log := log(ctx, app)
	log.Debug("Generating wordpress collage")

	video, ok := ctx.Value(encoder.CtxMp4File).(string)
	if !ok {
		return ctx, encoder.MissingMp4FileErr
	}

	target := encoder.ExtendFileName(video, "-wp-collage", "jpg")

	defer os.Remove(target)

	if err := app.ImageGenerationService.WordpressCollage(ctx, video, target, rec.GetPublishedFilename("mp4")); err != nil {
		return ctx, err
	}

	file, err := os.Open(target)
	if err != nil {
		return ctx, errors.Wrap(err, "Failed to open target for upload")
	}

	defer file.Close()

	id, err := app.MinervaClient.Upload(file, ecosystem.ExternalId(rec.Id), ecosystem.FileTypeWordpressCollage, nil)
	if err != nil {
		return ctx, err
	}

	rec.WordpressCollageUuid = id

	return ctx, nil
}
