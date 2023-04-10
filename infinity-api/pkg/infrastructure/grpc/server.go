package grpc

import (
	"net"

	"time"

	"git.misc.vee.bz/carnagel/go-ecosystem/protobuf/common"
	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/infrastructure/grpc/server"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	procedureResponseTime = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "grpc_method_duration",
			Help: "Number of seconds spent on each procedure call",
			Objectives: map[float64]float64{
				0.5:  0.05,
				0.9:  0.01,
				0.99: 0.001,
			},
		},
		[]string{"procedure"},
	)

	procedureErrorCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_method_errors",
			Help: "Number of errors that have occurred over time",
		},
		[]string{"procedure"},
	)
)

func init() {
	prometheus.MustRegister(procedureResponseTime)
	prometheus.MustRegister(procedureErrorCount)
}

func StartGrpcServer(ctx context.Context, app *infinity.Application) {
	listener, err := net.Listen("tcp", app.GrpcBindAddr)
	if err != nil {
		app.Logger.
			WithField("addr", app.GrpcBindAddr).
			WithError(err).
			Fatal("Failed to listen to grpc bind addr")
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(WrapUnaryInterceptor()))

	common.RegisterBandwidthTrackingServer(grpcServer, server.NewBandwidthTrackingServer(app))
	common.RegisterDeviceTrackingServer(grpcServer, server.NewDeviceTrackingServer(app))
	common.RegisterContentServer(grpcServer, server.NewContentServer(app))

	go func() {
		app.Logger.Infof("Starting grpc server: %s", app.GrpcBindAddr)

		if err := grpcServer.Serve(listener); err != nil {
			app.Logger.WithError(err).Fatal("Failed to start grpc server")
		}
	}()

	// block until we should terminate
	<-ctx.Done()

	listener.Close()
	grpcServer.GracefulStop()
}

func WrapUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		startedAt := time.Now()

		resp, err = handler(ctx, req)
		if err != nil {
			procedureErrorCount.WithLabelValues(info.FullMethod).Inc()
		}

		procedureResponseTime.
			WithLabelValues(info.FullMethod).
			Observe(time.Since(startedAt).Seconds())

		return resp, err
	}
}
