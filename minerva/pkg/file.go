package minerva

import (
	"time"

	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

var (
	FileNotFoundErr               = errors.New("The file was not found")
	FileAlreadyPendingUploadErr   = errors.New("The file is already pending upload")
	FileAlreadyPendingDeletionErr = errors.New("The file is already pending deletion")
)

type File struct {
	Uuid        uuid.UUID            `sql:",pk" json:"uuid"`
	Type        ecosystem.FileType   `sql:",notnull" json:"type"`
	ExternalId  ecosystem.ExternalId `json:"recordingUuid"`
	UpstoreHash string               `json:"upstoreHash"`
	Checksum    string               `json:"checksum"`

	PendingDeletion bool `sql:",notnull" json:"pendingDeletion"`
	PendingUpload   bool `sql:",notnull" json:"pendingUpload"`

	// Location
	Hostname         Hostname `sql:",notnull" json:"hostname"`
	Path             string   `sql:",notnull" json:"path"`
	OriginalFilename string   `json:"originalFilename"`

	// Meta
	Size uint64                 `sql:",notnull" json:"size"`
	Meta ecosystem.FileMetadata `json:"meta"`

	// Timestamp
	CreatedAt time.Time `sql:",notnull" json:"createdAt"`
	UpdatedAt time.Time `sql:",notnull" json:"updatedAt"`
}

type FileHit struct {
	Uuid     uuid.UUID `json:"uuid"`
	FileUuid uuid.UUID `json:"fileUuid"`

	RemoteAddress string `json:"remoteAddress"`

	CreatedAt time.Time `json:"createdAt"`
}

type FileRepository interface {
	GetByExternalIds(ids []ecosystem.ExternalId) ([]*File, error)
	GetByUuid(id uuid.UUID) (*File, error)
	GetByLocation(hostname, path string) (*File, error)
	GetByExternalId(id ecosystem.ExternalId) ([]*File, error)
	GetWithPendingOperations(hostname Hostname, limit uint64) ([]*File, error)
	GetOldestUpdatedAtByHostnameAndAccumulatedSize(hostname Hostname, amount uint64) ([]*File, error)

	TrackHit(id uuid.UUID, remoteAddr string) error

	Create(file *File) error
	Update(file *File) error
	Delete(file *File) error
}

type FileService interface {
	// CleanUp cleans up the file service by closing the amqp connection
	CleanUp()

	// Delete removes a file from the database
	Delete(id uuid.UUID) error

	// ScheduleDeletion marks a file for deletion and if the server is online pushes to it's DeletionRequest channel
	ScheduleDeletion(id uuid.UUID) error

	// ScheduleDeletion removes all files with the corresponding external and file type unless file type is TYPE_ALL
	// In which case it will remove all files
	ScheduleDeletionByExternalIdAndType(id ecosystem.ExternalId, fileType ecosystem.FileType) error

	// ScheduleUpload marks a file for uploading and if the server ie online pushes to it's UploadRequest channel
	ScheduleUpload(id uuid.UUID, name string) error

	// SetUpstoreHash sets the upstore hash of the file, and marks the file no longer pending upload
	SetUpstoreHash(id uuid.UUID, hash string) error

	// Create CRUD method
	Create(data *CreateFile) (*File, error)

	// Update CRUD method
	Update(data *UpdateFile) (*File, error)
}

type UpdateFile struct {
	Uuid uuid.UUID

	// Location
	Hostname string
	Path     string

	// Meta
	Size     uint64
	Meta     ecosystem.FileMetadata
	Checksum string
}

// CreateFile is a normalization of the fields required to create a new file entry in the database
type CreateFile struct {
	ExternalId      ecosystem.ExternalId
	Type            ecosystem.FileType
	PendingDeletion bool

	// Location
	Hostname string
	Path     string

	// Meta
	Size     uint64
	Meta     ecosystem.FileMetadata
	Checksum string
}
