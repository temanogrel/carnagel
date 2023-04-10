package middleware

import (
	"context"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"net/http"
	"strings"

	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	"github.com/satori/go.uuid"
)

const (
	RequestCtxRole          = "role"
	RequestCtxUserId        = "user"
	RequestCtxSessionId     = "session"
	RequestCtxPaymentPlanId = "paymentPlan"
)

type authorization struct {
	app         *infinity.Application
	nextHandler http.Handler
}

func Jwt(app *infinity.Application, handler http.Handler) http.Handler {
	return &authorization{
		app:         app,
		nextHandler: handler,
	}
}

func (middleware *authorization) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	// Don't attempt to authenticate option requests
	if r.Method == http.MethodOptions {
		middleware.nextHandler.ServeHTTP(rw, r)
		return
	}

	// rpc session stuff does not care about middleware and neither does the websocket service
	if strings.HasPrefix(r.URL.Path, "/rpc/session") || r.URL.Path == "/ws" {
		middleware.nextHandler.ServeHTTP(rw, r)
		return
	}

	// Handle crypto callbacks
	if strings.HasPrefix(r.URL.Path, "/blockcypher/webhook") {
		middleware.nextHandler.ServeHTTP(rw, r)
		return
	} else if strings.HasPrefix(r.URL.Path, "/blockcypher/forwarding-callback") {
		middleware.nextHandler.ServeHTTP(rw, r)
		return
	}

	// todo: add middleware here
	if strings.HasPrefix(r.URL.Path, "/metrics") {
		middleware.nextHandler.ServeHTTP(rw, r)
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Problem parsing the session cookie"))
		return
	}

	if cookie.Value == "" {
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Session cookie missing"))
		return
	}

	claims, err := middleware.app.UserSessionService.ParseToken(cookie.Value)
	if err != nil {
		if err == ecosystem.SessionExpiredErr {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Session cookie has expired"))
			return
		}

		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(err.Error()))
		return
	}

	ctx := r.Context()

	ctx = context.WithValue(ctx, RequestCtxRole, claims.Role)
	ctx = context.WithValue(ctx, "token", cookie.Value)

	if claims.User != uuid.Nil {
		ctx = context.WithValue(ctx, RequestCtxUserId, claims.User)
	}

	if claims.Session != uuid.Nil {
		ctx = context.WithValue(ctx, RequestCtxSessionId, claims.Session)
	}

	if claims.PaymentPlan != uuid.Nil {
		ctx = context.WithValue(ctx, RequestCtxPaymentPlanId, claims.PaymentPlan)
	}

	middleware.nextHandler.ServeHTTP(rw, r.WithContext(ctx))
}
