package crypto

import (
	"context"
	"encoding/json"
	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
	"github.com/sasha-s/go-deadlock"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type cryptoExchangeRateService struct {
	app        *infinity.Application
	mtx        deadlock.RWMutex
	logger     logrus.FieldLogger
	httpClient *retryablehttp.Client
	rate       float64
}

func NewCryptoExchangeRateService(app *infinity.Application) infinity.CryptoExchangeRateService {
	return &cryptoExchangeRateService{
		app:        app,
		mtx:        deadlock.RWMutex{},
		httpClient: retryablehttp.NewClient(),
		logger:     app.Logger.WithField("service", "CryptoExchangeRateService"),
	}
}

func (s *cryptoExchangeRateService) Run(ctx context.Context) {
	logger := s.logger.WithField("operation", "Run")
	timer := time.NewTimer(0)

	for {
		select {
		case <-ctx.Done():
			timer.Stop()
			return

		case <-timer.C:
			func() {
				defer timer.Reset(time.Minute * 20)

				if err := s.updateExchangeRate(); err != nil {
					logger.WithError(err).Errorf("Failed to update exchange rate")
				}
			}()
		}
	}
}

func (s *cryptoExchangeRateService) updateExchangeRate() error {
	request, err := retryablehttp.NewRequest(http.MethodGet, "https://blockchain.info/ticker", nil)
	if err != nil {
		return errors.Wrapf(err, "Failed to create request")
	}

	response, err := s.httpClient.Do(request)
	if err != nil {
		return errors.Wrapf(err, "Failed to get exchange rate")
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.Errorf("Unexpected status code: %d", response.StatusCode)
	}

	data := map[string]struct {
		Last float64 `json:"last"`
	}{}

	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return errors.Wrapf(err, "Failed to decode exchange rate")
	}

	if _, ok := data["USD"]; !ok {
		return errors.Wrapf(err, "No exchange rate for USD found")
	}

	s.mtx.Lock()
	s.rate = data["USD"].Last
	s.mtx.Unlock()

	return nil
}

func (s *cryptoExchangeRateService) ConvertUSDToBtc(amount float64) (float64, float64, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	if s.rate == 0 {
		return 0, 0, errors.New("Current exchange rate unavailable")
	}

	return amount / s.rate, s.rate, nil
}
