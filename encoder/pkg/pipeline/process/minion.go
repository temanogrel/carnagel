package process

import (
	"context"
	"git.misc.vee.bz/carnagel/encoder/pkg"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"os"
)

func MinionDeleteMp4(ctx context.Context, app *encoder.Application, rec *ecosystem.Recording) (context.Context, error) {
	if rec.VideoMp4Uuid == uuid.Nil {
		return ctx, encoder.MissingMp4FileErr
	}

	if err := app.MinervaClient.RequestDeletion(rec.VideoMp4Uuid); err != nil {
		return ctx, err
	}

	return ctx, nil
}

func MinionDeleteHls(ctx context.Context, app *encoder.Application, rec *ecosystem.Recording) (context.Context, error) {
	if rec.VideoHlsUuid == uuid.Nil {
		return ctx, encoder.MissingHlsFileErr
	}

	if err := app.MinervaClient.RequestDeletion(rec.VideoHlsUuid); err != nil {
		return ctx, err
	}

	return ctx, nil
}

func MinionDeleteFlv(ctx context.Context, app *encoder.Application, rec *ecosystem.Recording) (context.Context, error) {
	if rec.VideoMp4Uuid == uuid.Nil {
		return ctx, encoder.MissingFlvFileErr
	}

	if err := app.MinervaClient.RequestDeletion(rec.VideoMp4Uuid); err != nil {
		return ctx, err
	}

	rec.VideoMp4Uuid = uuid.Nil

	return ctx, nil
}

func MinionDeleteImages(ctx context.Context, app *encoder.Application, rec *ecosystem.Recording) (context.Context, error) {
	types := []ecosystem.FileType{
		ecosystem.FileTypeInfinityCollage,
		ecosystem.FileTypeInfinityImage,
		ecosystem.FileTypeInfinitySprite,
		ecosystem.FileTypeWordpressCollage,
	}

	for _, t := range types {
		if err := app.MinervaClient.RequestDeletionByType(ecosystem.ExternalId(rec.Id), t); err != nil {
			return ctx, errors.Wrapf(err, "Failed to delete file type %d for %d", t, rec.Id)
		}
	}

	return ctx, nil
}

func MinionDownloadMp4(ctx context.Context, app *encoder.Application, rec *ecosystem.Recording) (context.Context, error) {
	file, err := app.MinervaClient.Download(rec.VideoMp4Uuid)
	if err != nil {
		return ctx, err
	}

	ctx = context.WithValue(ctx, encoder.CtxMp4File, file)

	return ctx, nil
}

func MinionDownloadHls(ctx context.Context, app *encoder.Application, rec *ecosystem.Recording) (context.Context, error) {
	file, err := app.MinervaClient.Download(rec.VideoHlsUuid)
	if err != nil {
		return ctx, err
	}

	ctx = context.WithValue(ctx, encoder.CtxHlsFile, file)

	return ctx, nil
}

func MinionDownloadFlv(ctx context.Context, app *encoder.Application, rec *ecosystem.Recording) (context.Context, error) {
	file, err := app.MinervaClient.Download(rec.VideoMp4Uuid)
	if err != nil {
		return ctx, err
	}

	ctx = context.WithValue(ctx, encoder.CtxFlvFile, file)

	return ctx, nil
}

func MinionUploadHls(ctx context.Context, app *encoder.Application, rec *ecosystem.Recording) (context.Context, error) {
	path, ok := ctx.Value(encoder.CtxHlsFile).(string)
	if !ok {
		return ctx, encoder.MissingHlsFileErr
	}

	file, err := os.Open(path)
	if err != nil {
		return ctx, errors.Wrap(err, "Failed to open hls file for upload")
	}

	defer file.Close()

	id, err := app.MinervaClient.Upload(
		file,
		ecosystem.ExternalId(rec.Id),
		ecosystem.FileTypeRecordingHls,
		ecosystem.FileMetadata{},
	)

	if err != nil {
		return ctx, errors.Wrapf(err, "Failed to upload hls to storage")
	}

	rec.VideoHlsUuid = id

	return ctx, nil
}
