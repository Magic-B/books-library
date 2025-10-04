package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Magic-B/books-library/pkg/apperr"
	"github.com/Magic-B/books-library/pkg/op"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	User            string        `env:"POSTGRES_USER" env-required:"true"`
	Password        string        `env:"POSTGRES_PASSWORD" env-required:"true"`
	DBName          string        `env:"POSTGRES_DB_NAME" env-required:"true"`
	Port            string        `env:"POSTGRES_PORT" env-required:"true"`
	Host            string        `env:"POSTGRES_HOST" env-required:"true"`
	MaxConns        int32         `env:"POSTGRES_MAX_CONNS" env-default:"10"`
	MinConns        int32         `env:"POSTGRES_MIN_CONNS" env-default:"0"`
	MaxConnLifetime time.Duration `env:"POSTGRES_MAX_CONN_LIFETIME" env-default:"60m"`
	MaxConnIdleTime time.Duration `env:"POSTGRES_MAX_CONN_IDLE_TIME" env-default:"10m"`
	HealthTimeout   time.Duration `env:"POSTGRES_HEALTH_TIMEOUT" env-default:"3s"`
}

type Pool struct {
	pool *pgxpool.Pool
}

var operation = op.Namespace("adapter.postgres")

func New(ctx context.Context, cfg Config) (*Pool, error) {
	op := operation("New")
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	pgxCfg, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		return nil, apperr.Wrap(op, err)
	}

	pgxCfg.MaxConns = cfg.MaxConns
	pgxCfg.MinConns = cfg.MinConns
	pgxCfg.MaxConnLifetime = cfg.MaxConnLifetime
	pgxCfg.MaxConnIdleTime = cfg.MaxConnIdleTime

	p, err := pgxpool.NewWithConfig(ctx, pgxCfg)

	if err != nil {
		return nil, apperr.Wrap(op, err)
	}
	
	healthCtx, cancel := context.WithTimeout(ctx, cfg.HealthTimeout)
	defer cancel()
	if err := p.Ping(healthCtx); err != nil {
		p.Close()
		return nil, apperr.Wrap(op, err)
	}

	return &Pool{pool: p}, nil
}

func (p *Pool) Close() {
	if p.pool != nil {
		p.pool.Close()
	}
}
