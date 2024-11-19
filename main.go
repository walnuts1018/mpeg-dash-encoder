package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/walnuts1018/mpeg_dash-encoder/config"
	"github.com/walnuts1018/mpeg_dash-encoder/domain/logger"
	"github.com/walnuts1018/mpeg_dash-encoder/tracer"
	"github.com/walnuts1018/mpeg_dash-encoder/wire"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	logger.CreateAndSetLogger(cfg.LogLevel, cfg.LogType)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, os.Kill)
	defer stop()

	close, err := tracer.NewTracerProvider(ctx)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to create tracer provider: %v", err))
	}
	defer close()

	usecase, err := wire.CreateUsecase(ctx, cfg)
	if err != nil {
		slog.Error("Failed to create usecase", slog.Any("error", err))
		os.Exit(1)
	}

	router, err := wire.CreateRouter(ctx, cfg, usecase)
	if err != nil {
		slog.Error("Failed to create router", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("Server is running", slog.String("port", cfg.ServerPort))

	go func() {
		usecase.Run(ctx)
	}()

	if err := router.Run(fmt.Sprintf(":%s", cfg.ServerPort)); err != nil {
		slog.Error("Failed to run server", slog.Any("error", err))
		os.Exit(1)
	}
}
