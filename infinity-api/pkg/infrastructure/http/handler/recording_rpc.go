package handler

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/infrastructure/http/handler/internal"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/infrastructure/http/middleware"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type recordingRpc struct {
	app *infinity.Application
	log logrus.FieldLogger
}

func NewRecordingRpcHandler(app *infinity.Application) *recordingRpc {
	return &recordingRpc{
		app: app,
		log: app.Logger.WithFields(logrus.Fields{
			"component": "http",
			"handler":   "recordingRpc",
		}),
	}
}

func (endpoint *recordingRpc) IncrementViewCount(w http.ResponseWriter, r *http.Request) {
	logger := endpoint.log.WithField("RequestId", r.Context().Value("RequestId"))

	body := &internal.Uuid{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		jsonResponse(w, http.StatusBadRequest, RESP_MALFORMED_DATA)
		return
	}

	recording, err := endpoint.app.RecordingRepository.GetByUuid(body.Uuid)
	if err != nil {
		if err == infinity.RecordingNotFoundErr {
			jsonResponse(w, http.StatusNotFound, RESP_OBJECT_NOT_FOUND)
			return
		}

		logger.WithError(err).Error("Failed to retrieve recording for incrementing view count")

		jsonResponse(w, http.StatusInternalServerError, RESP_SERVER_ERROR)
		return
	}

	userHash := r.RemoteAddr + r.Header.Get("User-Agent")

	hasher := sha1.New()
	hasher.Write([]byte(userHash))

	userHash = hex.EncodeToString(hasher.Sum(nil))

	added, err := endpoint.app.RecordingRepository.AddView(recording, userHash)
	if err != nil {
		logger.WithError(err).Error("Failed to increment the view count")

		jsonResponse(w, http.StatusInternalServerError, RESP_SERVER_ERROR)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]bool{
		"added": added,
	})
}

func (endpoint *recordingRpc) ToggleLike(w http.ResponseWriter, r *http.Request) {
	logger := endpoint.log.WithField("RequestId", r.Context().Value("RequestId"))

	body := &internal.Uuid{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		jsonResponse(w, http.StatusBadRequest, RESP_MALFORMED_DATA)
		return
	}

	recording, err := endpoint.app.RecordingRepository.GetByUuid(body.Uuid)
	if err != nil {
		if err == infinity.RecordingNotFoundErr {
			jsonResponse(w, http.StatusNotFound, RESP_OBJECT_NOT_FOUND)
			return
		}

		logger.WithError(err).Error("Failed to retrieve recording for toggle like")

		jsonResponse(w, http.StatusInternalServerError, RESP_SERVER_ERROR)
		return
	}

	userId, ok := r.Context().Value(middleware.RequestCtxUserId).(uuid.UUID)
	if !ok {
		logger.WithField("userId", userId).Error("Failed to cast userId to uuid")

		jsonResponse(w, http.StatusInternalServerError, RESP_SERVER_ERROR)
		return
	}

	present, err := endpoint.app.UserRepository.ToggleLike(userId, recording.Uuid)
	if err != nil {
		logger.WithError(err).Error("Failed to toggle like for recording")

		jsonResponse(w, http.StatusInternalServerError, RESP_SERVER_ERROR)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]bool{
		"present": present,
	})
}

func (endpoint *recordingRpc) ToggleFavorite(w http.ResponseWriter, r *http.Request) {
	logger := endpoint.log.WithField("RequestId", r.Context().Value("RequestId"))

	body := &internal.Uuid{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		jsonResponse(w, http.StatusBadRequest, RESP_MALFORMED_DATA)
		return
	}

	recording, err := endpoint.app.RecordingRepository.GetByUuid(body.Uuid)
	if err != nil {
		if err == infinity.RecordingNotFoundErr {
			jsonResponse(w, http.StatusNotFound, RESP_OBJECT_NOT_FOUND)
			return
		}

		logger.WithError(err).Error("Failed to retrieve recording for ToggleFavorites")

		jsonResponse(w, http.StatusInternalServerError, RESP_SERVER_ERROR)
		return
	}

	userId, ok := r.Context().Value(middleware.RequestCtxUserId).(uuid.UUID)
	if !ok {
		logger.WithField("userId", userId).Error("Failed to cast userId to uuid")

		jsonResponse(w, http.StatusInternalServerError, RESP_SERVER_ERROR)
		return
	}

	present, err := endpoint.app.UserRepository.ToggleFavorite(userId, recording.Uuid)
	if err != nil {
		logger.WithError(err).
			WithFields(logrus.Fields{
				"userId":      userId,
				"recordingId": recording.Uuid,
			}).
			Error("Failed to toggle favorites for user")

		jsonResponse(w, http.StatusInternalServerError, RESP_SERVER_ERROR)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]bool{
		"present": present,
	})
}
