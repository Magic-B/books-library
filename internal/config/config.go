package config

import (
	"fmt"
	"log"

	"github.com/Magic-B/books-library/internal/adapter/postgres"
	"github.com/Magic-B/books-library/pkg/apperr"
	"github.com/Magic-B/books-library/pkg/httpserver"
	"github.com/Magic-B/books-library/pkg/op"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type App struct {
	Name string `env:"APP_NAME" env-required:"true"`
	Version string `env:"APP_VERSION" env-required:"true"`
	Env string `env:"APP_ENV" env-required:"true"`
}

type Config struct {
	App App
	HttpServer httpserver.Config
	Postgres postgres.Config
}

var operation = op.Namespace("config")

func readConfig[T any](cfg *T, name string) {
	op := operation("readConfig")
	if err := cleanenv.ReadEnv(cfg); err != nil {
		log.Fatal(apperr.Wrap(op, err, fmt.Sprintf("%s config", name)))
	}
}

func MustLoad() *Config {
	op := operation("MustLoad")

	if err := godotenv.Load(); err != nil {
		log.Fatal(apperr.Wrap(op, err))
	}
	
	var app App
	readConfig(&app, "app")
	
	var http httpserver.Config
	readConfig(&http, "httpserver")
	
	var postgres postgres.Config
	readConfig(&postgres, "postgres")

	return &Config{
		App: app,
		HttpServer: http,
		Postgres: postgres,
	}
}