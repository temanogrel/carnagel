package infinity

import (
	"time"

	"context"

	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

const NumberOfBasicPlansUsablePerDayForRemoteAddr = 2

var BlackListedSessionErr = errors.New("The provided session is black listed remainder of the day")

type BandwidthConsumptionCollector interface {
	// Synchronizer keeps the in-memory database synced with the real database at constant intervals and when the app
	// terminates so we don't loose data.
	Synchronizer(ctx context.Context)

	// RebuildFromToday queries the database and builds an in-memory database of today's usage so that we don't end up
	// with billions of entries for each chunk loaded when watching a video
	RebuildFromToday() error

	// AddConsumption for a session
	AddConsumption(consumption *BandwidthConsumption) (total, onRecording uint64)
	AddBlackListedSession(session uuid.UUID)

	// GetTotalConsumptionToday returns the in-memory number of bytes the user has consumed today
	GetTotalConsumptionToday(session uuid.UUID) (uint64, error)

	// GetTotalConsumptionTodayOnRecording returns the in-memory number of bytes the user has consumed on a specific recording today
	GetTotalConsumptionTodayOnRecording(session uuid.UUID, recordingId uuid.UUID) (uint64, error)
}

type BandwidthConsumptionRepository interface {
	GetTodaysEntries() ([]*BandwidthConsumption, error)

	// AddConsumption persists the used consumption to the database, but it will try and aggregate the daily usage to
	// minimize the number of entries in the database
	AddConsumption(consumption *BandwidthConsumption) error

	// GetConsumptionToday will return the usage of a session today
	GetConsumptionToday(session uuid.UUID) (uint64, error)
}

type BandwidthStatus struct {
	Remaining int64  `json:"remaining"`
	Total     uint64 `json:"total"`
}

type BandwidthConsumption struct {
	TableName struct{} `sql:"bandwidth_consumption,alias:bc"`

	UserUuid      uuid.NullUUID `json:"userUuid"`
	SessionUuid   uuid.UUID     `json:"sessionUuid"`
	RecordingUuid uuid.UUID     `json:"recordingId"`
	Bytes         uint64        `json:"bytes"`
	CreatedAt     time.Time     `json:"createdAt"`
}
