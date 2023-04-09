package server

import (
	"git.misc.vee.bz/carnagel/minerva-bindings/src"
	"git.misc.vee.bz/carnagel/minerva/pkg"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
)

type loadBalancerServer struct {
	app *minerva.Application
}

func NewLoadBalancerServer(application *minerva.Application) pb.LoadBalancerServer {
	return &loadBalancerServer{app: application}
}

func (server *loadBalancerServer) RecommendStorage(ctx context.Context, req *pb.RecommendStorageRequest) (*pb.RecommendStorageResponse, error) {
	suggestedServer, selectedInterface, err := server.app.LoadBalancer.RecommendStorage(minerva.Hostname(req.OriginHostname), req.Size)
	if err != nil {
		if err == minerva.ServerNotAvailableErr {
			return &pb.RecommendStorageResponse{
				Status: pb.StatusCode_NoServerAvailable,
			}, nil
		}

		return &pb.RecommendStorageResponse{
			Status: pb.StatusCode_InternalServerErr,
		}, nil
	}

	var selectedHostname minerva.Hostname
	if selectedInterface == minerva.InternalInterface {
		selectedHostname = suggestedServer.InternalHostname
	} else {
		selectedHostname = suggestedServer.ExternalHostname
	}

	return &pb.RecommendStorageResponse{
		Status:   pb.StatusCode_Ok,
		Hostname: string(selectedHostname),
	}, nil
}

func (server *loadBalancerServer) RecommendDownload(ctx context.Context, req *pb.RecommendDownloadRequest) (*pb.RecommendDownloadResponse, error) {
	logger := server.app.Logger.WithField("request", req)

	path, err := server.app.LoadBalancer.RecommendDownload(minerva.Hostname(req.OriginHostname), req.Uuid)
	switch err {
	case nil:
		if req.CountHit {
			if err := server.app.FileRepository.TrackHit(uuid.FromStringOrNil(req.Uuid), req.ClientAddr); err != nil {
				logger.WithError(err).Warn("Failed to register hit")
			}
		}

		resp := &pb.RecommendDownloadResponse{
			Status: pb.StatusCode_Ok,

			Path:     path.File.Path,
			Meta:     path.File.Meta.ToProtoBufferStruct(),
			FileType: pb.FileType(path.File.Type),
		}

		if path.Edge != nil {
			resp.Edge = string(path.Edge.ExternalHostname)
		}

		if path.InterfaceToUse == "internal" {
			resp.Origin = string(path.Origin.InternalHostname)
		} else {
			resp.Origin = string(path.Origin.ExternalHostname)
		}

		logger.WithField("path", path).Debug("Recommended download")

		return resp, nil

	case minerva.FileNotFoundErr:
		logger.Warn("File not found")

		return &pb.RecommendDownloadResponse{
			Status: pb.StatusCode_FileNotFound,
		}, nil

	default:
		logger.WithError(err).Warn("Unexpected server error")

		return &pb.RecommendDownloadResponse{
			Status: pb.StatusCode_InternalServerErr,
		}, nil
	}
}
