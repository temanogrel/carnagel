package handler

import (
	"net/http"
	"strconv"

	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/infrastructure/http/serializers"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type performer struct {
	app *infinity.Application
	log logrus.FieldLogger
}

func NewPerformerHandler(app *infinity.Application) *performer {
	return &performer{
		app: app,
		log: app.Logger.WithFields(logrus.Fields{
			"component": "http",
			"handler":   "performer",
		}),
	}
}

func (endpoint *performer) GetAll(w http.ResponseWriter, r *http.Request) {
	logger := endpoint.log.WithField("RequestId", r.Context().Value("RequestId"))

	criteria := &infinity.PerformerRepositoryCriteria{
		Limit: 10,
	}

	if query := r.FormValue("query"); query != "" {
		criteria.StageName = query
	}

	if limit, err := strconv.Atoi(r.FormValue("limit")); err == nil && limit > 0 && limit <= 30 {
		criteria.Limit = limit
	}

	if offset, err := strconv.Atoi(r.FormValue("offset")); err == nil {
		criteria.Offset = offset
	}

	performerIds, total, err := endpoint.app.PerformerSearchService.Matching(criteria)
	if err != nil {
		logger.
			WithField("criteria", criteria).
			WithError(err).
			Error("Failed to query the performer search service")

		jsonResponse(w, http.StatusInternalServerError, RESP_SERVER_ERROR)
		return
	}

	dbCriteria := &infinity.PerformerRepositoryCriteria{
		PerformerIds: performerIds,
	}

	if includeLatestRecording := r.FormValue("includeLatestRecording"); includeLatestRecording == "1" {
		dbCriteria.IncludeLatestRecording = true
	}

	var performers []*infinity.Performer

	// Only query the database if actually got some results, removing this will cause it to retrieve the entire db
	if total > 0 {
		performers, _, err = endpoint.app.PerformerRepository.Matching(dbCriteria)
		if err != nil {
			logger.
				WithField("criteria", criteria).
				WithError(err).
				Error("Failed to retrieve performers")

			jsonResponse(w, http.StatusInternalServerError, RESP_SERVER_ERROR)
			return
		}
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"data": serializers.NewPerformerSerializer().MustTransformArray(performers),
		"meta": ResponseMeta{
			Total:  total,
			Limit:  criteria.Limit,
			Offset: criteria.Offset,
		},
	})
}

func (endpoint *performer) GetById(w http.ResponseWriter, r *http.Request) {
	logger := endpoint.log.WithField("RequestId", r.Context().Value("RequestId"))

	id := mux.Vars(r)["uuidOrSlug"]

	var performer *infinity.Performer
	var err error

	if performerUuid, e := uuid.FromString(id); e == nil {
		performer, err = endpoint.app.PerformerRepository.GetByUuid(performerUuid)
	} else {
		performer, err = endpoint.app.PerformerRepository.GetBySlug(id)
	}

	switch err {
	case nil:
		jsonResponse(w, http.StatusOK, serializers.NewPerformerSerializer().Transform(performer))

	case infinity.PerformerNotFoundErr:
		jsonResponse(w, http.StatusNotFound, RESP_OBJECT_NOT_FOUND)

	default:
		logger.WithError(err).Error("Unexpected error occurred when retrieving performer")

		jsonResponse(w, http.StatusInternalServerError, RESP_SERVER_ERROR)
	}
}
