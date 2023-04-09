package encoder

import (
	"context"

	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"github.com/pkg/errors"
)

const (
	CtxMp4File   = "file:mp4"
	CtxHlsFile   = "file:hls"
	CtxFlvFile   = "file:flv"
	CtxProbeMeta = "video:probe"
)

var (
	MissingProbeErr                   = errors.New("Missing a video meta, a probe must be called prior to this process")
	MissingMp4FileErr                 = errors.New("Missing the mp4 file")
	MissingHlsFileErr                 = errors.New("Missing the hls file")
	MissingFlvFileErr                 = errors.New("Missing the flv file")
	ShortCircuitAndDeleteRecordingErr = errors.New("Recording filed the pipeline resulting in it being deleted")
)

type Pipeline interface {
	Process(ctx context.Context, app *Application, recording *ecosystem.Recording) (context.Context, error)
}

type PipelineFunc func(ctx context.Context, app *Application, recording *ecosystem.Recording) (context.Context, error)

func (f PipelineFunc) Process(ctx context.Context, app *Application, recording *ecosystem.Recording) (context.Context, error) {
	return f(ctx, app, recording)
}

type ComposedPipelines interface {
	NewRecordingPipeline(ctx context.Context, rec *ecosystem.Recording) error
	ConvertMp42Hls(ctx context.Context, rec *ecosystem.Recording) error
	ReplaceCorruptHls(ctx context.Context, rec *ecosystem.Recording) error
	RegenerateImages(ctx context.Context, rec *ecosystem.Recording, system string) error
}
