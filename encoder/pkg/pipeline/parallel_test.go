package pipeline

import (
	"context"
	"git.misc.vee.bz/carnagel/encoder/pkg"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParallelPipeline_Process(t *testing.T) {

	f := func(ctx context.Context, app *encoder.Application, rec *ecosystem.Recording) (context.Context, error) {
		return ctx, nil
	}

	pipeline := ParallelPipeline{
		f,
		f,
	}

	ctx := context.TODO()
	app := &encoder.Application{}
	rec := &ecosystem.Recording{}

	ctx2, err := pipeline.Process(ctx, app, rec)

	assert.Nil(t, err)
	assert.Equal(t, ctx, ctx2)
}

func TestParallelPipeline_ProcessWithError(t *testing.T) {

	success := func(ctx context.Context, app *encoder.Application, rec *ecosystem.Recording) (context.Context, error) {
		return ctx, nil
	}

	failure := func(ctx context.Context, app *encoder.Application, rec *ecosystem.Recording) (context.Context, error) {
		return nil, errors.New("I failed")
	}

	pipeline := ParallelPipeline{
		success,
		failure,
	}

	ctx := context.TODO()
	app := &encoder.Application{}
	rec := &ecosystem.Recording{}

	ctx2, err := pipeline.Process(ctx, app, rec)

	assert.Error(t, err, "i failed")
	assert.Equal(t, ctx, ctx2)
}
