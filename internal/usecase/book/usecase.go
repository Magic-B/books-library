package book

import (
	"context"

	"github.com/Magic-B/books-library/internal/domain"
	"github.com/Magic-B/books-library/pkg/apperr"
)

type Usecase struct {
	storage Storage
}

func New(storage Storage) *Usecase {
	return &Usecase{
		storage: storage,
	}
}

func (u *Usecase) Create(ctx context.Context, req CreateReq) (CreateRes, error) {
	book, err := domain.NewBook(req.Book.Title, req.Book.Description)
	if err != nil {
		return CreateRes{}, apperr.WithDesc("validation failed", err)
	}

	if err := u.storage.Create(ctx, &book); err != nil {
		return CreateRes{}, apperr.WithDesc("failed to create book", err)
	}

	return CreateRes{Book: book}, nil
}

func (u *Usecase) GetByID(ctx context.Context, id int) (GetByIDRes, error) {

	book, err := u.storage.GetByID(ctx, id)
	if err != nil {
		return GetByIDRes{}, apperr.WithDesc("failed to get book", err)
	}

	return GetByIDRes{Book: book}, nil
}
