package endpoint

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"strconv"

	"strings"

	"path/filepath"

	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"git.misc.vee.bz/carnagel/minerva-bindings/src"
	"git.misc.vee.bz/carnagel/minion/pkg"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type storageHandler struct {
	app *minion.Application
}

func NewStorageHandler(application *minion.Application) *storageHandler {
	return &storageHandler{app: application}
}

func (handler *storageHandler) Download(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// process query params
	r.ParseForm()

	path := r.FormValue("path")
	address := r.FormValue("address")

	logger := ecosystem.NewLogEntry(r.Context(), handler.app.Logger, logrus.Fields{
		"download": map[string]string{
			"path":    path,
			"address": address,
		},

		"component": "http",
		"handler":   "storage",
	})

	parts := strings.Split(address, ":")
	if len(parts) != 2 {
		logger.Warn("Invalid address specified")

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// if target host ends with current hostname, should be fine
	if strings.HasSuffix(parts[0], handler.app.Hostname) {
		logger.Debug("Serving local file")

		if !strings.HasPrefix(path, handler.app.DataDir) {
			logger.WithField("dir", handler.app.DataDir).Error("Trying to access file outside of data directory")

			w.WriteHeader(http.StatusNotFound)
			return
		}

		http.ServeFile(w, r, path)
		return
	}

	logger.Debug("Proxying file")
	proxy := httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Host = address
			req.URL.Scheme = "http"

			logger.
				WithField("url", req.URL.String()).
				Debug("Passing on via proxy")
		},
	}

	proxy.ServeHTTP(w, r)
}

func (handler *storageHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(handler.app.UploadFormMaxMemory)

	log := ecosystem.NewLogEntry(r.Context(), handler.app.Logger, logrus.Fields{
		"uuid":        r.FormValue("uuid"),
		"recordingId": r.FormValue("externalId"),
		"handler":     "storage",
		"component":   "http",
	})

	log.Debug("Received file transfer")

	id, err := uuid.FromString(r.FormValue("uuid"))
	if err != nil {
		log.
			WithError(err).
			Warn("Received transfer request with invalid uuid")

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	file, fileHandler, err := r.FormFile("file")
	if err != nil {
		log.
			WithError(err).
			Error("Failed to open uploaded file")

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer file.Close()

	if err := handler.app.FileService.HandleTransfer(r.Context(), id, file, filepath.Ext(fileHandler.Filename)); err != nil {
		log.
			WithError(err).
			Error("Failed to handle file transfer")

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (handler *storageHandler) Upload(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(handler.app.UploadFormMaxMemory)

	logger := ecosystem.NewLogEntry(r.Context(), handler.app.Logger, logrus.Fields{
		"upload": map[string]string{
			"externalId": r.FormValue("externalId"),
			"fileType":   r.FormValue("fileType"),
			"fileMeta":   r.FormValue("fileMeta"),
		},
		"recordingId": r.FormValue("externalId"),
		"component":   "http",
		"handler":     "storage",
	})

	externalId, err := strconv.ParseInt(r.FormValue("externalId"), 10, 64)
	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		logger.WithError(err).Warn("Failed to convert external id")
		return
	}

	fileType, err := strconv.Atoi(r.FormValue("fileType"))
	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		logger.WithError(err).Warn("Failed to convert external id")
		return
	}

	file, fileHandler, err := r.FormFile("file")
	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		logger.WithError(err).Error("Failed to open uploaded file")
		return
	}

	defer file.Close()

	uuid, err := handler.app.FileService.HandleUpload(file, fileHandler.Filename, uint64(externalId), pb.FileType(fileType))
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		logger.WithError(err).Error("Failed to handle uploaded file")
		return
	}

	w.WriteHeader(http.StatusCreated)

	// Return the minerva UUID
	json.NewEncoder(w).Encode(map[string]string{
		"uuid": uuid.String(),
	})

	logger.Debug("Uploaded new file")
}
