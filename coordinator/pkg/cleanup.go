package coordinator

import (
	"context"

	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"github.com/pkg/errors"
)

var (
	RecordingHasInfinityViewsErr                = errors.New("The recording has views in infinity")
	RecordingInInfinityPremiumUserCollectionErr = errors.New("The recording belongs to a premium user's collection in infinity")
)

type CleanupService interface {
	Run(ctx context.Context)
	DispatchMp42HlsConversion() error
	TransferUpstoreHashToAphrodite() error
	RinsePublishingFailures(remove bool) error
	RinseEncodingFailures(remove bool) error
	RinseUploadingFailures(remove bool) error
	RinseMp4Recordings() error
	Delete(recording *ecosystem.Recording, system string) error
}

type DeathFileReport struct {
	MissingInHermes    *uint64
	MissingInAphrodite *uint64
	Found              *uint64
	Skipped            *uint64
	Deleted            *uint64
}

type DeathFileService interface {
	Process(ctx context.Context, routines int, source string) error
	GetReport() *DeathFileReport
}
