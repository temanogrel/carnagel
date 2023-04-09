package ffmpeg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"io"

	"context"

	"git.misc.vee.bz/carnagel/encoder/pkg"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

const (
	wordpressCollageVideoFilters = `` +
		`scale=500:-1,` +
		`fps=1/{{ .FrameInterval }},` +
		`drawtext=fontfile='{{ .FontFile }}':fontsize=12:x=W-tw-10:y=H-th-10:fontcolor=white:text='%{pts\:hms}': box=0,tile=4x4:padding=4:margin=16:color=black@0.4,drawbox=y=0:x=0:color=#00000000:width=iw:height=16:t=max,` +
		`drawtext=fontfile='{{ .FontFile }}':text='Filename\: {{ .Filename }}':fontcolor=yellow:fontsize=14:x=16:y=1,drawbox=y=0:x=0:color=#00000000:width=iw:height=64:t=max,` +
		`drawtext=fontfile='{{ .FontFile }}':text='File size\: {{ .Size }}':fontcolor=yellow:fontsize=14:x=800:y=1,drawbox=y=0:x=0:color=#00000000:width=iw:height=64:t=max,` +
		`drawtext=fontfile='{{ .FontFile }}':text='Length\: {{ .Duration }}':fontcolor=yellow:fontsize=14:x=1100:y=1,drawbox=y=0:x=0:color=#00000000:width=iw:height=64:t=max,` +
		`drawtext=fontfile='{{ .FontFile }}':text='Dimensions\: {{ .Resolution }}':fontcolor=yellow:fontsize=14:x=1400:y=1,drawbox=y=0:x=0:color=#00000000:width=iw:height=64:t=max,` +
		`drawtext=fontfile='{{ .FontFile }}':text='Format\: {{ .Format }}':fontcolor=yellow:fontsize=14:x=1700:y=1,drawbox=y=0:x=0:color=#00000000:width=iw:height=64:t=max,` +
		`drawtext=fontfile='{{ .FontFile }}':text='FPS\: {{ .FrameRate }}':fontcolor=yellow:fontsize=14:x=W-tw-16:y=1`

	infinityCollageVideoFilters = `` +
		`scale=500:-1,` +
		`fps=1/{{ .FrameInterval }},` +
		`tile=3x3`
)

type ffmpegWrapper struct {
	log    logrus.FieldLogger
	config *encoder.EncoderConfig
}

func NewEncoderService(config *encoder.EncoderConfig, logger logrus.FieldLogger) encoder.EncoderService {
	return &ffmpegWrapper{
		log:    logger,
		config: config,
	}
}

func (service *ffmpegWrapper) exec(ctx context.Context, name string, args ...string) (io.Reader, error) {

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	log := ecosystem.NewLogEntry(ctx, service.log, logrus.Fields{
		"bin":  name,
		"args": strings.Join(args, " "),
	})

	cmd := exec.Command(name, args...)
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

	if service.config.OutputToConsole {
		log.Debug("Command executed successfully")
	}

	return stdout, nil
}

func (service *ffmpegWrapper) Probe(ctx context.Context, source string) (*encoder.VideoMeta, error) {
	args := []string{"-analyzeduration", "90M", "-probesize", "18M", "-show_format", "-show_streams", "-loglevel", "quiet", "-print_format", "json", source}

	stdout, err := service.exec(ctx, "ffprobe", args...)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to execute ffprobe command")
	}

	data, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read stdout from probe")
	}

	videoProbe := &ffprobeResult{}
	if err := json.Unmarshal(data, videoProbe); err != nil {
		return nil, errors.Wrap(err, "Failed to parse ffprobe output")
	}

	meta, err := ffprobeToVideoMeta(videoProbe)
	if err != nil {
		fmt.Println(string(data))
		service.log.
			WithField("ffprobeResult", string(data)).
			WithError(err).Error("Failed to parse ffprobeResult")

		return nil, err
	}

	return meta, nil
}

func (service *ffmpegWrapper) EncodeToH264(ctx context.Context, source, target string, meta *encoder.VideoMeta) error {

	args := []string{
		// Overwrite existing file
		"-y",

		// Input file
		"-i", source,

		// Output level
		"-loglevel", "24",

		// Video
		"-c:v", "libx264", "-crf", service.getCrfFromMeta(meta),

		// Audio
		"-c:a", "aac", "-preset", service.config.Preset,

		// ffmpeg fix for the following error
		//
		// Too many packets buffered for output stream 0:1.
		// [aac @ 0x3bc40a0] 2 frames left in the queue on closing
		"-max_muxing_queue_size", "9999",

		// Watermark the file
		"-vf", fmt.Sprintf("movie=%s [watermark]; [in][watermark] overlay=main_w-overlay_w-1:1 [out]", service.config.WatermarkFile),

		// Misc
		"-tune", service.config.Tune, "-threads", service.config.Threads,

		// Maxrate - set a maximum bitrate
		"-maxrate", "3M",

		// Bufsize - set a buffer size to control maxrate, usually 2 x maxrate
		"-bufsize", "6M",

		// Output file
		target,
	}

	if _, err := service.exec(ctx, "ffmpeg", args...); err != nil {
		return errors.Wrap(err, "Failed to encode to h264")
	}

	return nil
}

func (service *ffmpegWrapper) EncodeToH265(ctx context.Context, source, target string, meta *encoder.VideoMeta) error {

	args := []string{
		// Overwrite existing file
		"-y",

		// Input file
		"-i", source,

		// Output level
		"-loglevel", "24",

		// Video
		"-c:v", "libx265", "-crf", service.getCrfFromMeta(meta),

		// Audio
		"-c:a", "libfdk_aac", "-preset", service.config.Preset,

		// Watermark the file
		"-vf", fmt.Sprintf("movie=%s [watermark]; [in][watermark] overlay=main_w-overlay_w-1:1 [out]", service.config.WatermarkFile),

		// Misc
		"-tune", service.config.Tune, "-threads", service.config.Threads,

		// Output file
		target,
	}

	if _, err := service.exec(ctx, "ffmpeg", args...); err != nil {
		return errors.Wrap(err, "Failed to encode to h264")
	}

	return nil
}

func (service *ffmpegWrapper) EncodeMp42Hls(ctx context.Context, source, target, manifest string) error {
	args := []string{
		"--output-single-file",
		"--segment-filename-template", target,
		"--index-filename", manifest,
		"--iframe-index-filename", "/dev/null",
		source,
	}

	if _, err := service.exec(ctx, "mp42hls", args...); err != nil {
		return errors.Wrap(err, "Failed to encode mp4 to hls")
	}

	return nil
}

func (service *ffmpegWrapper) EncodeHls2Mp4(ctx context.Context, source, target, manifest string) error {
	path := fmt.Sprintf("/tmp/%s.m3u8", uuid.NewV4().String())
	f, err := os.Create(path)
	if err != nil {
		return errors.Wrapf(err, "Failed to create a temporary file for the manifest")
	}

	//defer os.Remove(path)

	f.WriteString(strings.Replace(manifest, "stream.ts", source, -1))
	f.Close()

	args := []string{
		"-i", source,
		"-c", "copy",
		target,
	}

	if _, err := service.exec(ctx, "ffmpeg", args...); err != nil {
		return errors.Wrapf(err, "Failed to encode the hls file to mp4")
	}

	return nil
}

func (service *ffmpegWrapper) GenerateInfinityThumbs(ctx context.Context, source string, interval uint) ([]string, string, error) {
	dir, err := ioutil.TempDir(os.TempDir(), "infinity_thumbs")
	if err != nil {
		return nil, "", errors.Wrap(err, "Failed to create temporary directory for generating collage")
	}

	args := []string{

		// Input file
		"-i", source,

		// Scale to 300 width and keep the same ratio
		"-vf", "scale=300:-1",

		// Interval at which to generate thumbs
		"-r", fmt.Sprintf("1/%d", interval),

		// Target
		filepath.Join(dir, "thumbs-%02d.jpg"),
	}

	if service.config.HardwareAcceleration != "" {
		args = append([]string{"-hwaccel", service.config.HardwareAcceleration}, args...)
	}

	if _, err := service.exec(ctx, "ffmpeg", args...); err != nil {
		return nil, dir, errors.Wrap(err, "Failed to generate thumbnails from video")
	}

	matches, err := filepath.Glob(filepath.Join(dir, "thumbs-*.jpg"))
	if err != nil {
		return nil, dir, errors.Wrap(err, "Failed to glob all the thumbs in the temporary folder")
	}

	return matches, dir, nil
}
func (service *ffmpegWrapper) GenerateInfinitySprites(ctx context.Context, source string) ([]string, string, error) {
	dir, err := ioutil.TempDir(os.TempDir(), "infinity_thumbs")
	if err != nil {
		return nil, "", errors.Wrap(err, "Failed to create temporary directory for generating collage")
	}

	args := []string{
		// Overwrite existing files if they exist
		"-y",

		// Input file
		"-i", source,

		// Scale the videos to the appropriate size
		"-vf", "scale=178:100,fps=1/5,tile=10x10",

		// Not sure, but I need it for the FPS filter
		"-vsync", "vfr", "-q:v", "1",

		// Output file
		filepath.Join(dir, "sprites-%05d.jpg"),
	}

	if service.config.HardwareAcceleration != "" {
		args = append([]string{"-hwaccel", service.config.HardwareAcceleration}, args...)
	}

	if _, err := service.exec(ctx, "ffmpeg", args...); err != nil {
		return nil, dir, errors.Wrap(err, "Failed to create infinity sprites")
	}

	matches, err := filepath.Glob(filepath.Join(dir, "sprites-*.jpg"))
	if err != nil {
		return nil, dir, errors.Wrap(err, "Failed to glob all the sprites in the temporary folder")
	}

	return matches, dir, nil
}

func (service *ffmpegWrapper) GenerateInfinityCollage(ctx context.Context, source, target string) error {
	probe, err := service.Probe(ctx, source)
	if err != nil {
		return errors.Wrap(err, "Failed to probe source file")
	}

	dir, err := ioutil.TempDir(os.TempDir(), "infinity_collage")
	if err != nil {
		return errors.Wrap(err, "Failed to create temporary directory for generating collage")
	}

	defer os.RemoveAll(dir)

	filterBuf := &bytes.Buffer{}
	filters := template.New("infinity collage filters")
	filters.Parse(infinityCollageVideoFilters)
	filters.Execute(filterBuf, struct {
		FrameInterval float64
	}{
		FrameInterval: probe.Duration / 8,
	})

	args := []string{
		// Input file
		"-i", source,

		// Scale the videos to the appropriate size
		"-vf", filterBuf.String(),

		// Not sure what this does, but it's required else it nine images of the same crap
		"-vsync", "vfr",

		// Output file
		target,
	}

	if service.config.HardwareAcceleration != "" {
		args = append([]string{"-hwaccel", service.config.HardwareAcceleration}, args...)
	}

	if _, err := service.exec(ctx, "ffmpeg", args...); err != nil {
		return errors.Wrap(err, "Failed to generate infinity collage")
	}

	return nil
}

func (service *ffmpegWrapper) GenerateWordpressCollage(ctx context.Context, source, target, filename string) error {
	probe, err := service.Probe(ctx, source)
	if err != nil {
		return errors.Wrap(err, "Failed to probe source file")
	}

	var format string
	if probe.HasAudio {
		format = fmt.Sprintf("%s / %s (%s)", probe.VideoCodecName, probe.AudioCodecName, probe.AudioChannelLayout)
	} else {
		format = fmt.Sprintf("%s / audio missing", probe.VideoCodecName)
	}

	size := humanize.Bytes(probe.Size)
	duration := durationToTimestamp(probe.Duration)
	resolution := fmt.Sprintf("%dx%d", probe.VideoWidth, probe.VideoHeight)

	filterBuf := &bytes.Buffer{}
	filters := template.New("video filters")
	filters.Parse(wordpressCollageVideoFilters)
	filters.Execute(filterBuf, struct {
		FontFile      string
		Filename      string
		Size          string
		Duration      string
		Resolution    string
		Format        string
		FrameInterval float64
		FrameRate     int
	}{
		FontFile:      service.config.FontFile,
		Filename:      filename,
		Size:          size,
		Duration:      strings.Replace(duration, ":", "\\:", 2),
		Resolution:    resolution,
		Format:        format,
		FrameRate:     probe.FrameRateAsInt(),
		FrameInterval: probe.Duration / 15,
	})

	args := []string{
		// Overwrite existing files if they exist
		"-y",

		// Input file
		"-i", source,

		// Scale the videos to the appropriate size
		"-vf", filterBuf.String(),

		// Not sure, but I need it for the FPS filter
		"-vsync", "vfr", "-q:v", "1",

		// Output file
		target,
	}

	if service.config.HardwareAcceleration != "" {
		args = append([]string{"-hwaccel", service.config.HardwareAcceleration}, args...)
	}

	if _, err := service.exec(ctx, "ffmpeg", args...); err != nil {
		return errors.Wrap(err, "Failed to create wordpress collage")
	}

	return nil
}

func (service *ffmpegWrapper) getCrfFromMeta(meta *encoder.VideoMeta) string {
	switch {

	// 200kbps
	case meta.BitRate < 200*1024:
		return "21"

		// 400kbps
	case meta.BitRate < 400*1024:
		return "21"

		// 1000kbps
	case meta.BitRate < 1000*1024:
		return "21"

		// 2000kbps
	case meta.BitRate < 2000*1024:
		return "21"

	default:
		return "21"
	}
}

func durationToTimestamp(duration float64) string {
	return time.Time{}.Add(time.Duration(duration) * time.Second).Format("15:04:05")
}
