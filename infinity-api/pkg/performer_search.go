package infinity

import (
	"context"
	"github.com/satori/go.uuid"
)

type PerformerSearchService interface {
	Run(ctx context.Context)

	AddPerformerAlias(performerId uuid.UUID, stageName string) error
	PerformerAliasExists(performerId uuid.UUID, stageName string) (bool, error)

	Matching(criteria *PerformerRepositoryCriteria) ([]uuid.UUID, int, error)
}
