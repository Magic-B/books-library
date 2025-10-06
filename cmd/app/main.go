package main

import (
	"context"
	"fmt"
	"log/slog"
	nethttp "net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Magic-B/books-library/internal/adapter/postgres"
	"github.com/Magic-B/books-library/internal/app"
	"github.com/Magic-B/books-library/internal/config"
	"github.com/Magic-B/books-library/internal/controller/http"
	"github.com/Magic-B/books-library/pkg/httpserver"
	"github.com/Magic-B/books-library/pkg/logger/slg"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func gracefulShutdown(
	ctx context.Context,
	serverError <-chan error,
	server *httpserver.Server,
	pgPool *postgres.Pool,
	logger *slog.Logger,
) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverError:
		if err != nil && err != nethttp.ErrServerClosed {
			logger.Error("server failed to start", slg.Error(err))
			os.Exit(1)
		}
	case s := <-sig:
		logger.Info(fmt.Sprintf("received signal: %v, shutting down", s))
	}

	logger.Info("shutting down server...")
	if err := server.Close(ctx); err != nil {
		logger.Error("failed to shutdown server gracefully", slg.Error(err))
	}

	pgPool.Close()
	logger.Info("server stopped")
}

func main() {
	cfg := config.MustLoad()

	logger := loggerSetup(cfg.App.Env)

	ctx := context.Background()

	pgPool, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		logger.Error("failed to init postgres", slg.Error(err))
		os.Exit(1)
	}
	defer pgPool.Close()
	logger.Info("postgres initialized")

	app := app.New(pgPool.Repos)

	router := http.NewRouter(http.RouterDeps{
		Usecases: app.Usecases,
		Logger:   logger,
	})
	server := httpserver.New(router, cfg.HttpServer)

	serverError := make(chan error, 1)
	go func() {
		logger.Info(fmt.Sprintf("starting http server on %s:%s", cfg.HttpServer.Host, cfg.HttpServer.Port))
		serverError <- server.Run()
	}()

	gracefulShutdown(ctx, serverError, server, pgPool, logger)
}

func loggerSetup(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
