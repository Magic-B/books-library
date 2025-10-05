package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Magic-B/books-library/internal/adapter/postgres"
	"github.com/Magic-B/books-library/internal/config"
	"github.com/Magic-B/books-library/internal/controller/http"
	"github.com/Magic-B/books-library/internal/usecase/book/createbook"
	"github.com/Magic-B/books-library/pkg/httpserver"
	"github.com/Magic-B/books-library/pkg/logger/slg"
)

const (
	envLocal = "local"
	envDev = "dev"
	envProd = "prod"
)

func main() {
	cfg := config.MustLoad()

	logger := loggerSetup(cfg.App.Env)

	ctx := context.Background()

	pgPool, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		logger.Error("faild to init postgres", slg.Error(err))
	}
	defer pgPool.Close()
	logger.Info("postgres has been inited")

	initUscases(pgPool.Repos)

	//http
	router := http.Router()
	server := httpserver.New(router, cfg.HttpServer)

	logger.Info("run http server")
	server.Run()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	<-sig

	pgPool.Close()
	server.Close(ctx)
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

func initUscases(pg *postgres.Repos) {
	createbook.New(pg.Book)
}
