package process

import (
	"context"
	"fmt"
	"git.misc.vee.bz/carnagel/encoder/pkg"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"io/ioutil"
	"os"
)

func ProbeFlv(ctx context.Context, app *encoder.Application, rec *ecosystem.Recording) (context.Context, error) {
	log := log(ctx, app)
	log.Debug("Running ProbeFlv")

	path, ok := ctx.Value(encoder.CtxFlvFile).(string)
	if !ok {
		return ctx, encoder.MissingFlvFileErr
	}

	meta, err := app.EncoderService.Probe(ctx, path)
	if err != nil {
		if err == encoder.MissingVideoStreamErr {
			log.Info("Deleting due to missing video stream")
		} else {
			log.
				WithError(err).
				Warn("Failed to do an initial probe, file will be deleted")
		}

		return ctx, encoder.ShortCircuitAndDeleteRecordingErr
	}

	ctx = context.WithValue(ctx, encoder.CtxProbeMeta, meta)

	return ctx, nil
}

func Probe(ctx context.Context, app *encoder.Application, recording *ecosystem.Recording) (context.Context, error) {
	log := log(ctx, app)
	log.Debug("Running probe")

	path, ok := ctx.Value(encoder.CtxMp4File).(string)
	if !ok {
		return ctx, encoder.MissingMp4FileErr
	}

	meta, err := app.EncoderService.Probe(ctx, path)
	if err != nil {
		if err == encoder.MissingVideoStreamErr {
			log.Info("Deleting due to missing video stream")
		} else {
			log.
				WithError(err).
				Warn("Failed to do an initial probe, file will be deleted")
		}

		return ctx, encoder.ShortCircuitAndDeleteRecordingErr
	}

	// store the video meta
	log = log.WithField("meta", meta)

	if meta.Duration <= app.MinimumDuration {
		log.Info("Deleting due to too short duration")

		return ctx, encoder.ShortCircuitAndDeleteRecordingErr
	}

	if meta.HasAudio {
		recording.Audio = fmt.Sprintf("%s, %d Hz, %s", meta.AudioCodecName, meta.AudioSampleRate, meta.AudioChannelLayout)
	}

	// Updating the size after encoding
	if meta.VideoCodecName == "h264" {
		recording.Size264 = meta.Size
	} else {
		recording.Size265 = meta.Size
	}

	recording.Video = fmt.Sprintf("%s, %s, %dx%d, %d fps", meta.VideoCodecName, meta.VideoPixelFormat, meta.VideoWidth, meta.VideoHeight, meta.FrameRateAsInt())
	recording.BitRate = uint32(meta.BitRate)
	recording.Duration = meta.Duration
	recording.Encoding = meta.VideoCodecName

	ctx = context.WithValue(ctx, encoder.CtxProbeMeta, meta)

	return ctx, nil
}

func VideoEncodeFlv2Mp4(ctx context.Context, app *encoder.Application, rec *ecosystem.Recording) (context.Context, error) {
	log := log(ctx, app)
	log.Debug("Running video encode flv2mp4")

	meta, ok := ctx.Value(encoder.CtxProbeMeta).(*encoder.VideoMeta)
	if !ok {
		return ctx, encoder.MissingProbeErr
	}

	sourcePath, ok := ctx.Value(encoder.CtxFlvFile).(string)
	if !ok {
		return ctx, encoder.MissingFlvFileErr
	}

	targetPath := encoder.ExtendFileName(sourcePath, "", "mp4")

	if err := app.EncoderService.EncodeToH264(ctx, sourcePath, targetPath, meta); err != nil {
		return ctx, encoder.ShortCircuitAndDeleteRecordingErr
	}

	ctx = context.WithValue(ctx, encoder.CtxMp4File, targetPath)

	return ctx, nil
}

func VideoEncodeHls2Mp4(ctx context.Context, app *encoder.Application, rec *ecosystem.Recording) (context.Context, error) {
	log := log(ctx, app)
	log.Debug("Running video encode hls2mp4")

	sourcePath, ok := ctx.Value(encoder.CtxHlsFile).(string)
	if !ok {
		return ctx, encoder.MissingHlsFileErr
	}

	mp4Path := encoder.ExtendFileName(sourcePath, "", "mp4")

	if err := app.EncoderService.EncodeHls2Mp4(ctx, sourcePath, mp4Path, rec.VideoManifest); err != nil {
		return ctx, err
	}

	ctx = context.WithValue(ctx, encoder.CtxMp4File, mp4Path)

	return ctx, nil
}

func VideoEncodeMp42Hls(ctx context.Context, app *encoder.Application, rec *ecosystem.Recording) (context.Context, error) {
	log := log(ctx, app)
	log.Debug("Running video encode mp42hls")

	source, ok := ctx.Value(encoder.CtxMp4File).(string)
	if !ok {
		return ctx, encoder.MissingMp4FileErr
	}

	hlsVideoFile := encoder.ReplaceFileExtension(source, "ts")
	hlsManifestFile := encoder.ReplaceFileExtension(source, "m3u8")

	defer os.Remove(hlsManifestFile)

	if err := app.EncoderService.EncodeMp42Hls(ctx, source, hlsVideoFile, hlsManifestFile); err != nil {
		return ctx, err
	}

	ctx = context.WithValue(ctx, encoder.CtxHlsFile, hlsVideoFile)

	manifest, err := ioutil.ReadFile(hlsManifestFile)
	if err != nil {
		return ctx, err
	}

	rec.VideoManifest = string(manifest)

	return ctx, nil
}
