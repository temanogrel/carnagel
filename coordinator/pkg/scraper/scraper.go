package scraper

import (
	"context"
	"time"

	"fmt"

	"encoding/json"

	"bytes"

	"net/http"

	"strconv"

	"git.misc.vee.bz/carnagel/coordinator/pkg"
	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"github.com/sethgrid/pester"
	"github.com/sirupsen/logrus"
)

type scraper struct {
	app *coordinator.Application
	log logrus.FieldLogger
}

func NewScraper(app *coordinator.Application) coordinator.Scraper {
	return &scraper{
		app: app,
		log: app.Logger.WithField("component", "scraper"),
	}
}

func (scraper *scraper) Scrape(ctx context.Context) error {
	go scraper.app.MyfreeCamsScraper.Scrape(ctx)

	timer := time.NewTimer(time.Second * 10)

	for {
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil

		case <-timer.C:
			startedAt := time.Now()

			scraper.sendSessionId()
			scraper.sendPerformers()

			scraper.log.
				WithField("duration", time.Since(startedAt).Seconds()).
				Debug("Processed scrape content")

			timer.Reset(time.Second * 10)
		}
	}
}

func (scraper *scraper) getModelServerApi() (string, error) {

	services, _, err := scraper.app.Consul.API().Catalog().Service("modelserver", "", &api.QueryOptions{})
	if err != nil {
		return "", errors.Wrap(err, "Failed to retrieve modelserver from consul")
	}

	if len(services) == 0 {
		return "", errors.New("No model server available")
	}

	return fmt.Sprintf("http://%s:%d", services[0].ServiceAddress, services[0].ServicePort), nil
}

func (scraper *scraper) sendSessionId() {
	api, err := scraper.getModelServerApi()
	if err != nil {
		scraper.log.WithError(err).Error("Failed to send session id")
		return
	}

	payload := struct {
		SessionId int64 `json:"session_id"`
	}{scraper.app.MyfreeCamsScraper.GetSessionId()}

	requestBody := &bytes.Buffer{}

	if err := json.NewEncoder(requestBody).Encode(payload); err != nil {
		scraper.log.WithError(err).Error("Failed to send session id")
		return
	}

	request, err := http.NewRequest(http.MethodPost, api+"/mfc/session_id", requestBody)
	if err != nil {
		scraper.log.WithError(err).Error("Failed to send session id")
		return
	}

	response, err := pester.Do(request)
	if err != nil {
		scraper.log.WithError(err).Error("Failed to send session id")
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		scraper.log.WithField("statusCode", response.StatusCode).Warn("Unexpected response code")
	}
}

func (scraper *scraper) sendPerformers() {
	api, err := scraper.getModelServerApi()
	if err != nil {
		scraper.log.WithError(err).Error("Failed to send performers")
		return
	}

	startedAt := time.Now()

	performers, total := scraper.app.MyfreeCamsScraper.GetPerformers()

	payload := []modelServerMfcPerformer{}

	for _, p := range performers {
		payload = append(payload, modelServerMfcPerformer{
			ServiceId:      strconv.FormatInt(p.UserId, 10),
			StageName:      p.StageName,
			CurrentViewers: p.ViewerCount,
			VideoState:     p.VideoState,
			AccessLevel:    p.AccessLevel,
			CamServer:      p.User.Camserv,
			CamScore:       p.Meta.Camscore,
			MissMfcRank:    p.Meta.Missmfc,
		})
	}

	requestBody := &bytes.Buffer{}

	if err := json.NewEncoder(requestBody).Encode(payload); err != nil {
		scraper.log.WithError(err).Error("Failed to send session id")
		return
	}

	request, err := http.NewRequest(http.MethodPost, api+"/mfc/models/_intersect", requestBody)
	if err != nil {
		scraper.log.WithError(err).Error("Failed to send session id")
		return
	}

	response, err := pester.Do(request)
	if err != nil {
		scraper.log.WithError(err).Error("Failed to send session id")
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		scraper.log.
			WithField("statusCode", response.StatusCode).
			Warn("Unexpected response code")
	}

	scraper.log.
		WithFields(logrus.Fields{
			"duration": time.Since(startedAt).Seconds(),
			"total":    total,
			"sent":     len(payload),
		}).
		Debug("Sent mfc performers to model server")
}
