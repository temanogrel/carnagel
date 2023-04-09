package minion

import (
	"context"
	"sync"
	"time"
)

type IntegrityReport struct {
	*sync.RWMutex

	Scanned  uint64
	Orphaned uint64

	UpdatedAt time.Time
}

type FilesystemIntegrity interface {
	GetReport() IntegrityReport
	Run(ctx context.Context)

	IsRunning() bool
	IsScanning() bool
}
