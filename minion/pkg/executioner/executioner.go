package executioner

import (
	context2 "context"
	"os"
	"time"

	"git.misc.vee.bz/carnagel/minerva-bindings/src"
	"git.misc.vee.bz/carnagel/minion/pkg"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type executionerService struct {
	app    *minion.Application
	logger logrus.FieldLogger
}

func NewExecutionerService(app *minion.Application) minion.ExecutionerService {
	return &executionerService{
		app:    app,
		logger: app.Logger.WithField("component", "executioner"),
	}
}

func (service *executionerService) ProcessRelocationRequests(ctx context2.Context) {

	service.logger.Debug("Starting process relocation listener")

	identifier := &pb.ServerIdentifier{
		Hostname: service.app.Hostname,
	}

	client, err := service.app.MinionDelegationClient.RelocateRequests(ctx, identifier)
	if err != nil {
		service.logger.WithError(err).Error("Failed to create deletion requests client")
		return
	}

	for {
		request, err := client.Recv()
		if err != nil {
			service.logger.WithError(err).Error("An error occurred receiving a deletion request")
			return
		}

		id, err := uuid.FromString(request.Uuid)
		if err != nil {
			service.logger.
				WithError(err).
				WithField("uuid", request.Uuid).
				Error("Invalid uuid provided")

			continue
		}

		if err := service.app.FileService.Transfer(ctx, id, request.TargetHost); err != nil {
			service.logger.WithError(err).Warn("Failed to transfer file")
			continue
		}
	}
}

func (service *executionerService) ProcessDeletions(ctx context2.Context) {

	service.logger.Debug("Starting process deletion listener")

	identifier := &pb.ServerIdentifier{
		Hostname: service.app.Hostname,
	}

	client, err := service.app.MinionDelegationClient.DeletionRequests(ctx, identifier)
	if err != nil {
		service.logger.WithError(err).Error("Failed to create deletion requests client")
		return
	}

	deleteFromMinerva := func(uuid string, logger logrus.FieldLogger) {
		resp, err := service.app.MinionDelegationClient.FileDeleted(context.TODO(), &pb.FileIdentifier{
			Uuid: uuid,
		})

		if err != nil {
			logger.WithError(err).Error("failed to notify minerva about deleted file")
			return
		}

		if resp.Status != pb.StatusCode_Ok {
			logger.WithField("status", resp.Status).Error("Unexpected response code")
			return
		}

		logger.Info("Deleted the file successfully")
	}

	for {
		request, err := client.Recv()
		if err != nil {
			service.logger.WithError(err).Error("An error occurred receiving a deletion request")
			return
		}

		go func() {

			logger := service.logger.WithField("request", request)
			logger.Debug("Received deletion request")

			if _, err := os.Stat(request.Path); os.IsNotExist(err) {
				deleteFromMinerva(request.Uuid, logger)
				return
			}

			if err := os.Remove(request.Path); err != nil {
				logger.WithError(err).Error("An error occurred trying to delete the file")
				return
			}

			deleteFromMinerva(request.Uuid, logger)
		}()
	}
}

// ProcessUploads gets incoming requests from minerva for files o the specified server to upload
func (service *executionerService) ProcessUploads(ctx context2.Context) {

	service.app.Logger.Debug("Starting process uploads listener")

	identifier := &pb.ServerIdentifier{
		Hostname: service.app.Hostname,
	}

	client, err := service.app.MinionDelegationClient.UploadRequests(ctx, identifier)
	if err != nil {
		service.app.Logger.WithError(err).Error("Failed to create uploads requests client")
		return
	}

	workerQueue := make(chan *pb.ReUploadRequest)

	for i := 0; i <= 10; i++ {
		go service.processUploadQueue(ctx, workerQueue)
	}

	for {
		request, err := client.Recv()
		if err != nil {
			service.app.Logger.WithError(err).Error("An error occurred receiving a upload request")
			return
		}

		workerQueue <- request
	}
}

func (service *executionerService) processUploadQueue(ctx context2.Context, queue chan *pb.ReUploadRequest) {

	for {
		select {
		case <-ctx.Done():
			return

		case request, ok := <-queue:
			if !ok {
				return
			}

			logger := service.app.Logger.WithField("request", request)
			logger.WithField("recordingId", request.ExternalId)

			if err := service.app.AphroditeClient.SetRecordingStateById(request.ExternalId, "uploading"); err != nil {
				logger.WithError(err).Error("Failed to notify aphrodite about state change")
				continue
			}

			if err := service.upload(request); err != nil {
				logger.WithError(err).Error("Failed to upload file")

				service.app.AphroditeClient.SetRecordingStateById(request.ExternalId, "uploading_failed")
			} else {

				service.app.AphroditeClient.SetRecordingStateById(request.ExternalId, "uploaded")
			}
		}
	}
}

func (service *executionerService) upload(request *pb.ReUploadRequest) error {

	logger := service.logger.WithFields(logrus.Fields{
		"request":     request,
		"recordingId": request.ExternalId,
	})

	logger.Debug("Received upload request")

	file, err := os.Open(request.Path)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Error("File does not exist on this server, it will be deleted")

			// delete all files with that external id
			resp, err := service.app.FileClient.RequestDeletion(context.TODO(), &pb.DeleteRequest{
				ExternalId: request.ExternalId,
			})

			if err != nil {
				logger.WithError(err).Error("Failed to delete all files with provided external id")
				return err
			}

			if resp.Status != pb.StatusCode_Ok {
				logger.WithField("status", resp.Status).Error("Unexpected response code")
				return nil
			}

			if err := service.app.AphroditeClient.DeleteRecording(request.ExternalId); err != nil {
				logger.WithError(err).Error("Failed to delete the recording")
				return err
			}

			logger.Info("Deleted the file and recording successfully")

		} else {
			logger.WithError(err).Error("Failed to open file")
		}

		return err
	}

	defer file.Close()

	startedAt := time.Now()

	hash, err := service.app.UpstoreClient.Upload(request.Name, file)
	if err != nil {
		return errors.Wrap(err, "Failed to upload to upstore")
	}

	setHashRequest := &pb.SetUpstoreHashRequest{
		Uuid: request.Uuid,
		Hash: hash,
	}

	resp, err := service.app.MinionDelegationClient.SetUpstoreHash(context.TODO(), setHashRequest)
	if err != nil {
		return errors.Wrapf(err, "Failed to set upstore hash in minerva")
	}

	if resp.Status != pb.StatusCode_Ok {
		return errors.Errorf("Unexpected status code %d received from minerva", resp.Status)
	}

	logger.
		WithField("hash", hash).
		WithField("duration", time.Since(startedAt).Seconds()).
		Debug("Successfully uploaded file to upstore")

	return nil
}
