package wire

import (
	"github.com/walnuts1018/mpeg_dash-encoder/infra/ffmpeg"
	"github.com/walnuts1018/mpeg_dash-encoder/infra/jwt"
	"github.com/walnuts1018/mpeg_dash-encoder/infra/minio"
	"github.com/walnuts1018/mpeg_dash-encoder/usecase"
)

var _ usecase.TokenIssuer = &jwt.Manager{}
var _ usecase.Encoder = &ffmpeg.FFMPEG{}
var _ usecase.SourceRepository = &minio.SourceClient{}
var _ usecase.EncodedObjectRepository = &minio.EncodedObjectClient{}
