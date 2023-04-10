package server

import (
	"git.misc.vee.bz/carnagel/go-ecosystem/protobuf/common"
	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	"github.com/sasha-s/go-deadlock"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type deviceTrackingServer struct {
	app *infinity.Application
	log logrus.FieldLogger
	mtx deadlock.RWMutex

	devicesPerSession map[uuid.UUID]uint8
}

func NewDeviceTrackingServer(app *infinity.Application) common.DeviceTrackingServer {
	return &deviceTrackingServer{
		app: app,
		log: app.Logger.WithFields(logrus.Fields{
			"component":   "grpc",
			"grpc_server": "device_tracking",
		}),
		mtx: deadlock.RWMutex{},
	}
}

func (server *deviceTrackingServer) IncDeviceCount(ctx context.Context, req *common.Token) (*common.IncDeviceCountResponse, error) {
	claims, err := server.app.UserSessionService.ParseToken(req.Token)
	if err != nil {
		return &common.IncDeviceCountResponse{
			Status: common.StatusCode_InvalidSessionToken,
		}, nil
	}

	plan, err := server.app.PaymentPlanRepository.GetByUuid(claims.PaymentPlan)
	if err != nil {
		server.log.
			WithError(err).
			WithField("claims", claims).
			Errorf("Failed to retrieve payment plan")

		return &common.IncDeviceCountResponse{
			Status: common.StatusCode_InternalServerErr,
		}, nil
	}

	server.mtx.Lock()
	defer server.mtx.Unlock()

	devices, _ := server.devicesPerSession[claims.Session]

	if devices >= plan.Devices {
		return &common.IncDeviceCountResponse{
			Status: common.StatusCode_TooManyActiveDevices,
		}, nil
	}

	devices++

	server.devicesPerSession[claims.Session] = devices

	return &common.IncDeviceCountResponse{
		Status: common.StatusCode_Ok,
	}, nil
}

func (server *deviceTrackingServer) DecDeviceCount(ctx context.Context, req *common.Token) (*common.DecDeviceCountResponse, error) {
	claims, err := server.app.UserSessionService.ParseToken(req.Token)
	if err != nil {
		return &common.DecDeviceCountResponse{
			Status: common.StatusCode_InvalidSessionToken,
		}, nil
	}

	server.mtx.Lock()
	defer server.mtx.Unlock()

	if _, ok := server.devicesPerSession[claims.Session]; ok {
		server.devicesPerSession[claims.Session] -= 1
	}

	return &common.DecDeviceCountResponse{
		Status: common.StatusCode_Ok,
	}, nil
}
