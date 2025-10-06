package book

import (
	"context"

	"github.com/Magic-B/books-library/internal/domain"
)

type Storage interface {
	Create(ctx context.Context, b *domain.Book) error
	GetByID(ctx context.Context, id int) (domain.Book, error)
}
