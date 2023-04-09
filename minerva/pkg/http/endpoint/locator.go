package endpoint

import (
	"fmt"
	"net/http"
	"net/url"

	"strconv"

	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"git.misc.vee.bz/carnagel/minerva/pkg"
	"github.com/sirupsen/logrus"
)

type locator struct {
	app *minerva.Application
	log logrus.FieldLogger
}

func NewLocator(app *minerva.Application) *locator {
	return &locator{
		app: app,
		log: app.Logger.WithField("component", "http"),
	}
}

func (endpoint *locator) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Query().Get("uuid")
	session := endpoint.extractSession(r)
	hostname := r.URL.Query().Get("hostname")

	logger := endpoint.log.WithFields(logrus.Fields{
		"method": "locator",
		"query": map[string]string{
			"uuid":     uuid,
			"session":  session,
			"hostname": hostname,
		},

		"uri": r.RequestURI,
	})

	if session == "" {
		logger.Debug("Locator failed, missing session")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if uuid == "" {
		logger.Debug("Locator failed, missing uuid")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if hostname == "" {
		logger.Debug("Locator failed, missing hostname")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	path, err := endpoint.app.LoadBalancer.RecommendDownload(minerva.Hostname(hostname), uuid)
	if err != nil {
		if err == minerva.FileNotFoundErr {
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		logger.WithError(err).Debug("Failed to locate a file")

		// it's important we don't expose any information to the user
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	u := url.URL{}
	u.Scheme = "http"
	u.Path = "/consume"

	if path.Edge != nil {
		u.Host = fmt.Sprintf("%s:%d", path.Edge.ExternalHostname, endpoint.app.MinionPort)
	} else {
		u.Host = fmt.Sprintf("%s:%d", path.Origin.ExternalHostname, endpoint.app.MinionPort)
	}

	query := u.Query()
	query.Set("path", path.File.Path)
	query.Set("token", endpoint.app.MinionToken)
	query.Set("session", session)
	query.Set("recording-id", strconv.FormatUint(uint64(path.File.ExternalId), 10))
	query.Set("bandwidth-tracking", strconv.FormatBool(path.File.Type == ecosystem.FileTypeRecordingHls))

	if path.InterfaceToUse == "internal" && path.Edge != nil {
		query.Set("address", fmt.Sprintf("%s:%d", path.Origin.InternalHostname, endpoint.app.MinionPort))
	} else {
		query.Set("address", fmt.Sprintf("%s:%d", path.Origin.ExternalHostname, endpoint.app.MinionPort))
	}

	u.RawQuery = query.Encode()

	http.Redirect(rw, r, u.String(), http.StatusFound)

	logger.
		WithField("u", u.String()).
		Debug("Locator found file and proposed a location")
}

// extractSession attempts to get the session from the cookie or defaulting to the request.FormValue
func (endpoint *locator) extractSession(r *http.Request) string {
	session := r.FormValue("session")

	if cookie, err := r.Cookie("session"); err == nil {
		session = cookie.Value
	}

	return session
}
