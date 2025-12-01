package app

import (
	"log/slog"
	"sso/internal/app/grpc"
	"time"
)

type App struct {
	GRPCServer *grpc.App
}

func New(logger *slog.Logger, port int, storagePath string, tokenTTL time.Duration) *App {
	// init storage
	// init auth service

	gRPCApp := grpc.New(logger, port)

	return &App{
		GRPCServer: gRPCApp,
	}
}
