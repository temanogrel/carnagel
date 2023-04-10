package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain/aphrodite"
	"git.misc.vee.bz/carnagel/go-ecosystem/infrastructure/factory"
	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/domain/cleanup"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/domain/payment"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/domain/recording"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/infrastructure/redis"
	"github.com/go-pg/pg"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"gopkg.in/olivere/elastic.v5"
	"os"
	"os/signal"
	"strings"
)

var (
	ctx, cancel = context.WithCancel(context.Background())
	hostname    string
	consul      ecosystem.ConsulClient
	logger      logrus.FieldLogger
	amqpCh      *amqp.Channel

	app *infinity.Application
)

var (
	command = flag.String("command", "", "Command to run")
	file    = flag.String("file", "", "Csv file to import")
)

func init() {
	flag.Parse()

	var err error

	consul, logger, hostname, err = factory.Bootstrap("infinity-client")
	if err != nil {
		panic(err)
	}

	amqpConn, err := connectToRabbitmq()
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to rabbitmq")
	}

	redis, err := redis.NewRedisConnection(consul.API())
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to redis")
	}

	elasticClient, err := getElasticSearchClient()
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to elasticsearch")
	}

	dbUser := consul.GetString("infinity/postgres/user", "infinity")
	dbName := consul.GetString("infinity/postgres/name", "infinity")
	dbPass := consul.GetString("infinity/postgres/pass", "infinity")
	dbHost := consul.GetString("infinity/postgres/host", "localhost")
	dbPort := consul.GetString("infinity/postgres/port", "5432")

	// use go.pg for the remainder of the time
	database := pg.Connect(&pg.Options{
		User:     dbUser,
		Addr:     fmt.Sprintf("%s:%s", dbHost, dbPort),
		Password: dbPass,
		Database: dbName,
	})

	app = &infinity.Application{
		Amqp:          amqpConn,
		Logger:        logger,
		Redis:         redis,
		DB:            database,
		Consul:        consul,
		ElasticSearch: elasticClient,

		// Ecosystem clients
		AphroditeClient: aphrodite.NewClient(
			consul.GetString("coordinator/aphrodite/api", "http://api.aphrodite.vee.bz"),
			consul.GetString("coordinator/aphrodite/token", "helloWorld"),
			logger,
		),
	}

	app.CleanupService = cleanup.NewCleanupService(app)
	app.CleanupRepository = cleanup.NewCleanupRepository(app)
	app.RecordingRepository = recording.NewRepository(app)
	app.PaymentPlanRepository = payment.NewPaymentPlanRepository(app)

	amqpCh, err = ecosystem.GetAmqpChannel(app.Amqp, 0)
	if err != nil {
		app.Logger.WithError(err).Fatal("Failed to create amqp channel")
	}
}

func main() {
	defer amqpCh.Close()
	defer cancel()

	switch *command {
	case "set-non-deletable-recordings":
		go func() {
			setNonDeletableRecordings()
			cancel()
		}()

	case "perform-cleanup":
		go func() {
			app.CleanupService.RunCleanup(ctx)
			cancel()
		}()

	default:
		logger.WithField("command", *command).Warn("Unknown command provided")
		return
	}

	signalCh := make(chan os.Signal, 1)

	// Listen for shutdown
	signal.Notify(signalCh, os.Interrupt)

	select {
	case <-signalCh:
		return
	case <-ctx.Done():
		return
	}
}

func setNonDeletableRecordings() {
	if *file == "" {
		logger.Fatal("No path to csv file provided")
	}

	file, err := os.Open(*file)
	if err != nil {
		logger.WithError(err).Fatalf("Failed to open file")
	}

	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = 5

	records, err := reader.ReadAll()
	if err != nil {
		logger.WithError(err).Fatalf("Failed to read csv records")
	}

	entries := make([]infinity.NonDeletableRecording, 0)

	for _, record := range records {
		hash := strings.Replace(record[0], "http://upsto.re/", "", -1)
		hash = strings.Replace(hash, "https://upstore.net/", "", -1)

		entries = append(entries, infinity.NonDeletableRecording{UpstoreHash: hash, Filename: record[1]})
	}

	if err := app.CleanupRepository.UpdateNonDeletableRecordings(ctx, entries); err != nil {
		logger.WithError(err).Fatalf("Failed to update non deletable recordings")
	}
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

	user := consul.GetString("infinity/rabbitmq/user", "infinity")
	pass := consul.GetString("infinity/rabbitmq/pass", "infinity")
	vhost := consul.GetString("infinity/rabbitmq/vhost", "infinity")

	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", user, pass, instances[0].ServiceAddress, instances[0].ServicePort, vhost)

	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to connect to rabbitmq service")
	}

	return conn, nil
}

func getElasticSearchClient() (*elastic.Client, error) {
	services, _, err := consul.API().Catalog().Service("elasticsearch", "", &consulapi.QueryOptions{})
	if err != nil {
		logger.WithError(err).Fatal("Failed to retrieve elasticsearch")
	}

	var addresses []string
	for _, service := range services {
		addresses = append(addresses, fmt.Sprintf("http://%s:%d", service.ServiceAddress, service.ServicePort))
	}

	return elastic.NewClient(elastic.SetURL(addresses...))
}
