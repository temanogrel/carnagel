package encoder

import (
	"strconv"
	"strings"

	"context"

	"github.com/pkg/errors"
)

var (
	MissingVideoStreamErr        = errors.New("Could not find the video stream")
	RecordingDurationTooShortErr = errors.New("The duration of the recording is too short")
)

var DefaultConfig = &EncoderConfig{
	Crf:           "20",
	Tune:          "film",
	Threads:       "8",
	Preset:        "veryfast",
	FontFile:      "files/DroidSansMono.ttf",
	WatermarkFile: "files/watermark.png",
}

type VideoMeta struct {
	Duration  float64
	BitRate   int
	Size      uint64
	FrameRate string

	VideoCodecName   string
	VideoPixelFormat string
	VideoWidth       int
	VideoHeight      int

	HasAudio           bool
	AudioSampleRate    uint
	AudioCodecName     string
	AudioChannelLayout string
}

func (meta *VideoMeta) FrameRateAsInt() int {
	parts := strings.Split(meta.FrameRate, "/")
	if len(parts) != 2 {
		return -1
	}

	first, err := strconv.ParseUint(parts[0], 10, 8)
	if err != nil {
		return -1
	}

	second, err := strconv.ParseUint(parts[1], 10, 8)
	if err != nil {
		return -1
	}

	return int(first / second)
}

type EncoderConfig struct {
	//  service preferences
	OutputToConsole      bool
	HardwareAcceleration string

	// encoder preferences
	Crf     string
	Tune    string
	Preset  string
	Threads string

	// File paths
	FontFile      string
	WatermarkFile string
}

// EncoderService is a utility interface that wraps around the ffmpeg binary
// and in future versions could be replaced by other systems
type EncoderService interface {
	// Probe uses ffprobe to get as much meta information as possible
	Probe(ctx context.Context, source string) (*VideoMeta, error)

	// EncodeToH264 Decodes the provided file and encodes it into h264
	EncodeToH264(ctx context.Context, source, target string, meta *VideoMeta) error

	// EncodeToH265 Decodes the provided file and encodes it into h265
	EncodeToH265(ctx context.Context, source, target string, meta *VideoMeta) error

	// Encode an mp4 file to HLS
	EncodeMp42Hls(ctx context.Context, source, target, manifest string) error

	// Encode a HLS to mp4
	EncodeHls2Mp4(ctx context.Context, source, target, manifest string) error

	// GenerateInfinityThumbs generates a bunch of thumbs to be listed under the recordings
	// Note, clean up should be handled post calling defer os.RemoveAll on the second return argument
	GenerateInfinityThumbs(ctx context.Context, source string, interval uint) ([]string, string, error)

	// GenerateInfinitySprites generates 25x25 grids of images for the video scrub
	// Note, clean up should be handled post calling defer os.RemoveAll on the second return argument
	GenerateInfinitySprites(ctx context.Context, source string) ([]string, string, error)

	// GenerateInfinityCollage generates a 3x3 grid of images to the target path
	GenerateInfinityCollage(ctx context.Context, source, target string) error

	// GenerateWordpressCollage generates a 5x5 grid of images with a header to be published
	// on the wordpress sites
	GenerateWordpressCollage(ctx context.Context, source, target, filename string) error
}
