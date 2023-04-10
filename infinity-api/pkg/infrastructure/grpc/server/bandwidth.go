package server

import (
	"git.misc.vee.bz/carnagel/go-ecosystem/protobuf/common"
	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type bandwidthTrackingServer struct {
	app *infinity.Application
	log logrus.FieldLogger
}

func NewBandwidthTrackingServer(app *infinity.Application) common.BandwidthTrackingServer {
	return &bandwidthTrackingServer{
		app: app,
		log: app.Logger.WithFields(logrus.Fields{
			"component":   "grpc",
			"grpc_server": "bandwidth_tracking",
		}),
	}
}

func (server *bandwidthTrackingServer) GetRemainingBandwidth(ctx context.Context, req *common.Token) (*common.RemainingBandwidthResponse, error) {
	log := server.log.WithField("request", req)

	claims, err := server.app.UserSessionService.ParseToken(req.Token)
	if err != nil {
		return &common.RemainingBandwidthResponse{
			Status: common.StatusCode_InvalidSessionToken,
		}, nil
	}

	log = log.WithField("claims", claims)

	plan, err := server.app.PaymentPlanRepository.GetByUuid(claims.PaymentPlan)
	if err != nil {
		log.
			WithError(err).
			Errorf("Failed to retrieve payment plan")

		return &common.RemainingBandwidthResponse{
			Status: common.StatusCode_PlanNotFound,
		}, nil
	}

	usedToday, err := server.app.BandwidthConsumptionCollector.GetTotalConsumptionToday(claims.Session)
	if err == infinity.BlackListedSessionErr {
		return &common.RemainingBandwidthResponse{
			Status:    common.StatusCode_Ok,
			Remaining: 0,
		}, nil
	}

	var remaining uint64
	if plan.Bandwidth > usedToday {
		remaining = plan.Bandwidth - usedToday
	}

	return &common.RemainingBandwidthResponse{
		Status:    common.StatusCode_Ok,
		Remaining: remaining,
	}, nil
}

func (server *bandwidthTrackingServer) AddConsumedBandwidth(ctx context.Context, req *common.ConsumedBandwidthRequest) (*common.RemainingBandwidthResponse, error) {
	log := server.log.WithField("request", req)
	log.Debug("Received AddConsumedBandwidth request")

	claims, err := server.app.UserSessionService.ParseToken(req.Token)
	if err != nil {
		log.Warn("Invalid token provided")

		return &common.RemainingBandwidthResponse{
			Status: common.StatusCode_InvalidSessionToken,
		}, nil
	}

	log = log.WithField("claims", claims)

	plan, err := server.app.PaymentPlanRepository.GetByUuid(claims.PaymentPlan)
	if err != nil {
		log.
			WithError(err).
			Errorf("Failed to retrieve payment plan")

		return &common.RemainingBandwidthResponse{
			Status: common.StatusCode_PlanNotFound,
		}, nil
	}

	recordingUuid, err := server.app.RecordingRepository.GetUuidByExternalId(req.RecordingId)
	if err != nil {
		log.
			WithError(err).
			Errorf("Failed to retrieve recording uuid from external id")

		return &common.RemainingBandwidthResponse{
			Status: common.StatusCode_RecordingNotFound,
		}, nil
	}

	consumption := &infinity.BandwidthConsumption{
		RecordingUuid: recordingUuid,
		SessionUuid:   claims.Session,
		Bytes:         req.Amount,
	}

	if claims.User != uuid.Nil {
		consumption.UserUuid = uuid.NullUUID{UUID: claims.User, Valid: true}
	}

	total, onRecording := server.app.BandwidthConsumptionCollector.AddConsumption(consumption)

	var remaining uint64
	var remainingFromTotal uint64
	var remainingOnRecording uint64

	if total <= plan.Bandwidth {
		remainingFromTotal = plan.Bandwidth - total
	}

	if plan.PerRecordingBandwidth != 0 && onRecording <= plan.PerRecordingBandwidth {
		remainingOnRecording = plan.PerRecordingBandwidth - onRecording
	}

	if remainingFromTotal < remainingOnRecording {
		remaining = remainingFromTotal
	} else {
		remaining = remainingOnRecording
	}

	return &common.RemainingBandwidthResponse{
		Status:    common.StatusCode_Ok,
		Remaining: remaining,
	}, nil
}
