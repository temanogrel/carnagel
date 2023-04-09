package pipeline

import (
	"context"
	"git.misc.vee.bz/carnagel/encoder/pkg"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
)

type SerialPipeline []encoder.PipelineFunc

func (pipeline SerialPipeline) Process(ctx context.Context, app *encoder.Application, recording *ecosystem.Recording) (context.Context, error) {
	var err error

	for _, f := range pipeline {
		ctx, err = f.Process(ctx, app, recording)
		if err != nil {
			return ctx, err
		}
	}

	return ctx, nil
}
