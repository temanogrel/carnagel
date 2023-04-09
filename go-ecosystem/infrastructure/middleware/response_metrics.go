package middleware

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

var endpointDurations = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Name:       "http_endpoint_duration_seconds",
		Help:       "Duration spent on each endpoint",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	},
	[]string{"endpoint", "method"},
)

type requestTracker struct {
	nextHandler http.HandlerFunc
}

func init() {
	prometheus.MustRegister(endpointDurations)
}

func ResponseMetrics(handler http.HandlerFunc) http.Handler {
	return &requestTracker{nextHandler: handler}
}

func (middleware *requestTracker) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	startedAt := time.Now()

	middleware.nextHandler(rw, r)

	// This will never have any long time duration
	if r.Method == http.MethodOptions {
		return
	}

	route := mux.CurrentRoute(r)
	if route == nil {
		log.Warnf("No route matched '%s' for request tracking", r.RequestURI)
		return
	}

	template, err := route.GetPathTemplate()
	if err != nil {
		log.Error(err)
		return
	}

	endpointDurations.
		WithLabelValues(template, r.Method).
		Observe(time.Since(startedAt).Seconds())
}
