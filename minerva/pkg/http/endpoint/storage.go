package endpoint

import (
	"net/http"

	"encoding/json"

	"git.misc.vee.bz/carnagel/minerva/pkg"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type relocateRequest struct {
	Top    int    `json:"top"`
	Bottom int    `json:"bottom"`
	Amount uint64 `json:"amount"`
}

type relocateSingleRequest struct {
	Uuid   uuid.UUID        `json:"uuid"`
	Origin minerva.Hostname `json:"origin"`
	Target minerva.Hostname `json:"target"`
}

type storage struct {
	app *minerva.Application
	log logrus.FieldLogger
}

func NewStorage(app *minerva.Application) *storage {
	return &storage{
		app: app,
		log: app.Logger.WithField("component", "http"),
	}
}

func (endpoint *storage) Relocate(w http.ResponseWriter, r *http.Request) {

	req := &relocateRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.
		NewEncoder(w).
		Encode(endpoint.app.LoadBalancer.RedistributeData(req.Top, req.Bottom, req.Amount))
}

func (endpoint *storage) RelocateSingle(w http.ResponseWriter, r *http.Request) {

	req := &relocateSingleRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	collection := endpoint.app.ServerCollection
	collection.ServersMtx.RLock()
	server, ok := collection.Servers[req.Origin]
	collection.ServersMtx.RUnlock()

	if ok {
		server.RelocateRequests <- minerva.RelocateRequest{
			Uuid:       req.Uuid,
			TargetHost: req.Target,
		}

		w.WriteHeader(http.StatusNoContent)
	}

	w.WriteHeader(http.StatusBadRequest)
}
