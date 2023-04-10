package http

import (
	"context"
	"net/http"
	"time"

	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/infrastructure/http/handler"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/infrastructure/http/middleware"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/infrastructure/http/rpc"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/infrastructure/http/rpc/websocket"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	connectedSockets = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "websocket_connection_count",
			Help: "Number of websocket connections alive",
		},
		[]string{},
	)
)

func init() {
	prometheus.MustRegister(connectedSockets)
}

func StartHttpServer(ctx context.Context, app *infinity.Application) {

	router := mux.NewRouter()

	injectHandlers(app, router)

	handler := handlers.RecoveryHandler(
		handlers.PrintRecoveryStack(true),
	)(router)

	handler = handlers.CORS(
		handlers.AllowedOrigins([]string{"http://dev.camtube.co:8080", "https://camtube.co"}),
		handlers.AllowedHeaders([]string{"Authorization", "Content-Type"}),
		handlers.AllowCredentials(),
		handlers.MaxAge(3600),
	)(handler)

	handler = handlers.ProxyHeaders(handler)
	handler = middleware.Jwt(app, handler)

	srv := http.Server{
		Addr:    app.HttpBindAddr,
		Handler: handler,
	}

	go func() {
		app.Logger.Infof("Starting http server: %s", app.HttpBindAddr)

		if err := srv.ListenAndServeTLS("/etc/ssl/server.crt", "/etc/ssl/server.key"); err != nil {
			app.Logger.WithError(err).Fatal("Failed to start http server")
		}
	}()

	<-ctx.Done()
	app.Logger.Warn("Main context.Done() called closing http server")

	shutdownCtx := context.Background()
	shutdownCtx, cancel := context.WithTimeout(shutdownCtx, time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		app.Logger.WithError(err).Errorf("Failed to shutdown http server cleanly")
	}
}

type route struct {
	endpoint string
	handler  http.HandlerFunc
	method   string
}

func injectHandlers(app *infinity.Application, router *mux.Router) {

	// user
	userHandler := handler.NewUserHandler(app)
	userRpcHandler := handler.NewUserRpcHandler(app)

	// content
	performerHandler := handler.NewPerformerHandler(app)
	recordingHandler := handler.NewRecordingHandler(app)
	recordingRpcHandler := handler.NewRecordingRpcHandler(app)

	// Payment plan
	paymentHandler := handler.NewPaymentsHandler(app)

	// misc
	sessionHandler := handler.NewSessionHandler(app)
	bandwidthHandler := handler.NewBandwidthHandler(app)

	routes := []route{
		// Misc
		{"/metrics", promhttp.Handler().ServeHTTP, http.MethodGet},

		// Session endpoints
		{"/rpc/session/new", sessionHandler.New, http.MethodGet},
		{"/rpc/session/renew", sessionHandler.Renew, http.MethodGet},
		{"/rpc/session/authenticate", sessionHandler.Authenticate, http.MethodPost},

		// recording
		{"/recordings", recordingHandler.GetAll, http.MethodGet},
		{"/recordings/{uuidOrSlug}", recordingHandler.GetById, http.MethodGet},
		{"/recordings/{uuid}/manifest.m3u8", recordingHandler.GetManifest, http.MethodGet},

		{"/rpc/recording.view", recordingRpcHandler.IncrementViewCount, http.MethodPost},
		{"/rpc/recording.like", recordingRpcHandler.ToggleLike, http.MethodPost},
		{"/rpc/recording.favorite", recordingRpcHandler.ToggleFavorite, http.MethodPost},

		// user
		{"/users", userHandler.Create, http.MethodPost},
		{"/users/{uuid}", userHandler.GetById, http.MethodGet},
		{"/users/{uuid}/favorites", recordingHandler.GetUserFavorites, http.MethodGet},

		{"/rpc/user.available", userRpcHandler.Available, http.MethodPost},
		{"/rpc/user.password-reset", userRpcHandler.PasswordReset, http.MethodPost},
		{"/rpc/user.new-password", userRpcHandler.NewPassword, http.MethodPost},

		// performer
		{"/performers", performerHandler.GetAll, http.MethodGet},
		{"/performers/{uuidOrSlug}", performerHandler.GetById, http.MethodGet},
		{"/performers/{uuidOrSlug}/recordings", recordingHandler.GetAllForPerformer, http.MethodGet},

		// payment
		{"/payment-plans", paymentHandler.GetAll, http.MethodGet},
		{"/payment-transactions/{uuid}", paymentHandler.GetPaymentTransaction, http.MethodGet},
		{"/blockcypher/webhook/{uuid}", paymentHandler.BlockCypherWebHook, http.MethodPost},
		{"/blockcypher/forwarding-callback/{uuid}", paymentHandler.BlockCypherForwardingCallback, http.MethodPost},
		{"/rpc/payments-plans.purchase", paymentHandler.InitiatePurchase, http.MethodPost},
	}

	for _, route := range routes {
		router.Handle(route.endpoint, middleware.RequestTracker(route.handler)).Methods(route.method)
	}

	onSocketDisconnected := func(socket rpc.Socket) {
		// track number of connected sockets
		connectedSockets.WithLabelValues().Dec()

		// if callback is set for socket, call it
		if socket.Value("onDisconnected") != nil {
			socket.Value("onDisconnected").(func(rpc.Socket))(socket)
		}
	}

	onSocketConnected := func(socket rpc.Socket) {
		// track number of connected sockets
		connectedSockets.WithLabelValues().Inc()
	}

	rpcMap := make(rpc.ProcedureMap)
	rpcMap["session:set"] = sessionHandler.SetSocketSession
	rpcMap["bandwidth:get-remaining"] = bandwidthHandler.GetRemainingBandwidth

	rpcServer := rpc.NewRpcServer(app.Logger, rpcMap)

	// setup rpc over websocket server
	websocketHandler := websocket.NewSocketHandler(app.Logger, rpcServer, onSocketConnected, onSocketDisconnected)

	// handler the websocket server
	router.Handle("/ws", websocketHandler).Methods(http.MethodGet)
}
