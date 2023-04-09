package main

import (
	"context"
	"fmt"
	"os"

	"os/signal"

	"git.misc.vee.bz/carnagel/encoder/pkg"
	"git.misc.vee.bz/carnagel/encoder/pkg/ffmpeg"
	"git.misc.vee.bz/carnagel/encoder/pkg/http"
	"git.misc.vee.bz/carnagel/encoder/pkg/image_generation"
	"git.misc.vee.bz/carnagel/encoder/pkg/pipeline"
	"git.misc.vee.bz/carnagel/encoder/pkg/worker"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain/aphrodite"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain/minerva"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain/upstore"
	"git.misc.vee.bz/carnagel/go-ecosystem/infrastructure/factory"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var (
	hostname string
	consul   ecosystem.ConsulClient
	logger   logrus.FieldLogger
)

func init() {

	var err error

	consul, logger, hostname, err = factory.Bootstrap("encoder")
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

	amqpConn, err := connectToRabbitmq()
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to rabbitmq")
	}

	minervaReadToken := consul.GetString("encoder/minerva/read-token", "read")
	minervaWriteToken := consul.GetString("encoder/minerva/write-token", "write")

	aphroditeUri := consul.GetString("encoder/aphrodite/api", "http://api.aphrodite.vee.bz")
	aphroditeToken := consul.GetString("encoder/aphrodite/token", "helloWorld")

	upstoreEmail := consul.GetString("upstore/email", "")
	upstorePassword := consul.GetString("upstore/password", "")

	encoderConfig := encoder.DefaultConfig
	encoderConfig.OutputToConsole = false

	app := &encoder.Application{
		Amqp:                   amqpConn,
		Logger:                 logger,
		ImageGenerationService: image_generation.NewVcsImageGeneration(logger),
		EncoderService:         ffmpeg.NewEncoderService(encoderConfig, logger),

		// Clients
		MinervaClient:   minerva.NewClient(minervaConn, hostname, minervaReadToken, minervaWriteToken),
		AphroditeClient: aphrodite.NewClient(aphroditeUri, aphroditeToken, logger),
		UpstoreClient:   upstore.NewUpstoreClient(upstoreEmail, upstorePassword, logger),

		Hostname:     hostname,
		HttpBindAddr: os.Getenv("NOMAD_ADDR_http"),

		MinimumDuration: 120,                    // At least two minutes
		H265Threshold:   1024 * 1024 * 1024 * 5, // 5GB
	}

	app.Pipelines = pipeline.NewComposedPipelines(app)
	app.Bootstrap()

	if cwd, err := os.Getwd(); err == nil {
		logger.WithField("cwd", cwd).Info("Starting the encoder")
	} else {
		logger.WithError(err).Fatal("Could not retrive working directory")
	}

	go worker.NewEncodingWorker(app).Run(ctx, 6)
	go worker.NewMp42hlsWorker(app).Run(ctx, 1)
	go worker.NewImageRegenerationWorker(app).Run(ctx, 1)
	go http.StartServer(ctx, app)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	<-signalCh

	app.Logger.Infof("Received SIGINT, terminating...")

	// cancel the "root" context, this will tell everything to terminate
	cancel()
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

	user := consul.GetString("encoder/rabbitmq/user", "encoder")
	pass := consul.GetString("encoder/rabbitmq/pass", "encoder")
	vhost := consul.GetString("encoder/rabbitmq/vhost", "encoder")

	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", user, pass, instances[0].ServiceAddress, instances[0].ServicePort, vhost)

	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to connect to rabbitmq service")
	}

	return conn, nil
}
