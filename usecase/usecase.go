package usecase

import (
	"context"
	"fmt"
	"io"
	"iter"
	"os"
	"time"

	"github.com/walnuts1018/mpeg-dash-encoder/config"
	"github.com/walnuts1018/mpeg-dash-encoder/domain/entity"
)

type Usecase struct {
	tokenIssuer TokenIssuer
	encoder     Encoder
	sourceRepo  SourceRepository
	encodedRepo EncodedObjectRepository

	encodeQueue   chan encodeRequest
	encodeTimeout time.Duration
	hostname      string
}

type TokenIssuer interface {
	CreateUserToken(mediaIDs []string) (string, error)
	GetMediaIDsFromToken(token string) ([]string, error)
}

type SourceRepository interface {
	ListUploadedFiles(ctx context.Context) iter.Seq2[entity.SourceFile, error]
	SetObjectTags(ctx context.Context, id string, tags map[string]string) error
	RemoveObjectTags(ctx context.Context, id string) error
	GetSourceContent(ctx context.Context, id string) (io.ReadSeekCloser, error)
	DeleteSourceContent(ctx context.Context, id string) error
}

type EncodedObjectRepository interface {
	Upload(ctx context.Context, mediaID string, localDir string) error
	GetObject(ctx context.Context, mediaID string, fileName string) (io.ReadSeekCloser, error)
}

type Encoder interface {
	Encode(id string, path string, audioOnly bool) (string, error)
	GetOutDirPrefix() string
}

func NewUsecase(
	cfg config.Config,
	tokenIssuer TokenIssuer,
	encoder Encoder,
	sourceRepo SourceRepository,
	encodedRepo EncodedObjectRepository,
) (*Usecase, error) {

	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname: %w", err)
	}

	return &Usecase{
		tokenIssuer:   tokenIssuer,
		encoder:       encoder,
		sourceRepo:    sourceRepo,
		encodedRepo:   encodedRepo,
		encodeQueue:   make(chan encodeRequest, 0),
		encodeTimeout: cfg.EncodeTimeout,
		hostname:      hostname,
	}, nil
}
