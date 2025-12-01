package main

import (
	"log/slog"
	"os"
	"os/signal"
	"sso/internal/app"
	"sso/internal/config"
	"syscall"
)

const (
	LOCAL = "local"
	DEV   = "dev"
	PROD  = "prod"
)

func main() {
	// init config
	config := config.New()

	// init logger
	logger := setupLogger(config.Env)

	// init app
	application := app.New(logger, config.GRPC.Port, config.StoragePath, config.TokenTTL)

	// launch gRPC server
	go application.GRPCServer.Run()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	<-ch

	application.GRPCServer.Stop()

	logger.Info("application was stopped")
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case LOCAL:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case DEV:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case PROD:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return logger
}
