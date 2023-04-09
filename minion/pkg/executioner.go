package minion

import "context"

type ExecutionerService interface {

	// ProcessDeletions receives deletion requests from the DelegationMinionClient and processes them
	ProcessDeletions(ctx context.Context)

	// ProcessUploads receives upload requests from the DelegationMinionClient and processes them
	ProcessUploads(ctx context.Context)

	// ProcessRelocationRequests receive a relocation request from the DelegationMinionClient and process them
	ProcessRelocationRequests(ctx context.Context)
}
