package handler

import (
	"encoding/json"
	"net/http"

	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/infrastructure/http/rpc"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type authenticationRequest struct {
	Identity         string `json:"identity"`
	Password         string `json:"password"`
	RememberMe       bool   `json:"rememberMe"`
	EndOtherSessions bool   `json:"endOtherSessions"`
}

type appAuthRequest struct {
	Token    string `json:"token"`
	Hostname string `json:"hostname"`
}

type session struct {
	app *infinity.Application
	log logrus.FieldLogger
}

func NewSessionHandler(app *infinity.Application) *session {
	return &session{
		app: app,
		log: app.Logger.WithFields(logrus.Fields{
			"component": "http",
			"handler":   "sessions",
		}),
	}
}

func (handler *session) New(rw http.ResponseWriter, r *http.Request) {
	token, err := handler.app.UserSessionService.New(r.RemoteAddr)
	if err != nil {
		handler.log.
			WithError(err).
			Error("Failed to generate a new session token")

		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonResponse(rw, http.StatusOK, struct {
		Token string `json:"token"`
	}{token})
}

func (handler *session) Renew(rw http.ResponseWriter, r *http.Request) {
	handler.log.Debug("Received session renew request")

	// jwt middleware does not handle /rpc/session routes
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

	claims, err := handler.app.UserSessionService.ParseToken(cookie.Value)
	if err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Session cookie is invalid or has expired"))
		return
	}

	log := handler.log.WithField("userUuid", claims.User)
	log.Debug("Attempting to renew token for user")

	token, err := handler.app.UserSessionService.Renew(claims.User, claims.Session, claims.BlackListedToday)
	if err != nil {
		switch err {
		case infinity.UserNotFoundErr:
			log.WithError(err).Warn("Guest tried to renew session")
			jsonResponse(rw, http.StatusBadRequest, RESP_MALFORMED_ID)
		default:
			log.WithError(err).Error("Failed to renew session token")
			jsonResponse(rw, http.StatusInternalServerError, RESP_SERVER_ERROR)
			return
		}
	}

	jsonResponse(rw, http.StatusOK, struct {
		Token string `json:"token"`
	}{token})
}

func (handler *session) Authenticate(rw http.ResponseWriter, r *http.Request) {

	payload := &authenticationRequest{}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		handler.log.
			WithError(err).
			Warn("Request with invalid json attempted")

		jsonResponse(rw, http.StatusBadRequest, RESP_MALFORMED_DATA)
		return
	}

	token, user, err := handler.app.UserSessionService.Authenticate(payload.Identity, payload.Password, r.RemoteAddr)
	switch err {
	case nil:
		if handler.userHasActiveSocketSession(user.Uuid) && !payload.EndOtherSessions {
			jsonResponse(rw, http.StatusBadRequest, &ApiResponse{
				Code:    400,
				Message: "User already has an active session",
				Errors: map[string]string{
					"userHasActiveSocketSession": "User already has an active session",
				},
			})
			return
		}

		jsonResponse(rw, http.StatusOK, struct {
			Token string `json:"token"`
		}{token})

	case infinity.InvalidCredentialsErr:
		jsonResponse(rw, http.StatusUnauthorized, RESP_UNAUTHORIZED)

	default:
		handler.log.
			WithError(err).
			Error("Unexpected error during authentication")

		jsonResponse(rw, http.StatusInternalServerError, RESP_SERVER_ERROR)
	}
}

func (handler *session) SetSocketSession(request *rpc.Request, socket rpc.Socket) (json.RawMessage, error) {

	body := struct {
		Token string `json:"token"`
	}{}

	if err := json.Unmarshal(request.Payload, &body); err != nil {
		return nil, errors.Wrap(err, "Failed to decode request body")
	}

	claims := &ecosystem.JwtClaims{}

	_, err := jwt.ParseWithClaims(body.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(handler.app.SessionSignKey), nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "Failed to decode token")
	}

	// End the current active socket session
	handler.endUserActiveSocketSession(socket)

	// Set a reference to the user
	socket.SetValue("socket", uuid.NewV4())
	socket.SetValue("user", claims.User)
	socket.SetValue("role", claims.Role)
	socket.SetValue("session", claims.Session)
	socket.SetValue("paymentPlan", claims.PaymentPlan)
	socket.SetValue("onDisconnected", handler.endUserActiveSocketSession)

	// Add the new session
	handler.setActiveSocketSession(socket)

	return nil, nil
}

func (handler *session) userHasActiveSocketSession(userUuid uuid.UUID) bool {
	handler.app.ActiveSessionCollection.Mtx.RLock()
	defer handler.app.ActiveSessionCollection.Mtx.RUnlock()

	_, ok := handler.app.ActiveSessionCollection.Sessions[userUuid]

	return ok
}

func (handler *session) endUserActiveSocketSession(socket rpc.Socket) {
	log := handler.log.WithField("operation", "endUserActiveSocketSession")

	// If we never had a signed in socket session, not much we can do
	if socket.Value("user") != nil && socket.Value("user").(uuid.UUID) != uuid.Nil {
		userUuid := socket.Value("user").(uuid.UUID)
		log = log.WithField("userUuid", userUuid)
		log.Debug("Removing active socket from user")

		handler.app.ActiveSessionCollection.Mtx.Lock()
		delete(handler.app.ActiveSessionCollection.Sessions, userUuid)
		handler.app.ActiveSessionCollection.Mtx.Unlock()
	}
}

func (handler *session) setActiveSocketSession(socket rpc.Socket) {
	log := handler.log.WithField("operation", "setActiveSocketSession")

	// If we don't have a signed in session, not much to do
	if socket.Value("user") != nil && socket.Value("user").(uuid.UUID) != uuid.Nil {
		userUuid := socket.Value("user").(uuid.UUID)
		log = log.WithField("userUuid", userUuid)

		handler.app.ActiveSessionCollection.Mtx.RLock()
		if s, ok := handler.app.ActiveSessionCollection.Sessions[userUuid]; ok {
			handler.app.ActiveSessionCollection.Mtx.RUnlock()
			if (*s).Value("socket") != socket.Value("socket") {
				log.WithField("socket", (*s).Value("socket")).Debug("Socket session ended")
				(*s).SetValue("session", nil)
				(*s).Send(rpc.CreateBroadcastResponse("session:ended", nil))
			}
		} else {
			handler.app.ActiveSessionCollection.Mtx.RUnlock()
		}

		log.WithField("socket", socket.Value("socket")).Debug("Setting active socket")
		handler.app.ActiveSessionCollection.Mtx.Lock()
		handler.app.ActiveSessionCollection.Sessions[userUuid] = &socket
		handler.app.ActiveSessionCollection.Mtx.Unlock()
	}
}
