package grpc

import (
	"net"

	"git.misc.vee.bz/carnagel/minerva-bindings/src"
	"git.misc.vee.bz/carnagel/minerva/pkg"
	"git.misc.vee.bz/carnagel/minerva/pkg/grpc/server"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func StartServer(ctx context.Context, app *minerva.Application) {
	listener, err := net.Listen("tcp", app.GrpcBindAddr)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterFileServer(grpcServer, server.NewFileServer(app))
	pb.RegisterMinionDelegationServer(grpcServer, server.NewMinionDelegationServer(app))
	pb.RegisterLoadBalancerServer(grpcServer, server.NewLoadBalancerServer(app))

	app.Logger.Infof("Started GrpcServer: %s", app.GrpcBindAddr)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			app.Logger.WithError(err).Fatal("Failed to start grpc server")
		}
	}()

	// block until we should terminate
	<-ctx.Done()

	grpcServer.GracefulStop()
}
