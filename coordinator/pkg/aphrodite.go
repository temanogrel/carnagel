package coordinator

import "context"

type RecordingRepository interface {
	GetPublishedRecordingIds(ctx context.Context, lastSeenId uint64) (<-chan uint64, error)
}
