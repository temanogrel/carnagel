package ffmpeg

import (
	"strconv"

	"git.misc.vee.bz/carnagel/encoder/pkg"
	"github.com/pkg/errors"
)

type stream struct {
	Index              int    `json:"index"`
	CodecName          string `json:"codec_name"`
	CodecLongName      string `json:"codec_long_name"`
	Profile            string `json:"profile,omitempty"`
	CodecType          string `json:"codec_type"`
	CodecTimeBase      string `json:"codec_time_base"`
	CodecTagString     string `json:"codec_tag_string"`
	CodecTag           string `json:"codec_tag"`
	Width              int    `json:"width,omitempty"`
	Height             int    `json:"height,omitempty"`
	CodedWidth         int    `json:"coded_width,omitempty"`
	CodedHeight        int    `json:"coded_height,omitempty"`
	HasBFrames         int    `json:"has_b_frames,omitempty"`
	SampleAspectRatio  string `json:"sample_aspect_ratio,omitempty"`
	DisplayAspectRatio string `json:"display_aspect_ratio,omitempty"`
	PixFmt             string `json:"pix_fmt,omitempty"`
	Level              int    `json:"level,omitempty"`
	ChromaLocation     string `json:"chroma_location,omitempty"`
	Refs               int    `json:"refs,omitempty"`
	IsAvc              string `json:"is_avc,omitempty"`
	NalLengthSize      string `json:"nal_length_size,omitempty"`
	RFrameRate         string `json:"r_frame_rate"`
	AvgFrameRate       string `json:"avg_frame_rate"`
	TimeBase           string `json:"time_base"`
	StartPts           int    `json:"start_pts"`
	StartTime          string `json:"start_time"`
	DurationTs         int    `json:"duration_ts"`
	Duration           string `json:"duration"`
	BitRate            string `json:"bit_rate"`
	BitsPerRawSample   string `json:"bits_per_raw_sample,omitempty"`
	NbFrames           string `json:"nb_frames"`
	SampleFmt          string `json:"sample_fmt,omitempty"`
	SampleRate         string `json:"sample_rate,omitempty"`
	Channels           int    `json:"channels,omitempty"`
	ChannelLayout      string `json:"channel_layout,omitempty"`
	BitsPerSample      int    `json:"bits_per_sample,omitempty"`
	MaxBitRate         string `json:"max_bit_rate,omitempty"`
}

type ffprobeResult struct {
	Streams []stream `json:"streams"`

	Format struct {
		Filename       string `json:"filename"`
		NbStreams      int    `json:"nb_streams"`
		NbPrograms     int    `json:"nb_programs"`
		FormatName     string `json:"format_name"`
		FormatLongName string `json:"format_long_name"`
		StartTime      string `json:"start_time"`
		Duration       string `json:"duration"`
		Size           string `json:"size"`
		BitRate        string `json:"bit_rate"`
		ProbeScore     int    `json:"probe_score"`
	} `json:"format"`
}

func findStream(probe *ffprobeResult, streamType string) (*stream, bool) {
	for _, stream := range probe.Streams {
		if stream.CodecType == streamType {
			return &stream, true
		}
	}

	return nil, false
}

func ffprobeToVideoMeta(probe *ffprobeResult) (*encoder.VideoMeta, error) {
	meta := &encoder.VideoMeta{}

	videoStream, ok := findStream(probe, "video")
	if !ok {
		return nil, encoder.MissingVideoStreamErr
	}

	if audioStream, ok := findStream(probe, "audio"); ok {
		meta.HasAudio = true

		if sampleRate, err := strconv.ParseUint(audioStream.SampleRate, 10, 32); err == nil {
			meta.AudioSampleRate = uint(sampleRate)
		} else {
			return nil, errors.Wrap(err, "failed to parse audio saimple rate")
		}

		meta.AudioCodecName = audioStream.CodecName
		meta.AudioChannelLayout = audioStream.ChannelLayout
	}

	if duration, err := strconv.ParseFloat(probe.Format.Duration, 64); err == nil {
		meta.Duration = duration
	} else {
		return nil, errors.Wrap(err, "Failed to parse duration")
	}

	if size, err := strconv.ParseUint(probe.Format.Size, 10, 64); err == nil {
		meta.Size = size
	} else {
		return nil, errors.Wrap(err, "Failed to parse file size")
	}

	// Do not use the bit rate provided from ffprobe, this will give a average
	meta.BitRate = int(uint64((meta.Size * 8) / uint64(meta.Duration)))

	// Misc
	meta.FrameRate = videoStream.RFrameRate

	// Video
	meta.VideoPixelFormat = videoStream.PixFmt
	meta.VideoCodecName = videoStream.CodecName
	meta.VideoWidth = videoStream.CodedWidth
	meta.VideoHeight = videoStream.CodedHeight

	return meta, nil
}
