package endpoint

import (
	"context"
	"net/http"
	"net/http/httputil"
	"strings"

	"strconv"

	"os"

	"git.misc.vee.bz/carnagel/go-ecosystem/protobuf/common"
	"git.misc.vee.bz/carnagel/minion/pkg"
	"git.misc.vee.bz/carnagel/minion/pkg/http/endpoint/internal"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type consumerHandler struct {
	app *minion.Application
	log logrus.FieldLogger
}

func NewConsumerHandler(app *minion.Application) http.Handler {
	return &consumerHandler{
		app: app,
		log: app.Logger.WithFields(logrus.Fields{
			"component": "http",
			"handler":   "infinity",
		}),
	}
}

func (handler *consumerHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	log := handler.app.Logger.WithField("requestId", r.Context().Value("requestId"))

	// process query params
	r.ParseForm()

	path := r.FormValue("path")
	address := r.FormValue("address")
	rangeHeader := r.Header.Get("Range")
	session := handler.extractSession(r)

	recordingId, err := strconv.ParseUint(r.FormValue("recording-id"), 10, 64)
	if err != nil {
		log.WithError(err).Warn("Failed to parse recording-id")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	trackBandwidth, err := strconv.ParseBool(r.FormValue("bandwidth-tracking"))
	if err != nil {
		log.WithError(err).Warn("Failed to parse bandwidth-tracking")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if trackBandwidth && rangeHeader == "" {
		rw.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}

	log = log.WithFields(logrus.Fields{
		"download": map[string]interface{}{
			"path":        path,
			"address":     address,
			"session":     session,
			"recordingId": recordingId,
		},
	})

	isProxy, err := handler.isProxyRequest(address, path)
	if err != nil {
		log.WithError(err).Warn("Failed to serve file")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if isProxy {
		log.Debug("Proxying remote file")

		proxy := httputil.ReverseProxy{
			Director: func(req *http.Request) {
				req.URL.Host = address
				req.URL.Scheme = "http"

				log.
					WithField("url", req.URL.String()).
					Debug("Passing on via proxy")
			},
		}

		proxy.ServeHTTP(rw, r)
		return
	}

	if trackBandwidth && !handler.hasAvailableBandwidth(rw, r, log, session) {
		return
	}

	log.Debug("Serving local file")

	http.ServeFile(rw, r, path)

	if trackBandwidth {
		stat, err := os.Stat(path)
		if err != nil {
			log.
				WithError(err).
				Error("Failed to get file size")

			return
		}

		ranges, err := internal.ParseRange(rangeHeader, stat.Size())
		if err != nil {
			log.
				WithError(err).
				Error("Failed to parse ranges")

			return
		}

		var bytes uint64

		for _, r := range ranges {
			bytes += uint64(r.Length)
		}

		handler.sendBandwidthUsage(log, session, bytes, recordingId)
	}
}

// extractSession attempts to get the session from the cookie or defaulting to the request.FormValue
func (handler *consumerHandler) extractSession(r *http.Request) string {
	session := r.FormValue("session")

	if cookie, err := r.Cookie("session"); err == nil {
		session = cookie.Value
	}

	return session
}

// isProxyRequest checks if we should proxy the current request, or if it's to be served by the local instance
func (handler *consumerHandler) isProxyRequest(address, path string) (bool, error) {

	parts := strings.Split(address, ":")
	if len(parts) != 2 {
		return false, errors.New("Invalid address supplied hostname:port is the expected format")
	}

	// if target host ends with current hostname, should be fine
	if strings.HasSuffix(parts[0], handler.app.Hostname) {
		if !strings.HasPrefix(path, handler.app.DataDir) {
			return false, errors.New("Attempt to access files outside of the data directory")
		}

		return false, nil
	}

	return true, nil
}

// hasAvailableBandwidth checks if the current session has
func (handler *consumerHandler) hasAvailableBandwidth(rw http.ResponseWriter, r *http.Request, log logrus.FieldLogger, session string) bool {

	tokenRequest := &common.Token{
		Token: session,
	}

	remainingBandwidth, err := handler.app.BandwidthTrackingClient.GetRemainingBandwidth(r.Context(), tokenRequest)
	if err != nil {
		log.WithError(err).Error("Failed to retrieve remaining bandwidth")
		rw.WriteHeader(http.StatusInternalServerError)

		return false
	}

	if remainingBandwidth.Status != common.StatusCode_Ok {
		log.
			WithField("status", remainingBandwidth.Status).
			Error("Unexpected response code from infinity")

		rw.WriteHeader(http.StatusInternalServerError)
		return false
	}

	if remainingBandwidth.Remaining == 0 {
		rw.WriteHeader(http.StatusPaymentRequired)
		return false
	}

	return true
}

func (handler *consumerHandler) sendBandwidthUsage(log logrus.FieldLogger, session string, bytes, recordingId uint64) {

	log.
		WithField("bytes", bytes).
		Debug("Tracking bandwidth")

	resp, err := handler.app.BandwidthTrackingClient.AddConsumedBandwidth(context.TODO(), &common.ConsumedBandwidthRequest{
		Token:       session,
		RecordingId: recordingId,
		Amount:      bytes,
	})

	if err != nil {
		log.
			WithError(err).
			Error("Failed to add consumed bandwidth")

		return
	}

	if resp.Status != common.StatusCode_Ok {
		log.
			WithField("status", resp.Status).
			Error("Unexpected response code on add consumed bandwidth")

		return
	}
}
