package handler

import (
	"net/http"

	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	"github.com/sirupsen/logrus"
)

type siteMapHandler struct {
	app *infinity.Application
	log logrus.FieldLogger
}

func NewSiteMapHandler(app *infinity.Application) *siteMapHandler {
	return &siteMapHandler{
		app: app,
		log: app.Logger.WithField("handler", "sitemap"),
	}
}

func (handler *siteMapHandler) Index(w http.ResponseWriter, r *http.Request) {

}
