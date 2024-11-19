package minio

import (
	"context"
	"fmt"
	"io/fs"
	"path"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/walnuts1018/mpeg-dash-encoder/config"
)

type EncodedObjectClient struct {
	bucketName string
	client     *minio.Client
}

func NewEncodedObjectClient(bucketName config.EncodedObjectBucketName, client *minio.Client) *EncodedObjectClient {
	return &EncodedObjectClient{
		bucketName: string(bucketName),
		client:     client,
	}
}

func (m *EncodedObjectClient) Upload(ctx context.Context, mediaID string, localDir string) error {
	if err := filepath.WalkDir(localDir, func(localFilePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}
		localRelativeFilePath, err := filepath.Rel(localDir, localFilePath)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		objectPath := path.Join(mediaID, filepath.ToSlash(localRelativeFilePath))
		if _, err := m.client.FPutObject(ctx, m.bucketName, objectPath, localFilePath, minio.PutObjectOptions{}); err != nil {
			return fmt.Errorf("failed to put object: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to upload directory: %w", err)
	}
	return nil
}
