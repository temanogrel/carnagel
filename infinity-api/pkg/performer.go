package infinity

import (
	"database/sql"
	"errors"
	"time"

	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"github.com/satori/go.uuid"
)

const OriginServiceChaturbate = "cbc"
const OriginServiceMyFreeCams = "mfc"

var (
	PerformerNotFoundErr    = errors.New("The performer was not found")
	UnknownOriginServiceErr = errors.New("Unknown origin service received")
)

type Performer struct {
	TableName struct{} `sql:"performers,alias:performer"`

	Uuid       uuid.UUID `sql:",pk" json:"uuid"`
	ExternalId uint64    `sql:",notnull" json:"externalId"`

	Slug string `json:"slug"`

	OriginService string         `sql:",notnull" json:"originService"`
	OriginSection sql.NullString `json:"originSection"`

	StageName string   `sql:",notnull" json:"stageName"`
	Aliases   []string `pg:",array" sql:",notnull" json:"aliases"`

	RecordingCount uint16 `sql:",notnull" json:"recordingCount"`

	UpdatedAt time.Time `sql:",notnull" json:"updatedAt"`
	CreatedAt time.Time `sql:",notnull" json:"createdAt"`

	// Relations
	Recordings []Recording `sql:"-" json:"recordings,omitempty"`
}

func (p *Performer) HasAlias(alias string) bool {
	for _, a := range p.Aliases {
		if a == alias {
			return true
		}
	}

	return false
}

func (p *Performer) GetFullOriginServiceName() (string, error) {
	switch p.OriginService {
	case OriginServiceChaturbate:
		return "Chaturbate", nil
	case OriginServiceMyFreeCams:
		return "MyFreeCams", nil
	default:
		return "", UnknownOriginServiceErr
	}
}

type PerformerRepositoryCriteria struct {
	Limit     int
	Offset    int
	StageName string

	IncludeLatestRecording bool
	PerformerIds           []uuid.UUID
}

type PerformerRepository interface {
	GetByUuid(id uuid.UUID) (*Performer, error)
	GetByExternalId(id uint64) (*Performer, error)
	GetBySlug(slug string) (*Performer, error)

	GetAllMissingSlug() ([]Performer, error)

	Create(performer *Performer) error
	Update(performer *Performer) error
	Matching(criteria *PerformerRepositoryCriteria) ([]*Performer, int, error)
}

type PerformerService interface {
	ImportFromAphrodite(performer *ecosystem.Performer) (*Performer, error)
}
