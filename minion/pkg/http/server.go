package http

import (
	"context"
	"fmt"
	"net/http"
	"os"

	ecosystem_middleware "git.misc.vee.bz/carnagel/go-ecosystem/infrastructure/middleware"
	"git.misc.vee.bz/carnagel/minion/pkg"
	"git.misc.vee.bz/carnagel/minion/pkg/http/endpoint"
	"git.misc.vee.bz/carnagel/minion/pkg/http/middleware"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartServer(ctx context.Context, app *minion.Application) {

	// Storage handler is the endpoint used internally
	storageHandler := endpoint.NewStorageHandler(app)

	// Consumer handler is the endpoint used by consumers of content and has bandwidth monitoring builtin
	consumerHandler := endpoint.NewConsumerHandler(app)

	// Expose the progress of the integrity handler
	integrityHandler := endpoint.NewIntegrityHandler(app)

	jwtSignKey := app.Consul.GetString("infinity/jwt-sign-key", "helloWorld")

	router := mux.NewRouter()
	router.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)
	router.Handle("/consume", ecosystem_middleware.ResponseMetrics(middleware.Session(jwtSignKey, consumerHandler).ServeHTTP))

	// Storage handler
	router.Handle("/upload", ecosystem_middleware.ResponseMetrics(middleware.Authenticate(app.HttpWriteToken, storageHandler.Upload).ServeHTTP))
	router.Handle("/transfer", ecosystem_middleware.ResponseMetrics(middleware.Authenticate(app.HttpWriteToken, storageHandler.Transfer).ServeHTTP))
	router.Handle("/download", ecosystem_middleware.ResponseMetrics(middleware.Authenticate(app.HttpReadToken, storageHandler.Download).ServeHTTP)).Methods(http.MethodGet, http.MethodHead)
	router.Handle("/integrity-report", ecosystem_middleware.ResponseMetrics(middleware.Authenticate(app.HttpReadToken, integrityHandler.GetReport).ServeHTTP))

	// Start the http server
	bindAddr := fmt.Sprintf("0.0.0.0:%s", os.Getenv("NOMAD_PORT_http"))

	wrappedHandler := handlers.RecoveryHandler(
		handlers.PrintRecoveryStack(true),
	)(router)

	wrappedHandler = ecosystem_middleware.RequestId(wrappedHandler.ServeHTTP)
	wrappedHandler = handlers.ProxyHeaders(wrappedHandler)
	wrappedHandler = handlers.CORS(
		handlers.AllowCredentials(),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{"Authorization", "Content-Type", "Content-Length", "Range"}),
	)(wrappedHandler)

	srv := http.Server{
		Addr:    bindAddr,
		Handler: wrappedHandler,
	}

	go func() {
		app.Logger.Infof("Starting http server: %s", bindAddr)

		if err := srv.ListenAndServe(); err != nil {
			app.Logger.WithError(err).Fatal("Failed to start http server")
		}
	}()

	srv.Shutdown(ctx)
}
