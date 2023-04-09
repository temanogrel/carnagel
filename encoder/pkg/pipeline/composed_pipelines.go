package pipeline

import (
	"context"
	"git.misc.vee.bz/carnagel/encoder/pkg"
	"git.misc.vee.bz/carnagel/encoder/pkg/pipeline/process"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

type composedPipelines struct {
	app *encoder.Application
}

func NewComposedPipelines(app *encoder.Application) encoder.ComposedPipelines {
	return &composedPipelines{
		app: app,
	}
}

func (cp *composedPipelines) exec(ctx context.Context, pipeline encoder.Pipeline, rec *ecosystem.Recording) error {

	ctx = context.WithValue(ctx, "recordingId", rec.Id)
	ctx = context.WithValue(ctx, "processUuid", uuid.NewV4().String())

	ctx, err := pipeline.Process(ctx, cp.app, rec)
	switch err {
	case nil:
		return nil

	case encoder.ShortCircuitAndDeleteRecordingErr:
		return cp.delete(ctx, rec)

	default:
		// Clean upp all local files on any random errors
		process.DeleteLocalFiles(ctx, cp.app, rec)

		cp.app.Logger.
			WithError(err).
			WithField("recordingId", rec.Id).
			WithField("component", "pipeline").
			Error("Failed to process pipeline")

		return err
	}
}

func (cp *composedPipelines) delete(ctx context.Context, rec *ecosystem.Recording) error {
	cp.app.Logger.
		WithField("recording", rec).
		Warn("Deleting recording")

	pipeline := SerialPipeline{}

	if rec.VideoMp4Uuid != uuid.Nil {
		pipeline = append(pipeline, process.MinionDeleteMp4)
	}

	if rec.VideoHlsUuid != uuid.Nil {
		pipeline = append(pipeline, process.MinionDeleteHls)
	}

	pipeline = append(pipeline,
		process.DeleteLocalFiles,
		process.MinionDeleteImages,
		process.AphroditeDeleteRecording,
	)

	if _, err := pipeline.Process(ctx, cp.app, rec); err != nil {
		return errors.Wrap(err, "Failed to cleanup recording")
	}

	return nil
}

func (cp *composedPipelines) NewRecordingPipeline(ctx context.Context, rec *ecosystem.Recording) error {
	pipeline := &SerialPipeline{
		process.MinionDownloadFlv,
		process.ProbeFlv,
		process.VideoEncodeFlv2Mp4,
		process.Probe,
		SerialPipeline{
			process.UpstoreUpload,
			process.WordpressGenerateCollage,
		}.Process,
		process.AphroditeUpdateRecording,
		process.MinionDeleteFlv,
		process.DeleteLocalFiles,
		process.DispatchPublish,
	}

	return cp.exec(ctx, pipeline, rec)
}

func (cp *composedPipelines) ConvertMp42Hls(ctx context.Context, rec *ecosystem.Recording) error {

	pipeline := &SerialPipeline{
		process.MinionDownloadMp4,
		process.MinionDeleteImages,
		process.Probe,
		SerialPipeline{
			SerialPipeline{
				process.VideoEncodeMp42Hls,
				process.MinionUploadHls,
			}.Process,
			process.InfinityGenerateCollage,
			process.InfinityGenerateThumbs,
			process.InfinityGenerateSprites,
			process.WordpressGenerateCollage,
		}.Process,
		process.AphroditeUpdateRecording,
		process.MinionDeleteMp4,
		process.DeleteLocalFiles,
		process.DispatchPublish,
	}

	return cp.exec(ctx, pipeline, rec)
}

func (cp *composedPipelines) ReplaceCorruptHls(ctx context.Context, rec *ecosystem.Recording) error {

	pipeline := &SerialPipeline{
		process.MinionDeleteHls,
		process.UpstoreDownload,
		process.VideoEncodeMp42Hls,
		process.MinionUploadHls,
		process.AphroditeUpdateRecording,
		process.DeleteLocalFiles,
		process.DispatchPublish,
	}

	return cp.exec(ctx, pipeline, rec)
}

func (cp *composedPipelines) RegenerateImages(ctx context.Context, rec *ecosystem.Recording, system string) error {
	var imageRegenerationPipeline ParallelPipeline

	switch system {
	case ecosystem.SystemInfinity:
		imageRegenerationPipeline = ParallelPipeline{
	
		}

	case ecosystem.SystemWordpress:
		imageRegenerationPipeline = ParallelPipeline{
			process.WordpressGenerateCollage,
		}

	default:
		imageRegenerationPipeline = ParallelPipeline{

			process.WordpressGenerateCollage,
		}
	}

	// Recreates images from corruption or other means
	pipeline := &SerialPipeline{
		process.MinionDownloadHls,
		process.VideoEncodeHls2Mp4,
		process.MinionDeleteImages,
		imageRegenerationPipeline.Process,
		process.AphroditeUpdateRecording,
		process.DeleteLocalFiles,
		process.DispatchPublish,
	}

	return cp.exec(ctx, pipeline, rec)
}
