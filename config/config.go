package config

import (
	"errors"
	"log/slog"
	"reflect"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
)

var ErrInvalidSessionSecretLength = errors.New("session secret must be 16, 24, or 32 bytes")

type Config struct {
	ServerPort string `env:"SERVER_PORT" envDefault:"8080"`

	// ------------------------ Log ------------------------
	LogLevel slog.Level `env:"LOG_LEVEL"`
	LogType  LogType    `env:"LOG_TYPE" envDefault:"json"`
	LogDir   string     `env:"LOG_DIR" envDefault:"/var/log/mucaron"`

	// ------------------------ Application ------------------------
	AdminToken    AdminToken    `env:"ADMIN_TOKEN,required"`
	MaxUploadSize uint64        `env:"MAX_UPLOAD_SIZE" envDefault:"1073741824"` //1GB
	EncodeTimeout time.Duration `env:"ENCODE_TIMEOUT" envDefault:"1h"`
	FFMPEGHwAccel string        `env:"FFMPEG_HW_ACCEL" envDefault:"qsv"`
	JWTSigningKey JWTSigningKey `env:"JWT_SIGN_SECRET,required"`

	// ------------------------ MinIO ------------------------
	MinIOEndpoint       string `env:"MINIO_ENDPOINT" envDefault:"localhost:9000"`
	MinIOAccessKey      string `env:"MINIO_ACCESS_KEY,required"`
	MinIOSecretKey      string `env:"MINIO_SECRET_KEY,required"`
	MinIOUseSSL         bool   `env:"MINIO_USE_SSL" envDefault:"false"`
	MinIOPublicEndpoint string `env:"MINIO_PUBLIC_ENDPOINT" envDefault:""` // localhost:9000

	MinIOSourceUploadBucket SourceClientBucketName  `env:"MINIO_SOURCE_UPLOAD_BUCKET" envDefault:"mpeg-dash-encoder-source-upload"`
	MinIOOutputBucket       EncodedObjectBucketName `env:"MINIO_OUTPUT_BUCKET" envDefault:"mpeg-dash-encoder-output"`
}

func Load() (Config, error) {
	var cfg Config
	if err := env.ParseWithOptions(&cfg, env.Options{
		FuncMap: map[reflect.Type]env.ParserFunc{
			reflect.TypeOf(slog.Level(0)):    returnAny(ParseLogLevel),
			reflect.TypeOf(time.Duration(0)): returnAny(time.ParseDuration),
			reflect.TypeOf(LogType("")):      returnAny(ParseLogType),
		},
	}); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func returnAny[T any](f func(v string) (t T, err error)) func(v string) (any, error) {
	return func(v string) (any, error) {
		t, err := f(v)
		return any(t), err
	}
}

func ParseLogLevel(v string) (slog.Level, error) {
	switch strings.ToLower(v) {
	case "":
		return slog.LevelInfo, nil
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		slog.Warn("Invalid log level, use default level: info")
		return slog.LevelInfo, nil
	}
}

type LogType string

const (
	LogTypeJSON LogType = "json"
	LogTypeText LogType = "text"
)

func ParseLogType(v string) (LogType, error) {
	switch strings.ToLower(v) {
	case "json":
		return LogTypeJSON, nil
	case "text":
		return LogTypeText, nil
	default:
		slog.Warn("Invalid log type, use default type: json")
		return LogTypeJSON, nil
	}
}
