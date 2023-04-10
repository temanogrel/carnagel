package handler

import (
	"github.com/blockcypher/gobcy"
	"net/http"

	"encoding/json"

	"time"

	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/infrastructure/http/handler/internal"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/infrastructure/http/middleware"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type payments struct {
	app *infinity.Application
	log logrus.FieldLogger
}

func NewPaymentsHandler(app *infinity.Application) *payments {
	return &payments{
		app: app,
		log: app.Logger.WithFields(logrus.Fields{
			"component": "http",
			"handler":   "payments",
		}),
	}
}

func (handler *payments) GetAll(rw http.ResponseWriter, r *http.Request) {
	logger := handler.log.WithField("RequestId", r.Context().Value("RequestId"))

	plans, err := handler.app.PaymentPlanRepository.GetAll()
	if err != nil {
		logger.
			WithError(err).
			Error("Failed to retrieve payment plans")

		jsonResponse(rw, http.StatusInternalServerError, RESP_SERVER_ERROR)
		return
	}

	jsonResponse(rw, http.StatusOK, map[string]interface{}{
		"data": plans,
		"meta": ResponseMeta{
			Total:  len(plans),
			Limit:  len(plans),
			Offset: 0,
		},
	})
}

func (handler *payments) InitiatePurchase(w http.ResponseWriter, r *http.Request) {
	userUuid, ok := r.Context().Value(middleware.RequestCtxUserId).(uuid.UUID)
	if !ok {
		handler.app.Logger.Warn("Unauthorized access")

		jsonResponse(w, http.StatusUnauthorized, RESP_UNAUTHORIZED)
		return
	}

	body := &internal.Uuid{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		handler.app.Logger.
			WithError(err).
			Warn("Failed to parse request body")

		jsonResponse(w, http.StatusBadRequest, RESP_MALFORMED_DATA)
		return
	}

	user, err := handler.app.UserRepository.GetByUuid(userUuid)
	if err != nil {
		handler.app.Logger.
			WithError(err).
			Error("Failed to retrieve the user")

		jsonResponse(w, http.StatusNotFound, RESP_OBJECT_NOT_FOUND)
		return
	}

	paymentPlan, err := handler.app.PaymentPlanRepository.GetByUuid(body.Uuid)
	switch err {
	case nil:
		if response, err := handler.app.PaymentService.InitiatePurchase(user, paymentPlan); err == nil {
			jsonResponse(w, http.StatusAccepted, response)
		} else {
			handler.app.Logger.
				WithError(err).
				Error("Failed to initiate purchase for unknown reason")

			jsonResponse(w, http.StatusInternalServerError, RESP_SERVER_ERROR)
		}

	case infinity.PaymentPlayNotFoundErr:
		jsonResponse(w, http.StatusNotFound, RESP_OBJECT_NOT_FOUND)

	default:
		handler.app.Logger.
			WithError(err).
			Error("Failed to retrieve payment plan")

		jsonResponse(w, http.StatusInternalServerError, RESP_SERVER_ERROR)
	}
}

func (handler *payments) BlockCypherForwardingCallback(w http.ResponseWriter, r *http.Request) {
	webhookUuid := uuid.FromStringOrNil(mux.Vars(r)["uuid"])

	logger := handler.log.WithField("endpoint", "BlockCypherForwardingCallback").WithField("webhookUuid", webhookUuid)
	logger.Info("BlockCypher forwarding callback triggered")

	transaction, err := handler.app.PaymentTransactionRepository.GetByWebhookUuid(webhookUuid)
	if err != nil {
		logger.WithError(err).Errorf("Failed to retrieve payment transaction by webhook uuid")
		jsonResponse(w, http.StatusNotFound, RESP_OBJECT_NOT_FOUND)
		return
	}

	data := &gobcy.Payback{}
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		logger.WithError(err).Errorf("Failed to decode payload")
		jsonResponse(w, http.StatusInternalServerError, RESP_SERVER_ERROR)
		return
	}

	if err := handler.app.PaymentService.AddressForwardingCallback(transaction, data); err != nil {
		logger.WithError(err).Errorf("Failed to process transaction")
		jsonResponse(w, http.StatusInternalServerError, RESP_SERVER_ERROR)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (handler *payments) BlockCypherWebHook(w http.ResponseWriter, r *http.Request) {
	webhookUuid := uuid.FromStringOrNil(mux.Vars(r)["uuid"])

	logger := handler.log.WithField("endpoint", "BlockCypherWebHook").WithField("webhookUuid", webhookUuid)
	logger.Info("BlockCypher webhook triggered")

	transaction, err := handler.app.PaymentTransactionRepository.GetByWebhookUuid(webhookUuid)
	if err != nil {
		logger.WithError(err).Errorf("Failed to retrieve payment transaction by webhook uuid")
		jsonResponse(w, http.StatusNotFound, RESP_OBJECT_NOT_FOUND)
		return
	}

	logger = logger.WithField("transactionUuid", transaction.Uuid).WithField("paymentAddress", transaction.PaymentAddress)

	go func() {
		for i := 0; i < 5; i++ {
			switch err := handler.app.PaymentService.UpdatePurchase(transaction); err {
			case infinity.NoPaymentReceivedErr:
				logger.Debug("No payment received on payment address, sleeping and attempting again in 30 seconds")
				time.Sleep(time.Second * 30)

			default:
				if err != nil {
					logger.WithError(err).Errorf("Failed to update payment transaction")
				}

				return
			}
		}
	}()

	w.WriteHeader(http.StatusNoContent)
}

func (handler *payments) GetPaymentTransaction(w http.ResponseWriter, r *http.Request) {
	userUuid, ok := r.Context().Value(middleware.RequestCtxUserId).(uuid.UUID)
	if !ok {
		handler.app.Logger.Warn("Unauthorized access")

		jsonResponse(w, http.StatusUnauthorized, RESP_UNAUTHORIZED)
		return
	}

	if _, err := handler.app.UserRepository.GetByUuid(userUuid); err != nil {
		handler.app.Logger.
			WithError(err).
			Error("Failed to retrieve the user")

		jsonResponse(w, http.StatusNotFound, RESP_OBJECT_NOT_FOUND)
		return
	}

	idFromRoute := mux.Vars(r)["uuid"]
	transactionUuid, err := uuid.FromString(idFromRoute)
	if err != nil {
		handler.app.Logger.WithError(err).Error("Failed to convert received payment transaction id to uuid")

		jsonResponse(w, http.StatusNotFound, RESP_OBJECT_NOT_FOUND)
		return
	}

	transaction, err := handler.app.PaymentTransactionRepository.GetByUuid(transactionUuid)
	if err != nil {
		handler.app.Logger.
			WithError(err).
			Error("Failed to retrieve payment transaction")

		jsonResponse(w, http.StatusNotFound, RESP_OBJECT_NOT_FOUND)
		return
	}

	// Check if current user owns the transaction
	if transaction.UserUuid != userUuid {
		handler.app.Logger.
			WithFields(logrus.Fields{"transactionUserUuid": transaction.UserUuid, "currentUserUuid": userUuid}).
			Debug("Current user is not the owner of this payment transaction")

		jsonResponse(w, http.StatusForbidden, RESP_FORBIDDEN)
		return
	}

	jsonResponse(w, http.StatusOK, transaction)
}
