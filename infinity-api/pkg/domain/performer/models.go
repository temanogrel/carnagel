package performer

import "github.com/satori/go.uuid"

type elasticPerformerDocument struct {
	PerformerId uuid.UUID `json:"performerId"`
	StageName   string    `json:"stageName"`
}
