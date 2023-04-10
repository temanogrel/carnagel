package infinity

import "context"

const (
	BlockcypherForwardingDestination = "3BsHWFknuVtzm2JBaei9PwvZqByJC55qhU"
)

type CryptoExchangeRateService interface {
	Run(ctx context.Context)
	ConvertUSDToBtc(amount float64) (float64, float64, error)
}
