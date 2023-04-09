package middleware

import (
	"crypto/subtle"
	"net/http"
)

type tokenMiddleware struct {
	expectedToken   string
	nextHandlerFunc http.HandlerFunc
}

func Authenticate(expectedToken string, handler http.HandlerFunc) http.Handler {
	return &tokenMiddleware{
		expectedToken:   expectedToken,
		nextHandlerFunc: handler,
	}
}

func (w *tokenMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Authorization")
	if token == "" {
		token = r.URL.Query().Get("token")
	}

	if subtle.ConstantTimeCompare([]byte(token), []byte(w.expectedToken)) == 1 {
		w.nextHandlerFunc(rw, r)
	} else {
		rw.WriteHeader(http.StatusUnauthorized)
	}
}
