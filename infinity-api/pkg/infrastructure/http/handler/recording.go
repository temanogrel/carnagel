package handler

import (
	"net/http"
	"strconv"

	"strings"

	"fmt"

	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/infrastructure/http/middleware"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/infrastructure/http/serializers"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type recording struct {
	app *infinity.Application
	log logrus.FieldLogger
}

func NewRecordingHandler(app *infinity.Application) *recording {
	return &recording{
		app: app,
		log: app.Logger.WithFields(logrus.Fields{
			"component": "http",
			"handler":   "recording",
		}),
	}
}

func (endpoint *recording) GetById(w http.ResponseWriter, r *http.Request) {
	logger := endpoint.app.Logger.WithField("RequestId", r.Context().Value("requestId"))

	id := mux.Vars(r)["uuidOrSlug"]

	var err error
	var recording interface{}

	if userId, ok := r.Context().Value(middleware.RequestCtxUserId).(uuid.UUID); ok {
		if recordingUuid, e := uuid.FromString(id); e == nil {
			recording, err = endpoint.app.RecordingRepository.GetByUuidWithUserContext(recordingUuid, userId)
		} else {
			recording, err = endpoint.app.RecordingRepository.GetBySlugWithUserContext(id, userId)
		}
	} else {
		if recordingUuid, e := uuid.FromString(id); e == nil {
			recording, err = endpoint.app.RecordingRepository.GetByUuid(recordingUuid)
		} else {
			recording, err = endpoint.app.RecordingRepository.GetBySlug(id)
		}
	}

	switch err {
	case nil:
		jsonResponse(w, http.StatusOK, serializers.NewRecordingSerializer().Transform(recording))

	case infinity.RecordingNotFoundErr:
		jsonResponse(w, http.StatusNotFound, RESP_OBJECT_NOT_FOUND)

	default:
		logger.WithError(err).Error("Unexpected error occurred when retrieving recording")

		jsonResponse(w, http.StatusInternalServerError, RESP_SERVER_ERROR)
	}
}

func (endpoint *recording) GetManifest(w http.ResponseWriter, r *http.Request) {
	logger := endpoint.app.Logger.WithField("RequestId", r.Context().Value("requestId"))

	id, err := uuid.FromString(mux.Vars(r)["uuid"])
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, RESP_MALFORMED_ID)
		return
	}

	recording, err := endpoint.app.RecordingRepository.GetByUuid(id)
	switch err {
	case nil:

		// todo: we shouldn't already redirect to p1 but we only have one proxy at the moment so it does not matter.
		streamUrl := fmt.Sprintf("https://p1.camtube.co/%s", recording.VideoUuid.String())

		w.Write([]byte(strings.Replace(recording.VideoManifest, "stream.ts", streamUrl, -1)))

	case infinity.RecordingNotFoundErr:
		jsonResponse(w, http.StatusNotFound, RESP_OBJECT_NOT_FOUND)

	default:
		logger.WithError(err).Error("Unexpected error occurred when retrieving recording")

		jsonResponse(w, http.StatusInternalServerError, RESP_SERVER_ERROR)
	}
}

func (endpoint *recording) GetAllForPerformer(w http.ResponseWriter, r *http.Request) {
	logger := endpoint.app.Logger.WithField("RequestId", r.Context().Value("requestId"))

	id := mux.Vars(r)["uuidOrSlug"]

	var performer *infinity.Performer
	var err error

	if performerUuid, e := uuid.FromString(id); e == nil {
		performer, err = endpoint.app.PerformerRepository.GetByUuid(performerUuid)
	} else {
		performer, err = endpoint.app.PerformerRepository.GetBySlug(id)
	}

	if err != nil {
		jsonResponse(w, http.StatusNotFound, RESP_OBJECT_NOT_FOUND)
		return
	}

	criteria := &infinity.RecordingRepositoryCriteria{
		PerformerId: performer.Uuid,
		Limit:       30,
		Offset:      0,
	}

	if limit, err := strconv.Atoi(r.FormValue("limit")); err == nil && limit > 0 && limit <= 90 {
		criteria.Limit = limit
	}

	if offset, err := strconv.Atoi(r.FormValue("offset")); err == nil && offset >= 0 {
		criteria.Offset = offset
	}

	var recordings interface{}
	var total int

	if userId, ok := r.Context().Value(middleware.RequestCtxUserId).(uuid.UUID); ok {
		recordings, total, err = endpoint.app.RecordingRepository.MatchingWithUserContext(criteria, userId)
	} else {
		recordings, total, err = endpoint.app.RecordingRepository.Matching(criteria)
	}

	if err != nil {
		logger.
			WithField("recordingRepositoryCriteria", criteria).
			WithError(err).
			Error("Failed to retrieve recordings of performer")

		jsonResponse(w, http.StatusInternalServerError, RESP_SERVER_ERROR)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"data": serializers.NewRecordingSerializer().MustTransformArray(recordings),
		"meta": ResponseMeta{
			Total:  total,
			Offset: criteria.Offset,
			Limit:  criteria.Limit,
		},
	})
}

func (endpoint *recording) GetAll(w http.ResponseWriter, r *http.Request) {
	logger := endpoint.app.Logger.WithField("RequestId", r.Context().Value("requestId"))

	criteria := &infinity.RecordingRepositoryCriteria{
		Limit: 90,
	}

	givenOffSet := 0

	if limit, err := strconv.Atoi(r.FormValue("limit")); err == nil && limit > 0 && limit <= 90 {
		criteria.Limit = limit
	}

	if offset, err := strconv.Atoi(r.FormValue("offset")); err == nil && offset >= 0 {
		// Since inject page cache resets the offset of the criteria
		givenOffSet = offset
		criteria.Offset = offset
	}

	sort := r.FormValue("sortMode")
	if sort != "" {
		interval, err := strconv.Atoi(r.FormValue("interval"))
		if err != nil {
			interval = 7
		}

		criteria.SortMode = sort
		criteria.Interval = uint8(interval)
	}

	var recordings interface{}
	var total int
	var err error

	if userId, ok := r.Context().Value(middleware.RequestCtxUserId).(uuid.UUID); ok {
		recordings, total, err = endpoint.app.RecordingRepository.MatchingWithUserContext(criteria, userId)
	} else {
		recordings, total, err = endpoint.app.RecordingRepository.Matching(criteria)
	}

	if err != nil {
		logger.
			WithField("recordingRepositoryCriteria", criteria).
			WithError(err).
			Error("Failed to retrieve recordings")

		jsonResponse(w, http.StatusInternalServerError, RESP_SERVER_ERROR)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"data": serializers.NewRecordingSerializer().MustTransformArray(recordings),
		"meta": ResponseMeta{
			Total:  total,
			Offset: givenOffSet,
			Limit:  criteria.Limit,
		},
	})
}

func (endpoint *recording) GetUserFavorites(w http.ResponseWriter, r *http.Request) {
	logger := endpoint.app.Logger.WithField("RequestId", r.Context().Value("requestId"))

	id, err := uuid.FromString(mux.Vars(r)["uuid"])
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, RESP_MALFORMED_ID)
		return
	}

	criteria := &infinity.RecordingRepositoryCriteria{
		Limit:         30,
		FavoritesOnly: true,
	}

	if stageName := r.FormValue("stageName"); stageName != "" {
		criteria.StageName = stageName
	}

	if limit, err := strconv.Atoi(r.FormValue("limit")); err == nil && limit > 0 && limit <= 90 {
		criteria.Limit = limit
	}

	if offset, err := strconv.Atoi(r.FormValue("offset")); err == nil && offset >= 0 {
		criteria.Offset = offset
	}

	recordings, total, err := endpoint.app.RecordingRepository.MatchingWithUserContext(criteria, id)
	if err != nil {
		logger.
			WithFields(logrus.Fields{"recordingRepositoryCriteria": criteria, "userUuid": id}).
			WithError(err).
			Error("Failed to retrieve user favorite recordings")

		jsonResponse(w, http.StatusInternalServerError, RESP_SERVER_ERROR)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"data": serializers.NewRecordingSerializer().MustTransformArray(recordings),
		"meta": ResponseMeta{
			Total:  total,
			Offset: criteria.Offset,
			Limit:  criteria.Limit,
		},
	})
}
