package ffmpeg

type Preset string

const (
	Ultrafast Preset = "ultrafast"
	Superfast Preset = "superfast"
	Veryfast  Preset = "veryfast"
	Faster    Preset = "faster"
	Fast      Preset = "fast"
	Medium    Preset = "medium"
	Slow      Preset = "slow"
	Slower    Preset = "slower"
	Veryslow  Preset = "veryslow"
)

type VideoQuality struct {
	Scale      string
	Bitrate    string
	MaxBitrate string
	Bufsize    string
}

type VideoQualityKey string

const (
	VideoQualityKey360P  VideoQualityKey = "360p"
	VideoQualityKey720P  VideoQualityKey = "720p"
	VideoQualityKey1080P VideoQualityKey = "1080p"
)

var videoQualities = map[VideoQualityKey]VideoQuality{
	VideoQualityKey360P: {
		Scale:      "-1:360",
		Bitrate:    "365k",
		MaxBitrate: "390k",
		Bufsize:    "640k",
	},
	VideoQualityKey720P: {
		Scale:      "-1:720",
		Bitrate:    "4.5M",
		MaxBitrate: "4.8M",
		Bufsize:    "8M",
	},
	VideoQualityKey1080P: {
		Scale:      "-1:1080",
		Bitrate:    "7.8M",
		MaxBitrate: "8.3M",
		Bufsize:    "14M",
	},
}
