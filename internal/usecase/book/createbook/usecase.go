package createbook

import (
	"context"
	"github.com/Magic-B/books-library/internal/domain"
)

type Storage interface {
	Create(ctx context.Context, b *domain.Book) error
}

type Usecase struct {
	storage Storage
}

func New(storage Storage) *Usecase {
	uc := &Usecase{
		storage: storage,
	}

	usecase = uc

	return uc
}

func (u *Usecase) Create(ctx context.Context, input Input) (Output, error) {
	book, err := domain.NewBook(input.Book.Title, input.Book.Description)
	if err != nil {
		return Output{}, err
	}

	err = u.storage.Create(ctx, &book)
	if err != nil {
		return Output{}, err
	}

	return Output{Book: book}, nil
}
