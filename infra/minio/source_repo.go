package minio

import (
	"context"
	"errors"
	"fmt"
	"io"
	"iter"

	"github.com/minio/minio-go/v7"
	miniotags "github.com/minio/minio-go/v7/pkg/tags"
	"github.com/walnuts1018/mpeg-dash-encoder/config"
	"github.com/walnuts1018/mpeg-dash-encoder/domain/entity"
)

type SourceClient struct {
	bucketName string
	client     *minio.Client
}

func NewSourceClient(bucketName config.SourceClientBucketName, client *minio.Client) *SourceClient {
	return &SourceClient{
		bucketName: string(bucketName),
		client:     client,
	}
}

func (m *SourceClient) ListUploadedFiles(ctx context.Context) iter.Seq2[entity.SourceFile, error] {
	infos := m.client.ListObjects(ctx, m.bucketName, minio.ListObjectsOptions{
		WithMetadata: true,
	})
	return func(yield func(entity.SourceFile, error) bool) {
		for info := range infos {
			if info.Err != nil {
				if !yield(entity.SourceFile{}, fmt.Errorf("failed to list objects: %w", info.Err)) {
					return
				}
			}
			if info.Key == "" {
				if !yield(entity.SourceFile{}, errors.New("source object not found")) {
					return
				}
			}
			if !yield(entity.SourceFile{ID: info.Key, Tags: info.UserTags}, nil) {
				return
			}
		}
	}
}

func (m *SourceClient) SetObjectTags(ctx context.Context, id string, tags map[string]string) error {
	newtag, err := miniotags.MapToObjectTags(tags)
	if err != nil {
		return fmt.Errorf("failed to create tags: %w", err)
	}

	if err := m.client.PutObjectTagging(ctx, m.bucketName, id, newtag, minio.PutObjectTaggingOptions{}); err != nil {
		return fmt.Errorf("failed to set tags: %w", err)
	}
	return nil
}

func (m *SourceClient) RemoveObjectTags(ctx context.Context, id string) error {
	if err := m.client.RemoveObjectTagging(ctx, m.bucketName, id, minio.RemoveObjectTaggingOptions{}); err != nil {
		return fmt.Errorf("failed to remove tags: %w", err)
	}
	return nil
}

func (m *SourceClient) GetSourceContent(ctx context.Context, id string) (io.ReadSeekCloser, error) {
	obj, err := m.client.GetObject(ctx, m.bucketName, id, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	return obj, nil
}

func (m *SourceClient) DeleteSourceContent(ctx context.Context, id string) error {
	if err := m.client.RemoveObject(ctx, m.bucketName, id, minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}
	return nil
}
