package middleware

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type authenticationToken struct {
	logger logrus.FieldLogger
}

func AuthenticationToken(logger logrus.FieldLogger) *authenticationToken {
	return &authenticationToken{
		logger: logger.WithFields(logrus.Fields{
			"middleware": "authentication_token",
			"component":  "http",
		}),
	}
}

func (middleware *authenticationToken) Wrap(next http.Handler, token string) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			next.ServeHTTP(rw, r)

			return
		}

		log := middleware.logger.WithFields(logrus.Fields{
			"url":        r.URL.String(),
			"remoteAddr": r.RemoteAddr,
		})

		if requestId, ok := r.Context().Value("RequestId").(string); ok {
			log = log.WithField("RequestId", requestId)
		}

		header := r.Header.Get("Authorization")
		if header == "" {
			log.Warn("Unauthorized access attempted")

			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		parts := strings.Split(header, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			log.Warn("Unauthorized access with invalid bearer")

			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		if subtle.ConstantTimeCompare([]byte(token), []byte(parts[1])) != 1 {
			log.Warn("Unauthorized access attempted")

			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(rw, r)
	})
}
