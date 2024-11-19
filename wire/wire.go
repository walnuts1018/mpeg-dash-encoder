//go:build wireinject
// +build wireinject

package wire

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/walnuts1018/mpeg-dash-encoder/config"
	"github.com/walnuts1018/mpeg-dash-encoder/infra/ffmpeg"
	"github.com/walnuts1018/mpeg-dash-encoder/infra/jwt"
	"github.com/walnuts1018/mpeg-dash-encoder/infra/minio"
	"github.com/walnuts1018/mpeg-dash-encoder/router"
	"github.com/walnuts1018/mpeg-dash-encoder/router/handler"
	"github.com/walnuts1018/mpeg-dash-encoder/router/middleware"
	"github.com/walnuts1018/mpeg-dash-encoder/usecase"
)

func CreateUsecase(
	ctx context.Context,
	cfg config.Config,
) (*usecase.Usecase, error) {
	wire.Build(
		UsecaseConfigSet,
		jwtSet,
		ffmpegSet,
		minio.NewMinIOClient,
		minioSourceClientSet,
		minioEncodedObjectClientSet,
		usecase.NewUsecase,
	)
	return &usecase.Usecase{}, nil
}

func CreateRouter(
	ctx context.Context,
	cfg config.Config,
	usecase *usecase.Usecase,
) (*gin.Engine, error) {
	wire.Build(
		RouterConfigSet,
		middleware.NewMiddleware,
		handler.NewHandler,
		router.NewRouter,
	)

	return &gin.Engine{}, nil
}

var jwtSet = wire.NewSet(
	jwt.NewManager,
	wire.Bind(new(usecase.TokenIssuer), new(*jwt.Manager)),
)

var minioSourceClientSet = wire.NewSet(
	minio.NewSourceClient,
	wire.Bind(new(usecase.SourceRepository), new(*minio.SourceClient)),
)

var minioEncodedObjectClientSet = wire.NewSet(
	minio.NewEncodedObjectClient,
	wire.Bind(new(usecase.EncodedObjectRepository), new(*minio.EncodedObjectClient)),
)

var ffmpegSet = wire.NewSet(
	ffmpeg.NewFFMPEG,
	wire.Bind(new(usecase.Encoder), new(*ffmpeg.FFMPEG)),
)

var UsecaseConfigSet = wire.FieldsOf(new(config.Config),
	"JWTSigningKey",
	"MinIOSourceUploadBucket",
	"MinIOOutputBucket",
)

var RouterConfigSet = wire.FieldsOf(new(config.Config),
	"AdminToken",
)
