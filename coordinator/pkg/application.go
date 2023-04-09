package coordinator

import (
	"context"

	"database/sql"

	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"git.misc.vee.bz/carnagel/go-ecosystem/protobuf/common"
	"github.com/go-pg/pg"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type Application struct {
	Ctx context.Context

	Amqp   *amqp.Connection
	Logger logrus.FieldLogger
	Consul ecosystem.ConsulClient

	// Http api client
	UltronClient         ecosystem.UltronClient
	HermesClient         ecosystem.HermesClient
	MinervaClient        ecosystem.MinervaClient
	AphroditeClient      ecosystem.AphroditeClient
	CamgirlGalleryClient ecosystem.CamgirlGalleryClient
	ProxyServerClient    ecosystem.ProxyServerClient

	// Grpc clients
	InfinityContentClient common.ContentClient

	// Content
	DeathFileService  DeathFileService
	CleanupService    CleanupService
	ContentService    ContentService
	ContentSubscriber ContentSubscriber

	// Scraper
	Scraper             Scraper
	MyfreeCamsScraper   MyfreeCamsScraper
	ScraperProxyService ScraperProxyService

	UltronDb    *sql.DB
	HermesDb    *sql.DB
	AphroditeDb *sql.DB
	MinervaDb   *pg.DB
}
