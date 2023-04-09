package server

import (
	"git.misc.vee.bz/carnagel/minerva-bindings/src"
	"git.misc.vee.bz/carnagel/minerva/pkg"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type minionDelegationServer struct {
	app    *minerva.Application
	logger logrus.FieldLogger
}

func NewMinionDelegationServer(application *minerva.Application) pb.MinionDelegationServer {
	return &minionDelegationServer{
		app:    application,
		logger: application.Logger.WithField("component", "minion_delegation"),
	}
}

func (server *minionDelegationServer) RelocateRequests(identifier *pb.ServerIdentifier, stream pb.MinionDelegation_RelocateRequestsServer) error {

	logger := server.logger.WithField("hostname", identifier.Hostname)
	logger.Debug("Streaming deletion requests")

	server.app.ServerCollection.ServersMtx.RLock()
	s, ok := server.app.ServerCollection.Servers[minerva.Hostname(identifier.Hostname)]
	server.app.ServerCollection.ServersMtx.RUnlock()

	if !ok {
		logger.Errorf("Unknown server %s tried to connect", identifier.Hostname)
		return nil
	}

	for {
		req, ok := <-s.RelocateRequests
		if !ok {
			logger.Warn("Relocation requests channel closed")
			return nil
		}

		request := &pb.RelocateRequest{
			Uuid:       req.Uuid.String(),
			TargetHost: string(req.TargetHost),
		}

		if err := stream.Send(request); err != nil {

			// re-queue the request
			s.RelocateRequests <- req

			logger.WithError(err).Error("Failed to send a request to the minion")
			return nil
		}

		logger.
			WithField("uuid", req.Uuid).
			WithField("target", req.TargetHost).
			Debug("Send file to be relocated")
	}
}

func (server *minionDelegationServer) DeletionRequests(identifier *pb.ServerIdentifier, stream pb.MinionDelegation_DeletionRequestsServer) error {

	logger := server.logger.WithField("hostname", identifier.Hostname)
	logger.Debug("Streaming deletion requests")

	server.app.ServerCollection.ServersMtx.RLock()
	s, ok := server.app.ServerCollection.Servers[minerva.Hostname(identifier.Hostname)]
	server.app.ServerCollection.ServersMtx.RUnlock()

	if !ok {
		logger.Errorf("Unknown server %s tried to connect", identifier.Hostname)
		return nil
	}

	for {
		id, ok := <-s.DeletionRequests
		if !ok {
			logger.Warn("Deletion requests channel closed")
			return nil
		}

		file, err := server.app.FileRepository.GetByUuid(id)
		if err != nil {
			server.logger.
				WithError(err).
				WithField("id", id.String()).
				Error("Failed to retrieve file for deletion")

			continue
		}

		request := &pb.DeletionRequest{
			Uuid: id.String(),
			Path: file.Path,
		}

		loggerCtx := logger.WithField("request", request)

		if err := stream.Send(request); err != nil {

			// Queue it again
			s.DeletionRequests <- id

			loggerCtx.WithError(err).Warn("Failed to send deletion request")
			return err
		}

		loggerCtx.Debug("Successfully send deletion request")
	}
}

func (server *minionDelegationServer) UploadRequests(identifier *pb.ServerIdentifier, stream pb.MinionDelegation_UploadRequestsServer) error {

	logger := server.logger.WithField("hostname", identifier.Hostname)
	logger.Debug("Streaming upload requests")

	server.app.ServerCollection.ServersMtx.RLock()
	s, ok := server.app.ServerCollection.Servers[minerva.Hostname(identifier.Hostname)]
	server.app.ServerCollection.ServersMtx.RUnlock()

	if !ok {
		logger.Errorf("Unknown server %s tried to connect", identifier.Hostname)
		return nil
	}

	for {
		id, ok := <-s.UploadRequests
		if !ok {
			logger.Warn("Upload requests channel closed")
			return nil
		}

		file, err := server.app.FileRepository.GetByUuid(id)
		if err != nil {
			logger.
				WithError(err).
				WithField("id", id.String()).
				Error("Failed to retrieve file for upload")

			continue
		}

		request := &pb.ReUploadRequest{
			Uuid:       id.String(),
			Path:       file.Path,
			Name:       file.OriginalFilename,
			ExternalId: uint64(file.ExternalId),
		}

		loggerCtx := logger.WithField("request", request)

		if err := stream.Send(request); err != nil {

			// Queue it again
			s.UploadRequests <- id

			loggerCtx.WithError(err).Warn("Failed to send upload request")
			return err
		}

		loggerCtx.Info("Successfully sent file upload request")
	}
}

func (server *minionDelegationServer) FileDeleted(ctx context.Context, identifier *pb.FileIdentifier) (*pb.FileDeletedResponse, error) {
	logger := server.logger.WithField("uuid", identifier.Uuid)

	if err := server.app.FileService.Delete(uuid.FromStringOrNil(identifier.Uuid)); err != nil {
		if err == minerva.FileNotFoundErr {
			logger.Warn("Failed to delete file record from database because it was not found")

			return &pb.FileDeletedResponse{Status: pb.StatusCode_FileNotFound}, nil
		}

		logger.
			WithError(err).
			Error("Unexpected error occurred during file deletion")

		return &pb.FileDeletedResponse{Status: pb.StatusCode_InternalServerErr}, err
	}

	logger.Debug("Finished deleted file")

	return &pb.FileDeletedResponse{
		Status: pb.StatusCode_Ok,
	}, nil
}

func (server *minionDelegationServer) SetUpstoreHash(ctx context.Context, request *pb.SetUpstoreHashRequest) (*pb.SetUpstoreHashResponse, error) {
	logger := server.logger.WithFields(logrus.Fields{
		"uuid": request.Uuid,
		"hash": request.Hash,
	})

	if err := server.app.FileService.SetUpstoreHash(uuid.FromStringOrNil(request.Uuid), request.Hash); err != nil {
		logger.WithError(err).Warn("Failed to set file upstore.net hash")

		return &pb.SetUpstoreHashResponse{
			Status: pb.StatusCode_InternalServerErr,
		}, nil
	}

	logger.Debug("finished set upstore hash")

	return &pb.SetUpstoreHashResponse{
		Status: pb.StatusCode_Ok,
	}, nil
}
