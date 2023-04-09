package http

import (
	"context"
	"net/http"
	"time"

	"git.misc.vee.bz/carnagel/go-ecosystem/infrastructure/middleware"
	"git.misc.vee.bz/carnagel/minerva/pkg"
	"git.misc.vee.bz/carnagel/minerva/pkg/http/endpoint"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartServer(ctx context.Context, app *minerva.Application) {
	go startPublicServer(ctx, app)
	go startPrivateServer(ctx, app)
}

// startPrivateServer will spin up a http server with no tls/ssl listening on the internal network of our cloud
func startPrivateServer(ctx context.Context, app *minerva.Application) {

	relocator := endpoint.NewStorage(app)

	router := mux.NewRouter()
	router.Handle("/metrics", promhttp.Handler())
	router.Handle("/rpc.storage.relocate", middleware.ResponseMetrics(relocator.Relocate))
	router.Handle("/rpc.storage.relocate-single", middleware.ResponseMetrics(relocator.RelocateSingle))

	srv := http.Server{
		Addr:    app.HttpBindAddrPrivate,
		Handler: applyMiddleware(router),
	}

	go func() {
		app.Logger.Infof("Started private http server: %s", app.HttpBindAddrPrivate)

		if err := srv.ListenAndServe(); err != nil {
			app.Logger.WithError(err).Fatal("Failed to start private http server")
		}
	}()

	<-ctx.Done()
	app.Logger.Warn("Main context.Done() called closing private http server")

	shutdownCtx := context.Background()
	shutdownCtx, cancel := context.WithTimeout(shutdownCtx, time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		app.Logger.WithError(err).Errorf("Failed to shutdown private http server cleanly")
	}
}

// startPublicServer will spin a http server with tls/ssl so it can receive communication from the servers in amsterdam
func startPublicServer(ctx context.Context, app *minerva.Application) {
	locator := endpoint.NewLocator(app)

	router := mux.NewRouter()
	router.Handle("/locator", middleware.ResponseMetrics(locator.ServeHTTP))

	srv := http.Server{
		Addr:    app.HttpBindAddrPublic,
		Handler: applyMiddleware(router),
	}

	go func() {
		app.Logger.Infof("Started public http server: %s", app.HttpBindAddrPublic)

		if err := srv.ListenAndServeTLS("/etc/ssl/server.crt", "/etc/ssl/server.key"); err != nil {
			app.Logger.WithError(err).Fatal("Failed to start public http server")
		}
	}()

	<-ctx.Done()
	app.Logger.Warn("Main context.Done() called closing public http server")

	shutdownCtx := context.Background()
	shutdownCtx, cancel := context.WithTimeout(shutdownCtx, time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		app.Logger.WithError(err).Errorf("Failed to shutdown public http server cleanly")
	}
}

func applyMiddleware(router *mux.Router) http.Handler {

	wrappedHandler := handlers.RecoveryHandler(
		handlers.PrintRecoveryStack(true),
	)(router)

	wrappedHandler = middleware.RequestId(wrappedHandler.ServeHTTP)
	wrappedHandler = handlers.ProxyHeaders(wrappedHandler)
	wrappedHandler = handlers.CORS(
		handlers.AllowCredentials(),
		handlers.AllowedOrigins([]string{"camtube.co"}),
		handlers.AllowedHeaders([]string{"Authorization", "Content-Type", "Content-Length", "Range"}),
	)(wrappedHandler)

	return wrappedHandler
}
