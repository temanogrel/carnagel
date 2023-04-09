package coordinator

import (
	"context"

	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"github.com/pkg/errors"
)

var (
	UpstoreHashMissingErr = errors.New("The upstore hash for the file is missing")

	ContentSubscriberNotRunningErr     = errors.New("Cannot stop the content subscriber, it's not running")
	ContentSubscriberAlreadyRunningErr = errors.New("Cannot started the content subscriber it's already running")
)

type ContentService interface {
	// Publish will publish the recording as new content as long as it does not already exist
	Publish(recording *ecosystem.Recording) error

	// Remove will remove the content from all sites
	Remove(recording *ecosystem.Recording, system string) error

	// RebuildUltron will attempt to rebuild the entire ultron database
	RebuildUltron(lastSeenId uint64) error
}

type ContentSubscriber interface {
	// Start will subscribe to the RabbitMQ uploaded queue, parse it and push it to the appropriate content method
	Run(ctx context.Context) error
}
