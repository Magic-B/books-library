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
	const q = `INSERT INTO books (title, description) VALUES ($1, $2) RETURNING id, created_at, updated_at`

	err := r.pool.QueryRow(ctx, q, b.Title, b.Description).Scan(&b.ID, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) GetByID(ctx context.Context, id int) (domain.Book, error) {
	const q = `SELECT id, title, description, created_at, updated_at FROM books WHERE id=$1`

	var b domain.Book
	err := r.pool.QueryRow(ctx, q, id).Scan(&b.ID, &b.Title, &b.Description, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return domain.Book{}, err
	}

	return b, nil
}
