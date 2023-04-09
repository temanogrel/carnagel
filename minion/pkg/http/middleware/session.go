package middleware

import (
	"context"
	"net/http"
	"strings"

	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"github.com/dgrijalva/jwt-go"
)

const (
	RequestCtxRole      = "role"
	RequestCtxUserId    = "user"
	RequestCtxSessionId = "session"
)

type sessionMiddleware struct {
	key         string
	nextHandler http.Handler
}

func Session(key string, handler http.Handler) http.Handler {
	return &sessionMiddleware{
		key:         key,
		nextHandler: handler,
	}
}

func (middleware *sessionMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	if strings.HasPrefix(r.URL.Path, "/metrics") {
		middleware.nextHandler.ServeHTTP(rw, r)
		return
	}

	session := r.URL.Query().Get("session")
	if session == "" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Missing session"))

		return
	}

	claims := &ecosystem.JwtClaims{}

	_, err := jwt.ParseWithClaims(session, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(middleware.key), nil
	})

	if err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Invalid session provided"))
		return
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, RequestCtxRole, claims.Role)
	ctx = context.WithValue(ctx, RequestCtxUserId, claims.User)
	ctx = context.WithValue(ctx, RequestCtxSessionId, claims.Session)

	middleware.nextHandler.ServeHTTP(rw, r.WithContext(ctx))
}
