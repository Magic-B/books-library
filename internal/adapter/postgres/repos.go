package postgres

import (
	"github.com/Magic-B/books-library/internal/adapter/postgres/repos/bookrepo"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repos struct {
	Book *bookrepo.Repo
}

func NewRepos(p *pgxpool.Pool) *Repos {
	return &Repos{
		Book: bookrepo.New(p),
	}
}