package server

import (
	"git.misc.vee.bz/carnagel/go-ecosystem/protobuf/common"
	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type contentServer struct {
	app *infinity.Application
	log logrus.FieldLogger
}

func NewContentServer(app *infinity.Application) common.ContentServer {
	return &contentServer{
		app: app,
		log: app.Logger.WithFields(logrus.Fields{
			"component": "grpc",
			"server":    "content",
		}),
	}
}

func (server *contentServer) UpsertRecording(ctx context.Context, req *common.RecordingIdentifier) (*common.UpsertRecordingResponse, error) {
	log := server.log.WithField("recordingId", req.Id)
	log.Debug("Received upsert recording")

	recording, performer, err := server.app.RecordingService.ImportRecording(ctx, req.Id)
	if err != nil {
		log.WithError(err).Error("Failed to import recording")

		return &common.UpsertRecordingResponse{
			Status: common.StatusCode_InternalServerErr,
		}, nil
	}

	log.WithFields(logrus.Fields{
		"recordingId": recording.ExternalId,
		"performerId": performer.ExternalId,
	}).Debug("Import recording")

	return &common.UpsertRecordingResponse{
		Status: common.StatusCode_Ok,
	}, nil
}

func (server *contentServer) DeleteRecording(ctx context.Context, req *common.RecordingIdentifier) (*common.DeleteRecordingResponse, error) {
	log := server.log.WithField("recordingId", req.Id)
	log.Debug("Received delete recording")

	err := server.app.RecordingService.DeleteRecording(context.Background(), req.Id)
	switch err {
	case nil:
		return &common.DeleteRecordingResponse{
			Status: common.StatusCode_Ok,
		}, nil

	case infinity.RecordingInPremiumUserCollectionErr:
		return &common.DeleteRecordingResponse{
			Status: common.StatusCode_RecordingInPremiumUserCollection,
		}, nil

	case infinity.RecordingHasViewsErr:
		return &common.DeleteRecordingResponse{
			Status: common.StatusCode_RecordingHasViews,
		}, nil

	default:
		log.WithError(err).Error("Unknown error deleting recording")

		return &common.DeleteRecordingResponse{
			Status: common.StatusCode_InternalServerErr,
		}, nil
	}
}
