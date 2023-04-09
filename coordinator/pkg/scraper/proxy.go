package scraper

import (
	"context"
	"fmt"
	"git.misc.vee.bz/carnagel/coordinator/pkg"
	"github.com/hashicorp/consul/api"
	"github.com/sasha-s/go-deadlock"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type scraperProxyService struct {
	app     *coordinator.Application
	mtx     deadlock.RWMutex
	logger  logrus.FieldLogger
	proxies []string
}

func NewScraperProxyService(app *coordinator.Application) coordinator.ScraperProxyService {
	return &scraperProxyService{
		app:     app,
		mtx:     deadlock.RWMutex{},
		logger:  app.Logger.WithField("service", "ScraperProxyService"),
		proxies: strings.Split(app.Consul.GetString("proxy-servers/ips", ""), ","),
	}
}

func (s *scraperProxyService) Run(ctx context.Context) {
	timer := time.NewTimer(0)

	for {
		select {
		case <-ctx.Done():
			timer.Stop()
			return

		case <-timer.C:
			s.logger.Info("Retrieving new proxy servers")
			s.syncProxyServers()
			timer.Reset(time.Minute * 30)
		}
	}
}

func (s *scraperProxyService) syncProxyServers() {
	proxies, err := s.app.ProxyServerClient.GetProxyIps()
	if err != nil {
		return
	}

	go s.updateProxiesInConsul(proxies)

	s.mtx.Lock()
	s.proxies = proxies
	s.mtx.Unlock()
}

func (s *scraperProxyService) updateProxiesInConsul(proxies []string) {
	logger := s.logger.WithField("method", "updateProxiesInConsul")
	logger.Debug("Updating proxies in consul")

	pair := &api.KVPair{
		Key:   "proxy-servers/ips",
		Value: []byte(strings.Join(proxies, ",")),
	}

	if _, err := s.app.Consul.API().KV().Put(pair, &api.WriteOptions{}); err != nil {
		logger.WithError(err).Errorf("Failed to update KV pair in consul")
	}

	logger.Info("Successfully updated proxies in consul")
}

func (s *scraperProxyService) GetProxy(r *http.Request) (*url.URL, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	if len(s.proxies) == 0 {
		return nil, coordinator.NoProxiesAvailableErr
	}

	proxy := s.proxies[rand.Intn(len(s.proxies))]

	return url.Parse(fmt.Sprintf("http://%s", proxy))
}
