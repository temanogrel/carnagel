package main

import (
	"context"
	"flag"
	"os"
	"os/signal"

	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain/aphrodite"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain/minerva"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain/upstore"
	"git.misc.vee.bz/carnagel/go-ecosystem/infrastructure/factory"
	"git.misc.vee.bz/carnagel/go-ecosystem/protobuf/common"
	"git.misc.vee.bz/carnagel/minerva-bindings/src"
	"git.misc.vee.bz/carnagel/minion/pkg"
	"git.misc.vee.bz/carnagel/minion/pkg/executioner"
	"git.misc.vee.bz/carnagel/minion/pkg/file"
	"git.misc.vee.bz/carnagel/minion/pkg/http"
	"git.misc.vee.bz/carnagel/minion/pkg/integrity"
	"github.com/sirupsen/logrus"
)

var (
	selfRegister = flag.Bool("self-register", false, "Register it's self in consulApi")

	logger   logrus.FieldLogger
	consul   ecosystem.ConsulClient
	hostname string
)

func init() {
	flag.Parse()

	var err error

	consul, logger, hostname, err = factory.Bootstrap("minion")
	if err != nil {
		panic(err)
	}
}

func main() {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	minervaConn, err := factory.NewMinervaConnection(ctx, consul)
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to minerva")
	}

	infinityConn, err := factory.NewInfinityConnection(ctx, consul)
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to infinity")
	}

	application := &minion.Application{
		Logger: logger,
		Consul: consul,

		// Config
		DataDir:             consul.GetString("minion/data-dir", "/tmp"),
		Hostname:            hostname,
		UploadFormMaxMemory: 1024 * 1024 * 1024 * 3, // 2GB

		// Http auth
		HttpReadToken:  consul.GetString("minion/http/read-token", "read"),
		HttpWriteToken: consul.GetString("minion/http/write-token", "write"),

		// grpc clients to minerva
		FileClient:             pb.NewFileClient(minervaConn),
		MinionDelegationClient: pb.NewMinionDelegationClient(minervaConn),

		// grpc clients to infinity
		DeviceTrackingClient:    common.NewDeviceTrackingClient(infinityConn),
		BandwidthTrackingClient: common.NewBandwidthTrackingClient(infinityConn),

		// http clients
		UpstoreClient: upstore.NewUpstoreClient(
			consul.GetString("minion/upstore/email", ""),
			consul.GetString("minion/upstore/password", ""),
			logger,
		),

		MinervaClient: minerva.NewClient(minervaConn,
			hostname,
			consul.GetString("minion/http/read-token", "read"),
			consul.GetString("minion/http/write-token", "write"),
		),

		AphroditeClient: aphrodite.NewClient(
			consul.GetString("aphrodite/api-root", "http://10g.api.aphrodite.vee.bz"),
			consul.GetString("aphrodite/api-token", "helloWorld"),
			logger,
		),
	}

	application.FileService = file.NewFileService(application)
	application.Executioner = executioner.NewExecutionerService(application)
	application.FilesystemIntegrity = integrity.NewFileSystemIntegrity(application)
	application.EnsureValidEnvironment()

	logger.Infof("Starting Minion")

	go http.StartServer(ctx, application)

	// Run the filesystem integrity based on consul configuration
	go application.FilesystemIntegrity.Run(ctx)

	// Run processing of the minion delegation events
	go application.ProcessDeletions(ctx)
	go application.ProcessUploads(ctx)
	go application.ProcessRelocations(ctx)

	// Don't register until all the subsystems are online
	if *selfRegister {
		application.RegisterInConsul()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	select {
	case <-c:
		logger.Info("Received SIGINT, terminating...")

		if *selfRegister {
			application.DeregisterInConsul()
		}

		// cleanup
		cancel()
	}
}
