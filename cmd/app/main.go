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
	"github.com/Magic-B/books-library/internal/config"
	"github.com/Magic-B/books-library/internal/controller/http"
	"github.com/Magic-B/books-library/internal/usecase/book"
	"github.com/Magic-B/books-library/pkg/httpserver"
	"github.com/Magic-B/books-library/pkg/logger/slg"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	logger := loggerSetup(cfg.App.Env)

	ctx := context.Background()

	// Init database
	pgPool, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		logger.Error("failed to init postgres", slg.Error(err))
		os.Exit(1)
	}
	defer pgPool.Close()
	logger.Info("postgres initialized")

	// Init usecases
	bookUsecase := book.New(pgPool.Repos.Book)

	// Init HTTP server
	router := http.NewRouter(http.RouterDeps{
		BookUsecase: bookUsecase,
		Logger:      logger,
	})
	server := httpserver.New(router, cfg.HttpServer)

	// Start server in goroutine
	serverErrors := make(chan error, 1)
	go func() {
		logger.Info(fmt.Sprintf("starting http server on %s:%s", cfg.HttpServer.Host, cfg.HttpServer.Port))
		serverErrors <- server.Run()
	}()

	// Wait for interrupt signal or server error
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if err != nil && err != nethttp.ErrServerClosed {
			logger.Error("server failed to start", slg.Error(err))
			os.Exit(1)
		}
	case s := <-sig:
		logger.Info(fmt.Sprintf("received signal: %v, shutting down", s))
	}

	// Graceful shutdown
	logger.Info("shutting down server...")
	if err := server.Close(ctx); err != nil {
		logger.Error("failed to shutdown server gracefully", slg.Error(err))
	}

	pgPool.Close()
	logger.Info("server stopped")
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
