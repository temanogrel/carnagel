package pipeline

import (
	"context"
	"sync"

	"git.misc.vee.bz/carnagel/encoder/pkg"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
)

type ParallelPipeline []encoder.PipelineFunc

func (pipeline ParallelPipeline) Process(ctx context.Context, app *encoder.Application, recording *ecosystem.Recording) (context.Context, error) {
	wg := sync.WaitGroup{}
	errCh := make(chan error)

	c, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, f := range pipeline {
		wg.Add(1)

		go func(segment encoder.PipelineFunc) {
			defer wg.Done()

			_, err := segment(c, app, recording)
			if err != nil {
				cancel()
				errCh <- err
			}
		}(f)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	err, ok := <-errCh
	if !ok {
		return ctx, nil
	}

	return ctx, err
}
