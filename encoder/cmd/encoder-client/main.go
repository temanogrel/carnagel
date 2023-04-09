package main

import (
	"context"
	"flag"
	"git.misc.vee.bz/carnagel/encoder/pkg"
	"git.misc.vee.bz/carnagel/encoder/pkg/ffmpeg"
	"git.misc.vee.bz/carnagel/encoder/pkg/image_generation"
	"git.misc.vee.bz/carnagel/encoder/pkg/pipeline"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain"

	"git.misc.vee.bz/carnagel/go-ecosystem/domain/aphrodite"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain/minerva"
	"git.misc.vee.bz/carnagel/go-ecosystem/infrastructure/factory"
	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

var (
	source    = flag.String("source", "", "Source file")
	target    = flag.String("target", "", "Target file")
	command   = flag.String("command", "", "Command to execute")
	recording = flag.Uint64("recording", 0, "recording id")

	application *encoder.Application
)

var (
	ctx      context.Context
	hostname string
	consul   ecosystem.ConsulClient
	logger   logrus.FieldLogger
)

func init() {
	flag.Parse()

	var err error

	consul, logger, hostname, err = factory.Bootstrap("encoder")
	if err != nil {
		panic(err)
	}

	ctx = context.Background()

	minervaConn, err := factory.NewMinervaConnection(ctx, consul)
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to minerva")
	}

	aphroditeUri := consul.GetString("encoder/aphrodite/api", "http://api.aphrodite.vee.bz")
	aphroditeToken := consul.GetString("encoder/aphrodite/token", "helloWorld")

	minervaReadToken := consul.GetString("encoder/minerva/read-token", "read")
	minervaWriteToken := consul.GetString("encoder/minerva/write-token", "write")

	config := encoder.DefaultConfig
	config.OutputToConsole = true

	application = &encoder.Application{
		Logger:                 logger,
		EncoderService:         ffmpeg.NewEncoderService(config, logger),
		ImageGenerationService: image_generation.NewVcsImageGeneration(logger),

		// Clients
		MinervaClient:   minerva.NewClient(minervaConn, hostname, minervaReadToken, minervaWriteToken),
		AphroditeClient: aphrodite.NewClient(aphroditeUri, aphroditeToken, logger),

		Hostname: hostname,

		MinimumDuration: 120,                    // At least two minutes
		H265Threshold:   1024 * 1024 * 1024 * 5, // 5GB
	}

	application.Pipelines = pipeline.NewComposedPipelines(application)
}

func main() {

	switch *command {
	case "probe":
		if meta, err := application.EncoderService.Probe(ctx, *source); err == nil {
			spew.Dump(meta)
		} else {
			application.Logger.WithError(err).Fatal("Failed to retrieve settings")
		}

	case "encode-264":
		meta, err := application.EncoderService.Probe(ctx, *source)
		if err != nil {
			application.Logger.WithError(err).Fatal("Failed to retrieve settings")
		}

		if err := application.EncoderService.EncodeToH264(ctx, *source, *target, meta); err != nil {
			application.Logger.WithError(err).Fatal("Encoding failed")
		}

	case "mp4-to-hls":
		if err := application.EncoderService.EncodeMp42Hls(ctx, *source, "test/video.ts", "test/stream.m3u8"); err != nil {
			application.Logger.WithError(err).Fatal("Encoding failed")
		}

	case "wordpress-collage":
		if err := application.ImageGenerationService.WordpressCollage(ctx, *source, *target, filepath.Base(*source)); err != nil {
			application.Logger.WithError(err).Fatal("Wordpress collage generation failed")
		}

	case "infinity-collage":
		if err := application.ImageGenerationService.InfinityCollage(ctx, *source, *target); err != nil {
			application.Logger.WithError(err).Fatal("Infinity collage generation failed")
		}

	case "infinity-thumbs":
		if err := application.ImageGenerationService.InfinityThumbs(ctx, *source, *target); err != nil {
			application.Logger.WithError(err).Fatal("Infinity thumbs generation failed")
		}

	case "infinity-sprites":
		if err := application.ImageGenerationService.InfinitySprites(ctx, *source, *target); err != nil {
			application.Logger.WithError(err).Fatal("Infinity sprites genertion failed")
		}

	case "mp42hls":
		rec, err := application.AphroditeClient.GetRecording(*recording)
		if err != nil {
			application.Logger.
				WithError(err).
				Error("Failed to retrieve recording")

			return
		}

		if err := application.Pipelines.ConvertMp42Hls(ctx, rec); err != nil {
			application.Logger.
				WithError(err).
				Error("Failed to convert mp42hls")
		}

	default:
		application.Logger.WithField("command", *command).Fatal("Invalid command provided")
	}
}
