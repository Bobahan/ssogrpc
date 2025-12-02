package grpc

import (
	"fmt"
	"log/slog"
	"net"
	"sso/internal/grpc/auth"
	authService "sso/internal/services/auth"
	"sso/internal/storage/sqlite"

	"google.golang.org/grpc"
)

type App struct {
	logger     *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(logger *slog.Logger, port int) *App {
	gRPCServer := grpc.NewServer()

	storage, err := sqlite.New("asd")

	aa := authService.New(logger)

	auth.PlugHandlers(gRPCServer, aa) // к gRPC сервису подключаем handlers

	return &App{
		logger:     logger,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) Run() error {
	const op = "app.Run"

	logger := a.logger.With(slog.String("op", op))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	logger.Info("gRPC server is running", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "app.Stop"

	logger := a.logger.With(slog.String("op", op), slog.Int("port", a.port))

	logger.Info("stopping gRPC server")

	a.gRPCServer.GracefulStop()
}
