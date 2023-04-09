package http

import (
	"net/http"

	"context"

	"git.misc.vee.bz/carnagel/encoder/pkg"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartServer(ctx context.Context, app *encoder.Application) {
	router := mux.NewRouter()
	router.Handle("/metrics", promhttp.Handler())

	srv := http.Server{
		Addr:    app.HttpBindAddr,
		Handler: router,
	}

	go func() {
		app.Logger.Infof("Starting http server: %s", app.HttpBindAddr)
		if err := srv.ListenAndServe(); err != nil {
			app.Logger.WithError(err).Fatal("Failed to start http server")
		}
	}()

	srv.Shutdown(ctx)
}
