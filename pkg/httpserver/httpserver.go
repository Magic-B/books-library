package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Config struct {
	Host            string        `env:"HTTP_ADDRESS" env-required:"true"`
	Port            string        `env:"HTTP_PORT" env-required:"true"`
	Timeout         time.Duration `env:"HTTP_TIMEOUT" env-required:"true"`
	IdleTimeout     time.Duration `env:"HTTP_IDLE_TIMEOUT" env-required:"true"`
	ShutdownTimeout time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT" env-required:"true"`
}

type Server struct {
	server *http.Server
	cfg    Config
}

func New(handler http.Handler, cfg Config) *Server {
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
		Handler:      handler,
	}

	return &Server{
		server: srv,
		cfg:    cfg,
	}
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

func (s *Server) Close(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, s.cfg.ShutdownTimeout)
	defer cancel()

	return s.server.Shutdown(shutdownCtx)
}
