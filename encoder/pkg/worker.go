package encoder

import (
	"context"

	"github.com/pkg/errors"
)

var (
	WorkerNotRunningErr     = errors.New("The worker is not running")
	WorkerAlreadyRunningErr = errors.New("The worker is already running")
)

type Worker interface {
	Run(ctx context.Context, routines int)
}
