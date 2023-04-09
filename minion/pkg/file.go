package minion

import (
	"io"

	"context"

	"git.misc.vee.bz/carnagel/minerva-bindings/src"
	"github.com/satori/go.uuid"
)

type FileService interface {

	// HandleUpload handles the upload from the http server, placing it in a random location within the data directory
	// and creates the file in minerva and returning the created file information
	HandleUpload(file io.Reader, fileName string, externalId uint64, fileType pb.FileType) (uuid.UUID, error)

	// HandleTransfer Places in the incoming file from the request into a random directory and if placed successfully
	// updates the file in minerva before returning a nil error state
	HandleTransfer(ctx context.Context, id uuid.UUID, file io.Reader, extension string) error

	// Transfer Moves a file from the local server to another server which is then handled by HandlerTransfer
	Transfer(ctx context.Context, id uuid.UUID, targetHost string) error
}
