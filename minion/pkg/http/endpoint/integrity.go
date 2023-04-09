package endpoint

import (
	"encoding/json"
	"net/http"

	"git.misc.vee.bz/carnagel/minion/pkg"
)

type integrityHandler struct {
	app *minion.Application
}

func NewIntegrityHandler(app *minion.Application) *integrityHandler {
	return &integrityHandler{
		app: app,
	}
}

func (handler *integrityHandler) GetReport(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(struct {
		Running  bool                   `json:"running"`
		Scanning bool                   `json:"scanning"`
		Report   minion.IntegrityReport `json:"report"`
	}{
		handler.app.FilesystemIntegrity.IsRunning(),
		handler.app.FilesystemIntegrity.IsScanning(),
		handler.app.FilesystemIntegrity.GetReport(),
	})

	if err != nil {
		handler.app.Logger.WithError(err).Errorf("Failed to generate json")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
