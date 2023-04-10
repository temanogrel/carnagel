package infinity

import (
	"context"
)

type NonDeletableRecording struct {
	TableName struct{} `sql:"non_deletable_recordings"`

	UpstoreHash string
	Filename    string
}

type CleanupService interface {
	Run(ctx context.Context)
	RunCleanup(ctx context.Context)
}

type CleanupRepository interface {
	UpdateNonDeletableRecordings(ctx context.Context, recordings []NonDeletableRecording) error
	IsUpstoreHashDeletable(ctx context.Context, hash string) (bool, error)
}
