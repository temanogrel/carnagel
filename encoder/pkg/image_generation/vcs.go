package image_generation

import (
	"bytes"
	"context"
	"io"
	"os/exec"
	"strings"

	"git.misc.vee.bz/carnagel/encoder/pkg"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type vcs struct {
	logger logrus.FieldLogger
}

const executable = "/usr/bin/vcs"

func createBasicArgs(source, target string) []string {
	return []string{
		// Input file
		source,

		// Disable shadows
		"-ds",
		// Disable padding
		"-dp",

		// Custom options
		"-O", "quality=75",

		// Output as jpeg
		"-j",

		// Output file
		"-o", target,
	}
}

func NewVcsImageGeneration(logger logrus.FieldLogger) encoder.ImageGenerationService {
	return &vcs{logger.WithField("service", "ImageGenerationService")}
}

func (s *vcs) InfinityThumbs(ctx context.Context, source, target string) error {
	args := []string{
		// Number of images
		"-n", "80",

		// Number of columns (20x4)
		"-c", "4",

		// Image height
		"-H150",
		
		// Set aspect ratio to 16/9
		"-a", "31/20",

		// Disable timestamps
		"-dt",

		// Disable the meta info (header / footer)
		"-dmetainfo",
	}

	if _, err := s.exec(ctx, source, target, args...); err != nil {
		return errors.Wrap(err, "Failed to generate thumbnails from video")
	}

	return nil
}

func (s *vcs) InfinitySprites(ctx context.Context, source, target string) error {
	args := []string{
		// Interval of images
		"-i", "5s",

		// Number of columns (11xY)
		"-c", "11",

		// Image height
		"-H65",
		
		// Set aspect ratio to 16/9
		"-a", "31/20",

		"-E", "0%",

		// Disable timestamps
		"-dt",

		// End offset 0% (generate for the entire duration of the video)
		"-E", "0%",

		// Disable the meta info (header)
		"-dmetainfo",
	}

	if _, err := s.exec(ctx, source, target, args...); err != nil {
		return errors.Wrap(err, "Failed to generate sprites from video")
	}

	return nil
}

func (s *vcs) InfinityCollage(ctx context.Context, source, target string) error {
	args := []string{
		// Number of images
		"-n", "9",

		// Number of columns (3x3)
		"-c", "3",

		// Image height
		"-H400",
		
		// Set aspect ratio to Auto
		"-A",

		// Disable timestamps
		"-dt",

		// Disable the meta info (header)
		"-dmetainfo",
	}

	if _, err := s.exec(ctx, source, target, args...); err != nil {
		return errors.Wrap(err, "Failed to generate collage from video")
	}

	return nil
}

func (s *vcs) WordpressCollage(ctx context.Context, source, target, filename string) error {
	args := []string{
		// Number of images
		"-n", "16",

		// Number of columns (4x4)
		"-c", "4",

		// Image height
		"-H300",
		
		// Set aspect ratio to Auto
		"-A",

		// Custom options
		"-O", "bg_all=black",
		"-O", "fg_heading=yellow",
		"-O", "fg_sign=black",
		"-O", "fg_tstamps=yellow",
		"-O", "fg_title=yellow",

		// Custom filename in header
		"--filename", filename,

		"-U0",
	}

	if _, err := s.exec(ctx, source, target, args...); err != nil {
		return errors.Wrap(err, "Failed to generate collage from video")
	}

	return nil
}

func (s *vcs) exec(ctx context.Context, source, target string, args ...string) (io.Reader, error) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	log := ecosystem.NewLogEntry(ctx, s.logger, logrus.Fields{
		"bin":  executable,
		"args": strings.Join(append(createBasicArgs(source, target), args...), " "),
	})

	cmd := exec.Command(executable, append(createBasicArgs(source, target), args...)...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		log.
			WithFields(logrus.Fields{
				"stdout": stdout.String(),
				"stderr": stderr.String(),
			}).
			Error("Failed to execute command")

		return nil, errors.Wrap(err, "Failed to execute command")
	}

	log.Debug("Command executed successfully")

	return stdout, nil
}
