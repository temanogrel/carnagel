package minerva

import (
	"os"
	"strconv"

	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"github.com/go-pg/pg"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/prometheus/client_golang/api/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type Application struct {
	// Misc
	Amqp             *amqp.Connection
	Logger           logrus.FieldLogger
	Database         *pg.DB
	Consul           ecosystem.ConsulClient
	PrometheusClient prometheus.QueryAPI

	// Services
	LoadBalancer  LoadBalancer
	ServerService ServerDiscovery

	ServerCollection *Servers

	// Configuration
	GrpcBindAddr        string
	HttpBindAddrPublic  string
	HttpBindAddrPrivate string

	DevMode     bool
	MinionPort  uint
	MinionToken string

	// Modules
	FileService    FileService
	FileRepository FileRepository
}

func (app *Application) RegisterInConsul() {
	app.Logger.Info("Registering from consul")

	port, err := strconv.Atoi(os.Getenv("NOMAD_PORT_grpc"))
	if err != nil {
		app.Logger.WithError(err).Fatal("Failed to determine port")
	}

	_, err = app.Consul.API().Catalog().Register(&consulapi.CatalogRegistration{
		Node:    "localhost",
		Address: "localhost",

		Service: &consulapi.AgentService{
			ID:      "minerva",
			Service: "minerva",
			Address: "localhost",
			Port:    port,
		},
	}, &consulapi.WriteOptions{})

	if err != nil {
		app.Logger.WithError(err).Fatal("Failed to register in consul")
	}
}

func (app *Application) DeregisterInConsul() {
	app.Logger.Info("Deregistering from consul")

	_, err := app.Consul.API().Catalog().Deregister(&consulapi.CatalogDeregistration{
		Node:      "localhost",
		Address:   "localhost",
		ServiceID: "minerva",
	}, &consulapi.WriteOptions{})

	if err != nil {
		app.Logger.WithError(err).Fatal("Failed to deregister in consul")
	}
}
