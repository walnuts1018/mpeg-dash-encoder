package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/walnuts1018/mpeg-dash-encoder/config"
)

func CreateAndSetLogger(logLevel slog.Level, logType config.LogType) {
	var hander slog.Handler
	switch logType {
	case config.LogTypeText:
		hander = tint.NewHandler(os.Stdout, &tint.Options{
			Level:     logLevel,
			AddSource: logLevel == slog.LevelDebug,
		})
	case config.LogTypeJSON:
		hander = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     logLevel,
			AddSource: logLevel == slog.LevelDebug,
		})
	}

	logger := slog.New(newTraceHandler(hander))
	slog.SetDefault(logger)
}
