package infinity

import (
	"context"
	"errors"
	"time"

	ecosystem "git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"github.com/satori/go.uuid"
)

var (
	RecordingHasViewsErr                = errors.New("The recording has views")
	RecordingInPremiumUserCollectionErr = errors.New("The recording is part of a premium user's collection")
	RecordingNotFoundErr                = errors.New("The recording was not found")
)

const (
	SortModePopularity    = "popularity"
	SortModeViews         = "views"
	SortModeOldest        = "oldest"
	SortModeExternalIdAsc = "external-id:asc"
)

type Recording struct {
	TableName struct{} `sql:"recordings,alias:recording" json:"-"`

	Uuid       uuid.UUID `sql:",pk" json:"uuid"`
	ExternalId uint64    `sql:",notnull" json:"externalId"`

	StageName string `json:"stageName"`

	VideoUuid     uuid.UUID   `json:"videoUuid"`
	VideoManifest string      `json:"-"`
	CollageUuid   uuid.UUID   `json:"collageUuid"`
	Sprites       []uuid.UUID `pg:",array" json:"sprites"`
	Images        []uuid.UUID `pg:",array" json:"images"`

	Slug string `json:"slug"`

	ViewCount uint `sql:",notnull" json:"viewCount"`
	LikeCount uint `sql:",notnull" json:"likeCount"`

	Duration float64 `json:"duration"`

	UpdatedAt        time.Time  `json:"updatedAt"`
	CreatedAt        time.Time  `json:"createdAt"`
	LastValidationAt *time.Time `json:"-"`

	// Relationship
	Performer     *Performer `json:"performer,omitempty"`
	PerformerUuid uuid.UUID  `sql:",notnull" json:"performerId"`
}

func (r *Recording) IsMinified() bool {
	return len(r.Images) == 1
}

type RecordingView struct {
	TableName struct{} `sql:"recording_views"`

	RecordingUuid uuid.UUID
	UserUuid      uuid.NullUUID
	UserHash      string
	CreatedAt     time.Time
}

type RecordingWithUserData struct {
	Recording `pg:",override"`

	IsLiked    bool `json:"isLiked"`
	IsFavorite bool `json:"isFavorite"`
}

type RecordingRepositoryCriteria struct {
	PerformerId   uuid.UUID
	StageName     string
	OriginSection string
	OriginService string

	Limit  int
	Offset int

	FavoritesOnly bool

	CreatedAfter  time.Time
	CreatedBefore time.Time

	ExternalIdMin uint64
	ExternalIdMax uint64

	LastValidationBefore time.Time

	// This is used by the filtering of popularity
	SortMode string
	Interval uint8
}

type RecordingRepository interface {
	LoadPageCacheFromRedis()
	GetByUuid(id uuid.UUID) (*Recording, error)
	GetByExternalId(id uint64) (*Recording, error)
	GetUuidByExternalId(id uint64) (uuid.UUID, error)
	GetBySlug(slug string) (*Recording, error)

	GetAllMissingSlug() ([]Recording, error)

	AddView(*Recording, string) (bool, error)

	RebuildPageCache(fromPage uint64)

	CanRemove(id uuid.UUID) error
	Remove(id uuid.UUID) (bool, error)
	Update(recording *Recording) error
	Create(*Recording, *Performer) error
	MarkAsValidated(ids []uuid.UUID) error

	// Matching retrieve all recordings matching the provided criteria
	Matching(*RecordingRepositoryCriteria) ([]Recording, int, error)

	// Optimized
	GetRecordingsCreatedBefore(time time.Time, limit int) ([]Recording, error)

	// GetByUuid methods but with a user context, this is so we can have the IsLiked/IsFavorite property populated
	GetBySlugWithUserContext(slug string, user uuid.UUID) (*RecordingWithUserData, error)
	GetByUuidWithUserContext(recording uuid.UUID, user uuid.UUID) (*RecordingWithUserData, error)
	MatchingWithUserContext(criteria *RecordingRepositoryCriteria, user uuid.UUID) ([]RecordingWithUserData, int, error)
}

type RecordingService interface {
	Run(context.Context)
	DeleteRecording(context.Context, uint64) error
	ImportRecording(context.Context, uint64) (*Recording, *Performer, error)
	ImportFromAphrodite(recording *ecosystem.Recording, performer *Performer) (*Recording, error)
}
