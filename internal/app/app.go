package app

import (
	"log"
	"log/slog"
	"sso/internal/app/grpc"
	"sso/internal/services/auth"
	"sso/internal/storage/sqlite"
	"time"
)

type App struct {
	GRPCServer *grpc.App
}

func New(logger *slog.Logger, port int, storagePath string, tokenTTL time.Duration) *App {
	// init storage
	storage, err := sqlite.New(storagePath)

	if err != nil {
		log.Fatal("failed to create storage")
	}

	// init auth service
	authService := auth.New(logger, storage, storage, storage, tokenTTL) // logger, IUserSaver, IUserProvider, IAppProvider, tokenTTL

	gRPCApp := grpc.New(logger, port, authService)

	return &App{
		GRPCServer: gRPCApp,
	}
}
