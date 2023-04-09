package middleware

import (
	"context"
	"net/http"

	"github.com/satori/go.uuid"
)

const RequestCtxRequestId = "requestId"

type requestId struct {
	nextHandler http.HandlerFunc
}

func RequestId(handler http.HandlerFunc) http.Handler {
	return &requestId{nextHandler: handler}
}

func (middleware *requestId) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	requestId := r.Header.Get("X-Request-Id")
	if requestId == "" {
		requestId = uuid.NewV4().String()
	}

	// return the request id for future tracking
	rw.Header().Set("X-Request-Id", requestId)

	ctx := r.Context()
	ctx = context.WithValue(ctx, RequestCtxRequestId, requestId)

	middleware.nextHandler(rw, r.WithContext(ctx))
}
