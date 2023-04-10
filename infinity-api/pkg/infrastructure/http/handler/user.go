package handler

import (
	"encoding/json"
	"net/http"

	"time"

	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/infrastructure/http/handler/internal"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/infrastructure/http/middleware"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type user struct {
	app *infinity.Application
	log logrus.FieldLogger
}

func NewUserHandler(app *infinity.Application) *user {
	return &user{
		app: app,
		log: app.Logger.WithFields(logrus.Fields{
			"component": "http",
			"handler":   "user",
		}),
	}
}

func (endpoint *user) GetById(w http.ResponseWriter, r *http.Request) {
	logger := endpoint.log.WithField("requestId", r.Context().Value("RequestId"))

	id, err := uuid.FromString(mux.Vars(r)["uuid"])
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, RESP_MALFORMED_ID)
		return
	}

	// Only authenticated users can retrieve profiles
	currentUserId, ok := r.Context().Value("user").(uuid.UUID)
	if !ok {
		jsonUnauthorizedResponse(w)
		return
	}

	currentUserRole, ok := r.Context().Value(middleware.RequestCtxRole).(string)
	if !ok {
		jsonUnauthorizedResponse(w)
		return
	}

	// Only admins can look up other users
	if currentUserId != id && currentUserRole != infinity.RoleAdmin {
		logger.WithFields(logrus.Fields{
			"currentUserId": currentUserId,
			"targetUserId":  id,
		}).Warn("Unauthorized access to user denied.")

		jsonForbiddenResponse(w)
		return
	}

	user, err := endpoint.app.UserRepository.GetByUuid(id)
	switch err {
	case nil:
		jsonResponse(w, http.StatusOK, user)

	case infinity.UserNotFoundErr:
		jsonResponse(w, http.StatusNotFound, RESP_OBJECT_NOT_FOUND)

	default:
		logger.WithError(err).Error("Failed to retrieve user by id")

		jsonResponse(w, http.StatusInternalServerError, RESP_SERVER_ERROR)
	}
}

func (endpoint *user) Create(w http.ResponseWriter, r *http.Request) {
	logger := endpoint.log.WithField("requestId", r.Context().Value("RequestId"))

	body := &internal.CreateUser{}

	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		jsonResponse(w, http.StatusBadRequest, RESP_MALFORMED_DATA)
		return
	}

	if _, err := endpoint.app.UserRepository.GetByUsername(body.Username); err == nil {
		jsonResponse(w, http.StatusConflict, &ApiResponse{
			Code:    409,
			Message: "Username is occupied",
			Error:   err,
		})

		return
	}

	if _, err := endpoint.app.UserRepository.GetByEmail(body.Email); err == nil {
		jsonResponse(w, http.StatusConflict, &ApiResponse{
			Code:    409,
			Message: "Username is occupied",
			Error:   err,
		})

		return
	}

	user := &infinity.User{
		Email:     body.Email,
		Username:  body.Username,
		Password:  body.Password,
		Role:      infinity.RoleUser,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := endpoint.app.UserRepository.Create(user); err != nil {
		logger.WithError(err).Error("Failed to retrieve user by id")

		jsonResponse(w, http.StatusInternalServerError, RESP_SERVER_ERROR)
		return
	}

	// Send welcome email
	endpoint.app.EmailService.SendAccountCreationEmail(user)

	jsonResponse(w, http.StatusCreated, user)
}

func (endpoint *user) Remove(w http.ResponseWriter, r *http.Request) {
	logger := endpoint.log.WithField("requestId", r.Context().Value("RequestId"))

	id, err := uuid.FromString(context.Get(r, "uuid").(string))
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, RESP_MALFORMED_ID)
		return
	}

	if _, err := endpoint.app.UserRepository.RemoveById(id); err != nil {
		logger.WithError(err).Error("Failed to retrieve user by id")

		jsonResponse(w, http.StatusInternalServerError, RESP_SERVER_ERROR)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
