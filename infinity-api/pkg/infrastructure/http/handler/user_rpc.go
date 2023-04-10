package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"time"

	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/infrastructure/http/handler/internal"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type passwordResetJwtToken struct {
	jwt.StandardClaims
	UsernameOrEmail string `json:"usernameOrEmail"`
}

type userRpc struct {
	app *infinity.Application
	log logrus.FieldLogger
}

func NewUserRpcHandler(app *infinity.Application) *userRpc {
	return &userRpc{
		app: app,
		log: app.Logger.WithFields(logrus.Fields{
			"component": "http",
			"handler":   "userRpc",
		}),
	}
}

func (endpoint *userRpc) Available(w http.ResponseWriter, r *http.Request) {
	logger := endpoint.app.Logger.WithField("RequestId", r.Context().Value("RequestId"))

	check := &internal.UserDetailAvailable{}
	if err := json.NewDecoder(r.Body).Decode(check); err != nil {
		jsonResponse(w, http.StatusBadRequest, RESP_MALFORMED_DATA)
		return
	}

	var err error

	switch strings.ToLower(check.Field) {
	case "username":
		_, err = endpoint.app.UserRepository.GetByUsername(check.Value)

	case "email":
		_, err = endpoint.app.UserRepository.GetByEmail(check.Value)

	default:
		jsonResponse(w, http.StatusBadRequest, &ApiResponse{
			Code:    400,
			Message: "Unsupported availability check, expected username or email",
		})

		return
	}

	switch err {
	case infinity.UserNotFoundErr:
		w.WriteHeader(http.StatusOK)

	case nil:
		w.WriteHeader(http.StatusConflict)

	default:
		logger.
			WithFields(logrus.Fields{
				"field": check.Field,
				"value": check.Value,
			}).
			WithError(err).Error("Failed to check user availability")

		jsonResponse(w, http.StatusInternalServerError, RESP_SERVER_ERROR)
	}
}

func (endpoint *userRpc) PasswordReset(w http.ResponseWriter, r *http.Request) {
	logger := endpoint.app.Logger.WithField("RequestId", r.Context().Value("RequestId"))

	data := &internal.PasswordReset{}
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		jsonResponse(w, http.StatusBadRequest, RESP_MALFORMED_DATA)
		return
	}

	user, err := endpoint.app.UserRepository.GetByUsernameOrEmail(data.UsernameOrEmail)
	if err != nil {
		switch err {
		// don't wanna expose emails in the system so just simulate that it succeeded
		case infinity.UserNotFoundErr:
			w.WriteHeader(http.StatusNoContent)
			return
		default:
			jsonResponse(w, http.StatusInternalServerError, RESP_SERVER_ERROR)
			return
		}
	}

	token, err := jwt.NewWithClaims(endpoint.app.SessionSignMethod, passwordResetJwtToken{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "camtube.co",
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
		UsernameOrEmail: data.UsernameOrEmail,
	}).SignedString(endpoint.app.SessionSignKey)

	if err != nil {
		logger.WithError(err).Error("Failed to create jwt password reset token")
		jsonResponse(w, http.StatusInternalServerError, RESP_SERVER_ERROR)
		return
	}

	endpoint.app.EmailService.SendPasswordResetEmail(user, token)

	w.WriteHeader(http.StatusNoContent)
}

func (endpoint *userRpc) NewPassword(w http.ResponseWriter, r *http.Request) {
	logger := endpoint.app.Logger.WithField("RequestId", r.Context().Value("RequestId"))

	data := &internal.NewPassword{}
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		jsonResponse(w, http.StatusBadRequest, RESP_MALFORMED_DATA)
		return
	}

	claims := &passwordResetJwtToken{}

	_, err := jwt.ParseWithClaims(data.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(endpoint.app.SessionSignKey), nil
	})

	if err != nil {
		logger.WithError(err).Error("Failed to parse jwt token to update password")
		jsonResponse(w, http.StatusInternalServerError, RESP_SERVER_ERROR)
		return
	}

	user, err := endpoint.app.UserRepository.GetByUsernameOrEmail(claims.UsernameOrEmail)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, RESP_SERVER_ERROR)
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.WithError(err).Error("Failed to encrypt the provided password")
		jsonResponse(w, http.StatusInternalServerError, RESP_SERVER_ERROR)
		return
	}

	user.Password = string(password)
	if err = endpoint.app.UserRepository.Update(user); err != nil {
		logger.WithError(err).Error("Failed to set new password of user")
		jsonResponse(w, http.StatusInternalServerError, RESP_SERVER_ERROR)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
