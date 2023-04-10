package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain/minerva"
	"git.misc.vee.bz/carnagel/go-ecosystem/infrastructure/factory"
	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/domain/recording"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/infrastructure/redis"
	"github.com/go-pg/pg"
	"github.com/hashicorp/consul/api"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"gopkg.in/olivere/elastic.v5"
	"os"
	"os/signal"
)

var (
	ctx, cancel         = context.WithCancel(context.Background())
	doneCh              = make(chan bool)
	recordingsProcessed uint64
	hostname            string
	consul              ecosystem.ConsulClient
	logger              logrus.FieldLogger
	amqpCh              *amqp.Channel

	app *infinity.Application
)

var (
	command       = flag.String("command", "", "Command to run")
	minExternalId = flag.Uint64("minExternalId", 0, "Min external id offset")
	maxRecordings = flag.Uint64("maxRecordings", 0, "Max amount of recordings to process")
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

	minervaConn, err := factory.NewMinervaConnection(ctx, consul)
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to minerva")
	}

	minervaReadToken := consul.GetString("infinity/minerva/read-token", "read")
	minervaWriteToken := consul.GetString("infinity/minerva/write-token", "write")

	app = &infinity.Application{
		Amqp:          amqpConn,
		Logger:        logger,
		Redis:         redis,
		DB:            database,
		Consul:        consul,
		ElasticSearch: elasticClient,

		// Ecosystem clients
		MinervaClient: minerva.NewClient(minervaConn, hostname, minervaReadToken, minervaWriteToken),
	}

	app.RecordingRepository = recording.NewRepository(app)

	amqpCh, err = ecosystem.GetAmqpChannel(app.Amqp, 0)
	if err != nil {
		app.Logger.WithError(err).Fatal("Failed to create amqp channel")
	}
}

func main() {
	switch *command {
	case "validate-recordings":
		go processAllRecordings(validateRecordings)

	case "regenerate-images":
		go processAllRecordings(regenerateImages)

	default:
		logger.WithField("command", *command).Warn("Unknown command provided")

		return
	}

	signalCh := make(chan os.Signal, 1)

	// Listen for shutdown
	signal.Notify(signalCh, os.Interrupt)

	select {
	case success := <-doneCh:
		app.Logger.
			WithFields(logrus.Fields{"success": success, "recordingsProcessed": recordingsProcessed}).
			Info("Finished running command")
		return

	case <-signalCh:
		app.Logger.Info("Received SIGINT, terminating...")

		// cleanup
		cancel()
	}

	amqpCh.Close()
}

func validateRecordings(recordings []infinity.Recording) {
	ids := make([]ecosystem.ExternalId, len(recordings))
	for i, recording := range recordings {
		ids[i] = ecosystem.ExternalId(recording.ExternalId)
	}

	result, err := app.MinervaClient.HasBulk(ids)
	if err != nil {
		logger.WithError(err).Error("Failed to check if recordings exist while querying minerva")
		return
	}

	for _, recording := range recordings {
		log := logger.WithFields(logrus.Fields{
			"recordingUuid":       recording.Uuid,
			"recordingExternalId": recording.ExternalId,
		})

		exists, ok := result[ecosystem.ExternalId(recording.ExternalId)]
		if !ok {
			log.Warn("Did not receive existence status of recording from Minerva")

			continue
		}

		if exists {
			continue
		}

		log.Debug("Deleting recording")

		// Delete it in our repository first to avoid in premium user collection / view errors
		// when calling delete from coordinator
		_, err := app.RecordingRepository.Remove(recording.Uuid)
		if err != nil {
			log.WithError(err).Error("Failed to delete recording from database")

			continue
		}

		body, err := json.Marshal(&ecosystem.SystemScopedRecordingPayload{
			RecordingId: uint64(recording.ExternalId),
			System:      ecosystem.SystemInfinity,
		})

		if err != nil {
			log.WithError(err).Error("Failed to encode amqp payload")

			continue
		}

		payload := amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		}

		if err := amqpCh.Publish("", ecosystem.AmqpQueueRemoveRecording, false, false, payload); err != nil {
			log.WithError(err).Errorf("Failed to publish to queue %s", ecosystem.AmqpQueueRemoveRecording)
		}
	}
}

func regenerateImages(recordings []infinity.Recording) {
	for _, recording := range recordings {
		log := logger.WithFields(logrus.Fields{
			"recordingUuid":       recording.Uuid,
			"recordingExternalId": recording.ExternalId,
		})

		if recording.IsMinified() {
			log.Debug("Recording is already minified")

			continue
		}

		body, err := json.Marshal(&ecosystem.SystemScopedRecordingPayload{
			RecordingId: uint64(recording.ExternalId),
			System:      ecosystem.SystemInfinity,
		})

		if err != nil {
			log.WithError(err).Error("Failed to encode amqp payload")

			continue
		}

		payload := amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		}

		if err := amqpCh.Publish("", ecosystem.AmqpQueueImageRegeneration, false, false, payload); err != nil {
			log.WithError(err).Errorf("Failed to publish to queue %s", ecosystem.AmqpQueueImageRegeneration)
		}
	}
}

func processAllRecordings(processFunc func([]infinity.Recording)) {
	criteria := &infinity.RecordingRepositoryCriteria{
		SortMode:      infinity.SortModeExternalIdAsc,
		ExternalIdMin: *minExternalId,
		Limit:         90,
	}

	logger.WithFields(logrus.Fields{
		"maxRecordings": *maxRecordings,
		"minExternalId": *minExternalId,
	}).Info("Processing recordings")

	for {
		select {
		case <-ctx.Done():
			return

		default:
			recordings, _, err := app.RecordingRepository.Matching(criteria)
			if err != nil {
				logger.WithError(err).Error("Failed to retrieve recordings by criteria")
				doneCh <- false
				return
			}

			// Check if we ran it with a max number of recordings
			if *maxRecordings != 0 {
				// If we will surpass the max number of recordings, cut the array short
				if recordingsProcessed+uint64(len(recordings)) > *maxRecordings {
					logger.Debug("Cutting array short due to passed max number of recordings")

					limit := *maxRecordings - recordingsProcessed

					// Create temporary shorter array
					tmp := make([]infinity.Recording, limit)

					// Move recordings to temporary array
					for i := range tmp {
						tmp[i] = recordings[i]
					}

					// Set recordings to shortened temporary array
					recordings = tmp
				}
			}

			// If we cut it so much we have no recordings left, stop
			if len(recordings) == 0 {
				doneCh <- true
				return
			}

			// Call the processor function
			processFunc(recordings)

			// Add recordings processed
			recordingsProcessed += uint64(len(recordings))

			// Shortened arrays, due to max recordings to process, will also fall under this condition
			if len(recordings) < 90 {
				doneCh <- true
				return
			}

			criteria.ExternalIdMin = recordings[len(recordings)-1].ExternalId + 1
		}
	}
}

func getElasticSearchClient() (*elastic.Client, error) {
	services, _, err := consul.API().Catalog().Service("elasticsearch", "", &api.QueryOptions{})
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
