package main

import (
	"context"
	"flag"
	"fmt"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/domain/cleanup"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/domain/crypto"
	"github.com/blockcypher/gobcy"
	"os"
	"os/signal"
	"time"

	"runtime"

	"github.com/sasha-s/go-deadlock"

	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain/aphrodite"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain/minerva"
	"git.misc.vee.bz/carnagel/go-ecosystem/infrastructure/factory"
	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/domain/bandwidth"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/domain/email"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/domain/payment"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/domain/performer"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/domain/recording"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/domain/user"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/infrastructure/grpc"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/infrastructure/http"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/infrastructure/http/rpc"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/infrastructure/redis"
	"git.misc.vee.bz/carnagel/sql-migrations"
	"git.misc.vee.bz/carnagel/sql-migrations/driver"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-pg/pg"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"gopkg.in/olivere/elastic.v5"
)

var (
	logger   logrus.FieldLogger
	consul   ecosystem.ConsulClient
	hostname string
)

func init() {
	flag.Parse()

	var err error

	consul, logger, hostname, err = factory.Bootstrap("infinity-api")
	if err != nil {
		panic(err)
	}
}

func main() {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

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

	// migrate
	migrate(dbHost, dbPort, dbUser, dbPass, dbName)

	// use go.pg for the remainder of the time
	database := pg.Connect(&pg.Options{
		User:     dbUser,
		Addr:     fmt.Sprintf("%s:%s", dbHost, dbPort),
		Password: dbPass,
		Database: dbName,
		// Everything below here is very experimental to see if we can get rid of connection timeout errors
		// Default is 10 * runtime.NumCPU(), so 50% increased pool
		PoolSize: 15 * runtime.NumCPU(),
		// Close connections after being idle for 15 mins
		IdleTimeout: time.Minute * 15,
	})

	minervaConn, err := factory.NewMinervaConnection(ctx, consul)
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to minerva")
	}

	minervaReadToken := consul.GetString("infinity/minerva/read-token", "read")
	minervaWriteToken := consul.GetString("infinity/minerva/write-token", "write")

	app := &infinity.Application{
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
		MinervaClient: minerva.NewClient(minervaConn, hostname, minervaReadToken, minervaWriteToken),

		// todo: get this from nomad
		HttpBindAddr: fmt.Sprintf("0.0.0.0:%s", getEnv("NOMAD_PORT_http", "9000")),
		GrpcBindAddr: fmt.Sprintf("0.0.0.0:%s", getEnv("NOMAD_PORT_grpc", "9001")),

		// JWT Session
		SessionSignKey:       []byte(consul.GetString("infinity/jwt-sign-key", "helloWorld")),
		SessionSignMethod:    jwt.SigningMethodHS256,
		SessionTokenDuration: time.Hour,

		ActiveSessionCollection: &infinity.ActiveSessionCollection{
			Mtx:      deadlock.RWMutex{},
			Sessions: map[uuid.UUID]*rpc.Socket{},
		},
	}

	app.PaymentService = payment.NewService(app)
	app.PaymentPlanRepository = payment.NewPaymentPlanRepository(app)
	app.PaymentTransactionRepository = payment.NewPaymentTransactionRepository(app)

	app.UserRepository = user.NewRepository(app)
	app.UserSessionService = user.NewSessionService(app)

	app.PerformerService = performer.NewService(app)
	app.PerformerRepository = performer.NewRepository(app)
	app.PerformerSearchService = performer.NewElasticSearchService(app)

	app.BandwidthConsumptionCollector = bandwidth.NewConsumptionCollector(app)
	app.BandwidthConsumptionRepository = bandwidth.NewConsumptionRepository(app)

	app.CryptoExchangeRateService = crypto.NewCryptoExchangeRateService(app)
	app.RecordingService = recording.NewService(app)
	app.RecordingRepository = recording.NewRepository(app)

	app.CleanupService = cleanup.NewCleanupService(app)
	app.CleanupRepository = cleanup.NewCleanupRepository(app)

	app.EmailService = email.NewEmailService(
		app,
		consul.GetString("infinity/email/host", ""),
		int(consul.GetInt16("infinity/email/port", 465)),
		consul.GetString("infinity/email/username", ""),
		consul.GetString("infinity/email/password", ""),
		consul.GetString("infinity/email/fromAddress", ""),
	)

	app.BlockcypherClient = gobcy.API{
		Token: "f7cd8fdbd3164c5686cc50bcd8de2f0e",
		Coin:  "btc",
		Chain: "main",
	}

	// Load page cache from redis
	go app.RecordingRepository.LoadPageCacheFromRedis()

	// Get exchange rates for btc/usd
	go app.CryptoExchangeRateService.Run(ctx)

	// Rebuild bandwidth consumption collector & start the database sync
	go app.BandwidthConsumptionCollector.RebuildFromToday()
	go app.BandwidthConsumptionCollector.Synchronizer(ctx)

	go app.CleanupService.Run(ctx)
	go app.RecordingService.Run(ctx)
	go app.PaymentService.Run(ctx)
	go app.PerformerSearchService.Run(ctx)

	// Start out services
	go http.StartHttpServer(ctx, app)
	go grpc.StartGrpcServer(ctx, app)

	signalCh := make(chan os.Signal, 1)

	// Listen for shutdown
	signal.Notify(signalCh, os.Interrupt)

	for {
		select {
		case <-signalCh:
			app.Logger.Info("Received SIGINT, terminating...")

			// cleanup
			cancel()
		}
	}
}

func migrate(host, port, user, password, dbname string) {

	url := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	conn, err := sqlx.Connect("postgres", url)
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to postgres to perform migrations")
	}

	defer conn.Close()

	mig := migrations.CreateFromDirectory(consul.GetString("infinity/migrations-path", "./migrations"))
	mig.Up(driver.NewPostgresDriver(conn), true)
}

func getEnv(key, defaultValue string) string {
	if os.Getenv(key) == "" {
		return defaultValue
	}

	return os.Getenv(key)
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
