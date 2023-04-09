package main

import (
	"context"
	"flag"
	"fmt"
	"git.misc.vee.bz/carnagel/encoder/pkg"
	"git.misc.vee.bz/carnagel/encoder/pkg/image_generation"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"strings"
)

var application *encoder.Application

var (
	ctx    = context.Background()
	logger = logrus.New()
)

const sampleVideosDir = "/sample-videos"

func init() {
	flag.Parse()

	application = &encoder.Application{
		Logger:                 logger,
		ImageGenerationService: image_generation.NewVcsImageGeneration(logger),
	}
}

func main() {
	files, err := ioutil.ReadDir(sampleVideosDir)
	if err != nil {
		application.Logger.WithError(err).Fatal("Failed to read sample video directory")
	}

	application.Logger.WithField("count", len(files)).Info("Found sample videos")

	for _, f := range files {
		application.Logger.WithField("filename", f.Name()).Info("Video found in sample video directory")
	}

	for _, f := range files {
		nameWithoutExtension := strings.Replace(f.Name(), ".mp4", "", -1)
		source := fmt.Sprintf("%s/%s", sampleVideosDir, f.Name())

		// Create sample images for all wp / infinity images using the sample-video found under files/sample-video.mp4
		if err := application.ImageGenerationService.WordpressCollage(ctx, source, fmt.Sprintf("/samples/%s-wp-collage.jpg", nameWithoutExtension), source); err != nil {
			application.Logger.WithError(err).Fatal("Wordpress collage sample generation failed")
		}

		if err := application.ImageGenerationService.InfinityCollage(ctx, source, fmt.Sprintf("/samples/%s-inf-collage.jpg", nameWithoutExtension)); err != nil {
			application.Logger.WithError(err).Fatal("Infinity collage sample generation failed")
		}

		if err := application.ImageGenerationService.InfinitySprites(ctx, source, fmt.Sprintf("/samples/%s-inf-sprites.jpg", nameWithoutExtension)); err != nil {
			application.Logger.WithError(err).Fatal("Infinity sprites sample generation failed")
		}

		if err := application.ImageGenerationService.InfinityThumbs(ctx, source, fmt.Sprintf("/samples/%s-inf-thumbs.jpg", nameWithoutExtension)); err != nil {
			application.Logger.WithError(err).Fatal("Infinity thumbs sample generation failed")
		}
	}
}
