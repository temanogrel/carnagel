package process

import (
	"context"
	"git.misc.vee.bz/carnagel/encoder/pkg"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"github.com/satori/go.uuid"
	"os"
)

func InfinityGenerateCollage(ctx context.Context, app *encoder.Application, recording *ecosystem.Recording) (context.Context, error) {
	log := log(ctx, app)
	log.Debug("Generating infinity collage")

	video, ok := ctx.Value(encoder.CtxMp4File).(string)
	if !ok {
		return ctx, encoder.MissingMp4FileErr
	}

	target := encoder.ExtendFileName(video, "-inf-collage", "jpg")
	defer os.Remove(target)

	if err := app.ImageGenerationService.InfinityCollage(ctx, video, target); err != nil {
		return ctx, err
	}

	file, err := os.Open(target)
	if err != nil {
		return ctx, err
	}

	defer file.Close()

	// remove the existing infinity collage uuid if it exists
	if recording.InfinityCollageUuid != uuid.Nil {

		// it's not critical if this fails since it will be caught by the clean up process on the minions
		if err := app.MinervaClient.RequestDeletion(recording.InfinityCollageUuid); err != nil {
			log.WithError(err).Error("Failed to request deletion")
		}
	}

	id, err := app.MinervaClient.Upload(file, ecosystem.ExternalId(recording.Id), ecosystem.FileTypeInfinityCollage, nil)
	if err != nil {
		return ctx, err
	}

	recording.InfinityCollageUuid = id

	return ctx, nil
}

func InfinityGenerateThumbs(ctx context.Context, app *encoder.Application, recording *ecosystem.Recording) (context.Context, error) {
	log := log(ctx, app)
	log.Debug("Generating infinity thumbs")

	video, ok := ctx.Value(encoder.CtxMp4File).(string)
	if !ok {
		return ctx, encoder.MissingMp4FileErr
	}

	target := encoder.ExtendFileName(video, "-inf-thumbs", "jpg")
	defer os.Remove(target)

	err := app.ImageGenerationService.InfinityThumbs(ctx, video, target)
	if err != nil {
		return ctx, err
	}

	file, err := os.Open(target)
	if err != nil {
		return ctx, err
	}

	defer file.Close()

	id, err := app.MinervaClient.Upload(file, ecosystem.ExternalId(recording.Id), ecosystem.FileTypeInfinityImage, nil)
	if err != nil {
		return ctx, err
	}

	recording.Images = []uuid.UUID{id}

	return ctx, nil
}

func InfinityGenerateSprites(ctx context.Context, app *encoder.Application, recording *ecosystem.Recording) (context.Context, error) {
	log := log(ctx, app)
	log.Debug("Generating infinity sprites")

	video, ok := ctx.Value(encoder.CtxMp4File).(string)
	if !ok {
		return ctx, encoder.MissingMp4FileErr
	}

	target := encoder.ExtendFileName(video, "-inf-sprites", "jpg")
	defer os.Remove(target)

	err := app.ImageGenerationService.InfinitySprites(ctx, video, target)
	if err != nil {
		return ctx, err
	}

	file, err := os.Open(target)
	if err != nil {
		return ctx, err
	}

	defer file.Close()

	id, err := app.MinervaClient.Upload(file, ecosystem.ExternalId(recording.Id), ecosystem.FileTypeInfinitySprite, nil)
	if err != nil {
		return ctx, err
	}

	recording.Sprites = []uuid.UUID{id}

	return ctx, nil
}
