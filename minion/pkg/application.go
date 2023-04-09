package minion

import (
	"context"
	"time"

	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"git.misc.vee.bz/carnagel/go-ecosystem/protobuf/common"
	"git.misc.vee.bz/carnagel/minerva-bindings/src"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

type Application struct {
	// Misc
	Logger logrus.FieldLogger
	Consul ecosystem.ConsulClient

	// Http clients
	UpstoreClient   ecosystem.UpstoreClient
	MinervaClient   ecosystem.MinervaClient
	AphroditeClient ecosystem.AphroditeClient

	// Grpc clients
	FileClient              pb.FileClient
	MinionDelegationClient  pb.MinionDelegationClient
	BandwidthTrackingClient common.BandwidthTrackingClient
	DeviceTrackingClient    common.DeviceTrackingClient

	// Business logic
	Executioner         ExecutionerService
	FileService         FileService
	FilesystemIntegrity FilesystemIntegrity

	// Config
	Hostname            string
	DataDir             string
	UploadFormMaxMemory int64

	// Auth
	HttpReadToken  string
	HttpWriteToken string
}

func (app *Application) EnsureValidEnvironment() {

	if unix.Access(app.DataDir, unix.W_OK) != nil {
		app.Logger.WithField("dataDir", app.DataDir).Fatal("No write access to data dir")
	}

	if unix.Access(app.DataDir, unix.R_OK) != nil {
		app.Logger.WithField("dataDir", app.DataDir).Fatal("No read access to data dir")
	}
}

func (app *Application) RegisterInConsul() {
	app.Logger.Info("Register in consul")

	_, err := app.Consul.API().Catalog().Register(&consulapi.CatalogRegistration{
		Node:    "localhost",
		Address: "localhost",
		NodeMeta: map[string]string{
			"internal_hostname": "localhost",
			"external_hostname": "localhost",
		},

		Service: &consulapi.AgentService{
			ID:      "minion",
			Service: "minion",
			Address: "localhost",
			Port:    6000,
		},
	}, &consulapi.WriteOptions{})

	if err != nil {
		app.Logger.WithError(err).Fatal("Failed to register in consul")
	}
}

func (app *Application) DeregisterInConsul() {
	app.Logger.Info("Deregister in consul")

	_, err := app.Consul.API().Catalog().Deregister(&consulapi.CatalogDeregistration{
		Node:      "localhost",
		Address:   "localhost",
		ServiceID: "minion",
	}, &consulapi.WriteOptions{})

	if err != nil {
		app.Logger.WithError(err).Fatal("Failed to deregister in consul")
	}
}

func (app *Application) ProcessRelocations(ctx context.Context) {
	for {
		app.Executioner.ProcessRelocationRequests(ctx)
		app.Logger.Warn("Process relocations terminated, sleeping 3 seconds and starting again")
		time.Sleep(3 * time.Second)
	}
}

func (app *Application) ProcessDeletions(ctx context.Context) {
	for {
		app.Executioner.ProcessDeletions(ctx)
		app.Logger.Warn("Process deletions terminated, sleeping 3 seconds and starting again")
		time.Sleep(3 * time.Second)
	}
}

func (app *Application) ProcessUploads(ctx context.Context) {
	for {
		app.Executioner.ProcessUploads(ctx)
		app.Logger.Warn("Process uploads terminated, sleeping 3 seconds and starting again")
		time.Sleep(3 * time.Second)
	}
}
