package minio

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/walnuts1018/mpeg_dash-encoder/config"
)

const (
	accessKey              = "mockaccesskey"
	secretKey              = "mocksecretkey"
	sourceClientBucketName = "mpeg-dash-encoder-source-upload"
	outputBucketName       = "mpeg-dash-encoder-output"
)

var (
	hostAndPort string
	minioClient *minio.Client
	testData    = []byte("testdata")
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		slog.Error("failed to create pool", slog.Any("error", err))
		os.Exit(1)
	}

	if err := pool.Client.Ping(); err != nil {
		slog.Error("failed to connect to Docker", slog.Any("error", err))
		os.Exit(1)
	}

	resource, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: "minio/minio",
			Tag:        "latest",
			Env: []string{
				fmt.Sprintf("MINIO_ACCESS_KEY=%s", accessKey),
				fmt.Sprintf("MINIO_SECRET_KEY=%s", secretKey),
			},
			Cmd: []string{"server", "/export", "--console-address", ":9001"},
		},
		func(config *docker.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
		},
	)
	if err != nil {
		slog.Error("failed to create pool", slog.Any("error", err))
		os.Exit(1)
	}

	hostAndPort = resource.GetHostPort("9000/tcp")

	if err := pool.Retry(func() error {
		var err error
		minioClient, err = NewMinIOClient(config.Config{
			MinIOEndpoint:  hostAndPort,
			MinIOAccessKey: accessKey,
			MinIOSecretKey: secretKey,
			MinIOUseSSL:    false,
		})
		if err != nil {
			slog.Error("failed to create minio client", slog.Any("error", err))
			os.Exit(1)
		}

		cancel, err := minioClient.HealthCheck(10 * time.Second)
		if err != nil {
			return fmt.Errorf("failed to call healthcheck: %v", err)
		}
		defer cancel()

		online := minioClient.IsOnline()
		if online {
			return nil
		} else {
			return fmt.Errorf("minio is offline")
		}

	}); err != nil {
		slog.Error("failed to retry", slog.Any("error", err))
		os.Exit(1)
	}

	ctx := context.Background()
	for _, bucketName := range []string{sourceClientBucketName, outputBucketName} {
		bucketExist, err := minioClient.BucketExists(ctx, bucketName)
		if err != nil {
			slog.Error("failed to check bucket", slog.Any("error", err))
			os.Exit(1)
		}
		if !bucketExist {
			if err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
				slog.Error("failed to create bucket", slog.Any("error", err))
				os.Exit(1)
			}
		}
	}

	defer func() {
		if err := pool.Purge(resource); err != nil {
			slog.Error("failed to purge resources", slog.Any("error", err))
			os.Exit(1)
		}

	}()

	m.Run()
}
