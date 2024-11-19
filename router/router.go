package router

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"github.com/walnuts1018/mpeg_dash-encoder/config"
	"github.com/walnuts1018/mpeg_dash-encoder/consts"
	"github.com/walnuts1018/mpeg_dash-encoder/router/handler"
	"github.com/walnuts1018/mpeg_dash-encoder/router/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func NewRouter(config config.Config, handler handler.Handler, m *middleware.Middleware) (*gin.Engine, error) {
	if config.LogLevel != slog.LevelDebug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(sloggin.NewWithConfig(slog.Default(), sloggin.Config{
		DefaultLevel:     config.LogLevel,
		ClientErrorLevel: slog.LevelWarn,
		ServerErrorLevel: slog.LevelError,

		WithUserAgent:      false,
		WithRequestID:      true,
		WithRequestBody:    false,
		WithRequestHeader:  false,
		WithResponseBody:   false,
		WithResponseHeader: false,
		WithSpanID:         true,
		WithTraceID:        true,

		Filters: []sloggin.Filter{
			sloggin.IgnorePath("/healthz"),
		},
	}))
	r.Use(otelgin.Middleware(consts.ApplicationName))

	r.GET("/healthz", handler.Health)
	v1 := r.Group("/v1")

	admin := v1.Group("/admin")
	admin.Use(m.AdminAuth())
	{
		admin.POST("/create_user_token", handler.CreateUserToken)
	}

	// user := v1.Group("/user")
	// {
	// 	// user.GET("/:media_id/dash.mpd", handler.GetDashManifest)
	// 	// user.GET("/:media_id/:segment", handler.GetSegmentFile)
	// }

	return r, nil
}
