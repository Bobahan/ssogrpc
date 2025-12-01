package main

import (
	"log/slog"
	"os"
	"sso/internal/config"
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

	logger.Info("server has been started")

	// init app
	// launch gRPC server
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
