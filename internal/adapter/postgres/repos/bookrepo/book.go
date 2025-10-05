package bookrepo

import (
	"context"

	"github.com/Magic-B/books-library/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

func (r *Repo) Create(ctx context.Context, b *domain.Book) error {
	const q = `INSERT INTO books (title, description) VALUES ($1, $2) RETURNING id`

	err := r.pool.QueryRow(ctx, q, b.Title, b.Description).Scan(&b.ID)
	if err != nil {
		return err
	}

	return nil
}


func (r *Repo) Get(ctx context.Context, id int) (domain.Book, error) {
	const q = `SELECT id, title, description FROM books WHERE id=$1`

	var b domain.Book
	err := r.pool.QueryRow(ctx, q, id).Scan(&b.ID, &b.Title, &b.Description)
	if err != nil {
		return domain.Book{}, err
	}

	return b, nil
}


