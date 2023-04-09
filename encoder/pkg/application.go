package encoder

import (
	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type Application struct {
	Amqp                   *amqp.Connection
	Logger                 logrus.FieldLogger
	EncoderService         EncoderService
	ImageGenerationService ImageGenerationService
	Pipelines              ComposedPipelines

	// Api clients
	MinervaClient   ecosystem.MinervaClient
	AphroditeClient ecosystem.AphroditeClient
	UpstoreClient   ecosystem.UpstoreClient

	// General config
	Hostname     string
	HttpBindAddr string

	// Encoding settings
	H265Threshold   uint64
	MinimumDuration float64
}

func (app *Application) ensureProductionSettings() {

	if app.Hostname == "" {
		app.Logger.Fatal("No hostname specified")
	}

	if app.H265Threshold == 0 {
		app.Logger.Fatal("No H265Treshold specified")
	}

	if app.MinimumDuration == 0.0 {
		app.Logger.Fatal("No MinimumDuration specified")
	}
}

func (app *Application) Bootstrap() {
	app.ensureProductionSettings()

	ch, err := ecosystem.GetAmqpChannel(app.Amqp, 0)
	if err != nil {
		app.Logger.WithError(err).Fatal("Failed to create amqp channel")
	}

	ch.Close()
}
