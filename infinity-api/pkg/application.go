package infinity

import (
	"time"

	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-pg/pg"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"gopkg.in/olivere/elastic.v5"

	"github.com/blockcypher/gobcy"
)

type Application struct {
	// Infrastructure
	Amqp          *amqp.Connection
	DB            *pg.DB
	Redis         *redis.Client
	Consul        ecosystem.ConsulClient
	Logger        logrus.FieldLogger
	ElasticSearch *elastic.Client

	// Ecosystem clients
	MinervaClient   ecosystem.MinervaClient
	AphroditeClient ecosystem.AphroditeClient

	// Config
	HttpBindAddr string
	GrpcBindAddr string

	// Auth
	SessionSignKey       []byte
	SessionSignMethod    jwt.SigningMethod
	SessionTokenDuration time.Duration

	// Performer
	PerformerService       PerformerService
	PerformerRepository    PerformerRepository
	PerformerSearchService PerformerSearchService

	// Recording
	RecordingService    RecordingService
	RecordingRepository RecordingRepository

	// Payment
	PaymentService               PaymentService
	PaymentTransactionRepository PaymentTransactionRepository
	PaymentPlanRepository        PaymentPlanRepository

	// Bandwidth
	BandwidthConsumptionCollector  BandwidthConsumptionCollector
	BandwidthConsumptionRepository BandwidthConsumptionRepository

	// Session
	ActiveSessionCollection *ActiveSessionCollection

	// User
	UserRepository     UserRepository
	UserSessionService UserSessionService

	// Bitcoin stuff
	BlockcypherClient         gobcy.API
	CryptoExchangeRateService CryptoExchangeRateService

	// Email
	EmailService EmailService

	// Cleanup
	CleanupService    CleanupService
	CleanupRepository CleanupRepository
}

// EnsureValidSettings will make sure that the application is in a valid state to start
func (app *Application) EnsureValidSettings() {
	if app.HttpBindAddr == "" {
		app.Logger.Fatal("HttpBindAddr not set")
	}

	if app.GrpcBindAddr == "" {
		app.Logger.Fatal("GrpcBindAddr not set")
	}
}
