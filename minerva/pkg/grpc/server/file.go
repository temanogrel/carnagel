package server

import (
	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"git.misc.vee.bz/carnagel/minerva-bindings/src"
	"git.misc.vee.bz/carnagel/minerva/pkg"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type fileServer struct {
	app *minerva.Application
}

func NewFileServer(app *minerva.Application) pb.FileServer {
	return &fileServer{app}
}

func (server *fileServer) HasBulk(ctx context.Context, request *pb.HasBulkRequest) (*pb.HasBulkResponse, error) {
	logger := server.app.Logger.WithFields(logrus.Fields{
		"method":    "HasBulk",
		"request":   request,
		"requestId": uuid.NewV4().String(),
	})

	result := map[uint64]bool{}
	ids := make([]ecosystem.ExternalId, len(request.Ids))
	for i, id := range request.Ids {
		result[id] = false
		ids[i] = ecosystem.ExternalId(id)
	}

	recordings, err := server.app.FileRepository.GetByExternalIds(ids)
	if err != nil {
		logger.WithError(err).Error("Failed to retrieve files by external id")

		return &pb.HasBulkResponse{
			Status: pb.StatusCode_InternalServerErr,
		}, nil
	}

	for _, recording := range recordings {
		result[uint64(recording.ExternalId)] = true
	}

	return &pb.HasBulkResponse{
		Status: pb.StatusCode_Ok,
		Result: result,
	}, nil
}

func (server *fileServer) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	logger := server.app.Logger.WithFields(logrus.Fields{
		"method":    "Get",
		"request":   request,
		"requestId": uuid.NewV4().String(),
	})

	var err error
	var file *minerva.File

	if request.Uuid != "" {

		logger.Debug("Retrieve file by uuid from database")

		id, err := uuid.FromString(request.Uuid)
		if err != nil {
			return &pb.GetResponse{
				Status: pb.StatusCode_InvalidUuid,
			}, nil
		}

		file, err = server.app.FileRepository.GetByUuid(id)
	} else {

		logger.Debug("Retrieving file by location from database")
		file, err = server.app.FileRepository.GetByLocation(request.Hostname, request.Path)
	}

	switch err {
	case nil:

		// Trying to figure out what causes the panic in #19
		if file == nil && err == nil {
			logger.Error("File and err are both nil, this should never happen")

			return &pb.GetResponse{
				Status: pb.StatusCode_InternalServerErr,
			}, nil
		}

		return &pb.GetResponse{
			Status: pb.StatusCode_Ok,
			File: &pb.FileData{
				Uuid:            file.Uuid.String(),
				ExternalId:      uint64(file.ExternalId),
				Checksum:        file.Checksum,
				Type:            pb.FileType(file.Type),
				Hostname:        string(file.Hostname),
				Path:            file.Path,
				UpstoreHash:     file.UpstoreHash,
				PendingUpload:   file.PendingUpload,
				PendingDeletion: file.PendingDeletion,
				Size:            file.Size,
				Meta:            file.Meta.ToProtoBufferStruct(),
				CreatedAt:       file.CreatedAt.String(),
				UpdatedAt:       file.UpdatedAt.String(),
			},
		}, nil

	case minerva.FileNotFoundErr:
		return &pb.GetResponse{
			Status: pb.StatusCode_FileNotFound,
		}, nil

	default:
		logger.WithError(err).Warn("Failed to retrieve file data")

		return &pb.GetResponse{
			Status: pb.StatusCode_InternalServerErr,
		}, nil
	}
}

// CRUD tasks
func (server *fileServer) Create(ctx context.Context, request *pb.CreateRequest) (*pb.CreateResponse, error) {
	logger := server.app.Logger.WithField("request", request)

	data := minerva.CreateFile{
		ExternalId: ecosystem.ExternalId(request.ExternalId),

		Meta: ecosystem.FileMetadata{}.FromProtoBufferStruct(request.Meta),
		Type: ecosystem.FileType(request.Type),

		Size:     request.Size,
		Path:     request.Path,
		Hostname: request.Hostname,
		Checksum: request.Checksum,
	}

	if file, err := server.app.FileService.Create(&data); err == nil {
		logger.Debug("Successfully registered new file")

		return &pb.CreateResponse{
			Status: pb.StatusCode_Ok,
			Uuid:   file.Uuid.String(),
		}, nil
	} else {
		logger.WithError(err).Warn("Failed to save recording")
	}

	return &pb.CreateResponse{
		Status: pb.StatusCode_InternalServerErr,
	}, nil
}

func (server *fileServer) Update(ctx context.Context, request *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	logger := server.app.Logger.WithField("request", request).WithField("method", "Update")
	logger.Debug("File update request")

	id, err := uuid.FromString(request.Uuid)
	if err != nil {
		logger.Warn("Invalid UUID")

		return &pb.UpdateResponse{Status: pb.StatusCode_InvalidUuid}, err
	}

	data := &minerva.UpdateFile{
		Uuid: id,
		Meta: ecosystem.FileMetadata{}.FromProtoBufferStruct(request.Meta),

		Size:     request.Size,
		Path:     request.Path,
		Hostname: request.Hostname,
		Checksum: request.Checksum,
	}

	_, err = server.app.FileService.Update(data)
	switch err {
	case nil:
		logger.Info("Successfully updated file")
		return &pb.UpdateResponse{Status: pb.StatusCode_Ok}, nil

	case minerva.FileNotFoundErr:
		logger.Warn("File not found")
		return &pb.UpdateResponse{Status: pb.StatusCode_FileNotFound}, nil

	default:
		logger.WithError(err).Warn("Unexpected error occurred during Update")
		return &pb.UpdateResponse{Status: pb.StatusCode_InternalServerErr}, err
	}
}

func (server *fileServer) RequestDeletion(ctx context.Context, request *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	logger := server.app.Logger.WithField("request", request).WithField("method", "RequestDeletion")
	logger.Debug("File deletion request")

	var err error
	if request.Uuid != "" {
		id, err := uuid.FromString(request.Uuid)
		if err != nil {
			logger.Warn("Invalid UUID")

			return &pb.DeleteResponse{Status: pb.StatusCode_InvalidUuid}, err
		}

		err = server.app.FileService.ScheduleDeletion(id)
	} else {

		// Delete all files by external id and type
		err = server.app.FileService.ScheduleDeletionByExternalIdAndType(
			ecosystem.ExternalId(request.ExternalId),
			ecosystem.FileType(request.Type),
		)
	}

	switch err {
	case nil, minerva.FileAlreadyPendingDeletionErr:
		logger.Info("Successfully scheduled for deletion")
		return &pb.DeleteResponse{Status: pb.StatusCode_Ok}, nil

	case minerva.FileNotFoundErr:
		logger.Warn("File not found")
		return &pb.DeleteResponse{Status: pb.StatusCode_FileNotFound}, nil

	default:
		logger.WithError(err).Warn("Unexpected error occurred during RequestDeletion")
		return &pb.DeleteResponse{Status: pb.StatusCode_InternalServerErr}, err
	}
}

// Upload a specific file
func (server *fileServer) RequestUpload(ctx context.Context, request *pb.UploadRequest) (*pb.UploadResponse, error) {
	logger := server.app.Logger.WithField("request", request).WithField("method", "RequestUpload")
	logger.Debug("File upload request")

	id, err := uuid.FromString(request.Uuid)
	if err != nil {
		logger.Warn("Invalid UUID")

		return &pb.UploadResponse{Status: pb.StatusCode_InvalidUuid}, err
	}

	err = server.app.FileService.ScheduleUpload(id, request.Name)
	switch err {
	case nil, minerva.FileAlreadyPendingUploadErr:
		logger.Info("Successfully scheduled file upload")
		return &pb.UploadResponse{Status: pb.StatusCode_Ok}, nil

	case minerva.FileNotFoundErr:
		logger.Warn("File not found")
		return &pb.UploadResponse{Status: pb.StatusCode_FileNotFound}, nil

	default:
		logger.WithError(err).Warn("Unexpected error occurred during RequestUpload")
		return &pb.UploadResponse{Status: pb.StatusCode_InternalServerErr}, err
	}
}
