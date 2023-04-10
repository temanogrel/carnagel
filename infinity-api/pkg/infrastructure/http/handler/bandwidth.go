package handler

import (
	"encoding/json"

	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/infrastructure/http/rpc"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type bandwidthHandler struct {
	app *infinity.Application
	log logrus.FieldLogger
}

func NewBandwidthHandler(app *infinity.Application) *bandwidthHandler {
	return &bandwidthHandler{
		app: app,
		log: app.Logger.WithFields(logrus.Fields{
			"component": "http",
			"handler":   "bandwidth",
		}),
	}
}

func (handler *bandwidthHandler) GetRemainingBandwidth(request *rpc.Request, socket rpc.Socket) (json.RawMessage, error) {
	session, ok := socket.Value("session").(uuid.UUID)
	if !ok {
		return nil, errors.New("Missing session identifier")
	}

	paymentPlan, ok := socket.Value("paymentPlan").(uuid.UUID)
	if !ok {
		return nil, errors.New("Missing payment plan identifier")
	}

	plan, err := handler.app.PaymentPlanRepository.GetByUuid(paymentPlan)
	if err != nil {
		return nil, errors.New("Failed to retrieve payment plan")
	}

	today, err := handler.app.BandwidthConsumptionCollector.GetTotalConsumptionToday(session)
	if err == infinity.BlackListedSessionErr {
		return rpc.ToRawMessage(&infinity.BandwidthStatus{
			Total:     plan.Bandwidth,
			Remaining: 0,
		})
	}

	remaining := int64(plan.Bandwidth) - int64(today)
	if remaining < 0 {
		remaining = 0
	}

	return rpc.ToRawMessage(&infinity.BandwidthStatus{
		Total:     plan.Bandwidth,
		Remaining: remaining,
	})
}

func (handler *bandwidthHandler) MayWatchVideo(request *rpc.Request, socket rpc.Socket) (json.RawMessage, error) {
	session, ok := socket.Value("session").(uuid.UUID)
	if !ok {
		return nil, errors.New("Missing session identifier")
	}

	body := struct {
		RecordingUuid uuid.UUID `json:"recordingUuid"`
	}{}

	if err := json.Unmarshal(request.Payload, &body); err != nil {
		return nil, errors.Wrapf(err, "Failed to decode request payload")
	}

	paymentPlan, ok := socket.Value("paymentPlan").(uuid.UUID)
	if !ok {
		return nil, errors.New("Missing payment plan identifier")
	}

	plan, err := handler.app.PaymentPlanRepository.GetByUuid(paymentPlan)
	if err != nil {
		return nil, errors.New("Failed to retrieve payment plan")
	}

	todayTotal, err := handler.app.BandwidthConsumptionCollector.GetTotalConsumptionToday(session)
	if todayTotal >= plan.Bandwidth || err == infinity.BlackListedSessionErr {
		return nil, errors.New("Exceeded daily usage")
	}

	todayOnRecording, err := handler.app.BandwidthConsumptionCollector.GetTotalConsumptionTodayOnRecording(session, body.RecordingUuid)

	if (plan.PerRecordingBandwidth != 0 && todayOnRecording >= plan.PerRecordingBandwidth) || err == infinity.BlackListedSessionErr {
		return nil, errors.New("Exceeded daily usage on this recording")
	}

	return nil, nil
}
