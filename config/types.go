package config

type JWTSigningKey string

type SourceClientBucketName string

type EncodedObjectBucketName string

type AdminToken string

type FFmpegHWAccel string

const (
	FFmpegHWAccelNone FFmpegHWAccel = ""
	FFmpegHWAccelQSV  FFmpegHWAccel = "qsv"
)

type FFmpegPreset string

const (
	Ultrafast FFmpegPreset = "ultrafast"
	Superfast FFmpegPreset = "superfast"
	Veryfast  FFmpegPreset = "veryfast"
	Faster    FFmpegPreset = "faster"
	Fast      FFmpegPreset = "fast"
	Medium    FFmpegPreset = "medium"
	Slow      FFmpegPreset = "slow"
	Slower    FFmpegPreset = "slower"
	Veryslow  FFmpegPreset = "veryslow"
)

type FFmpegConfig struct {
	LogDir     string        `env:"LOG_DIR" envDefault:"/var/log/mpeg-dash-encoder/ffmpeg"`
	FPS        int           `env:"FPS" envDefault:"30"`
	Preset     FFmpegPreset  `env:"PRESET" envDefault:"medium"`
	HWAccel    FFmpegHWAccel `env:"FFMPEG_HW_ACCEL"`
	AudioCodec string        `env:"AUDIO_CODEC" envDefault:"aac"`
}
