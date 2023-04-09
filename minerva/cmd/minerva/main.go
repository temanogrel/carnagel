package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"

	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"git.misc.vee.bz/carnagel/go-ecosystem/infrastructure/factory"
	"git.misc.vee.bz/carnagel/minerva/pkg"
	"git.misc.vee.bz/carnagel/minerva/pkg/file"
	"git.misc.vee.bz/carnagel/minerva/pkg/grpc"
	"git.misc.vee.bz/carnagel/minerva/pkg/http"
	"git.misc.vee.bz/carnagel/minerva/pkg/loadbalancer"
	"git.misc.vee.bz/carnagel/minerva/pkg/server"
	"git.misc.vee.bz/carnagel/sql-migrations"
	"git.misc.vee.bz/carnagel/sql-migrations/driver"
	"github.com/go-pg/pg"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/api/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var (
	devMode      = flag.Bool("dev", false, "If we are running in dev mode")
	selfRegister = flag.Bool("self-register", false, "If the system should register it's self in consul")

	consul ecosystem.ConsulClient
	logger logrus.FieldLogger
)

func init() {
	flag.Parse()

	var err error

	consul, logger, _, err = factory.Bootstrap("minerva")
	if err != nil {
		panic(err)
	}
}

func main() {

	publicIp := getIpv4OfDevice(os.Getenv("PUBLIC_NIC"))
	if publicIp == "" {
		logger.Fatal("Failed to start minerva, couldnt determine public network interface")
	}

	dbUser := consul.GetString("minerva/postgres/user", "minerva")
	dbName := consul.GetString("minerva/postgres/name", "minerva")
	dbPass := consul.GetString("minerva/postgres/pass", "minerva")
	dbHost := consul.GetString("minerva/postgres/host", "localhost")
	dbPort := consul.GetString("minerva/postgres/port", "5432")

	grpcBindAddr := os.Getenv("NOMAD_ADDR_grpc")
	httpBindAddrPublic := fmt.Sprintf("%s:%s", publicIp, os.Getenv("NOMAD_PORT_http_public"))
	httpBindAddrPrivate := os.Getenv("NOMAD_ADDR_http_private")

	// go.pg does not expose the underlying sqlx connection
	go migrate(dbHost, dbPort, dbUser, dbPass, dbName)

	// use go.pg for the remainder of the time
	database := pg.Connect(&pg.Options{
		User:     dbUser,
		Addr:     fmt.Sprintf("%s:%s", dbHost, dbPort),
		Password: dbPass,
		Database: dbName,
		PoolSize: 400,
	})

	prometheusService, err := getPrometheusService()
	if err != nil {
		logger.WithError(err).Fatal("Failed to retrieve prometheus")
	}

	prometheusConfig := prometheus.Config{
		Address: fmt.Sprintf("http://%s:%d", prometheusService.ServiceAddress, prometheusService.ServicePort),
	}

	prometheusClient, err := prometheus.New(prometheusConfig)
	if err != nil {
		logger.WithError(err).Fatal("Failed to create prometheus client")
	}

	amqpConn, err := connectToRabbitmq()
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to rabbitmq")
	}

	prometheusQueryApi := prometheus.NewQueryAPI(prometheusClient)

	application := &minerva.Application{
		Amqp:             amqpConn,
		Logger:           logger,
		Database:         database,
		Consul:           consul,
		PrometheusClient: prometheusQueryApi,

		// Settings
		GrpcBindAddr:        grpcBindAddr,
		HttpBindAddrPublic:  httpBindAddrPublic,
		HttpBindAddrPrivate: httpBindAddrPrivate,

		DevMode:     *devMode,
		MinionPort:  6000,
		MinionToken: "read",

		ServerCollection: minerva.NewServers(),
	}

	application.FileService = file.NewFileService(application)
	application.FileRepository = file.NewFileRepository(application)
	application.LoadBalancer = loadbalancer.NewService(application)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	logger.Infof("Starting Minerva")

	if *selfRegister {
		application.RegisterInConsul()
	}

	go server.NewMetricCollector(application).Run(ctx)
	go server.NewServiceService(application).Run(ctx)

	go http.StartServer(ctx, application)
	go grpc.StartServer(ctx, application)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	select {
	case <-c:
		logger.Info("Received SIGINT, terminating...")

		if *selfRegister {
			application.DeregisterInConsul()
		}

		// cleanup
		application.FileService.CleanUp()
		cancel()
	}
}

func migrate(host, port, user, password, dbname string) {

	url := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	conn, err := sqlx.Connect("postgres", url)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	var migrationsPath string
	if *devMode {
		migrationsPath = "./migrations"
	} else {
		migrationsPath = consul.GetString("minerva/migrations-path", "./migrations")
	}

	mig := migrations.CreateFromDirectory(migrationsPath)
	mig.Up(driver.NewPostgresDriver(conn), true)
}

func getPrometheusService() (*consulapi.CatalogService, error) {
	service, _, err := consul.API().Catalog().Service("prometheus", "", &consulapi.QueryOptions{})
	if err != nil {
		return nil, err
	}

	if len(service) != 1 {
		return nil, errors.New("Unexpected number of prometheus services returned")
	}

	return service[0], nil
}

// connectToRabbitmq queries consul and connects to the available rabbitmq
func connectToRabbitmq() (*amqp.Connection, error) {
	instances, _, err := consul.API().Catalog().Service("rabbitmq", "", &consulapi.QueryOptions{})
	if err != nil {
		return nil, err
	}

	if len(instances) == 0 {
		return nil, errors.New("No rabbitmq service available in consul")
	}

	user := consul.GetString("minerva/rabbitmq/user", "minerva")
	pass := consul.GetString("minerva/rabbitmq/pass", "minerva")
	vhost := consul.GetString("minerva/rabbitmq/vhost", "minerva")

	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", user, pass, instances[0].ServiceAddress, instances[0].ServicePort, vhost)

	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to connect to rabbitmq service")
	}

	return conn, nil
}

func getIpv4OfDevice(name string) string {
	interfaces, _ := net.Interfaces()
	for _, inter := range interfaces {
		if inter.Name != name {
			continue
		}

		if addrs, err := inter.Addrs(); err == nil {
			for _, addr := range addrs {
				if tmp, ok := addr.(*net.IPNet); ok {
					if ip := tmp.IP.To4(); ip != nil {
						return ip.String()
					}
				}
			}
		}
	}

	return ""
}
