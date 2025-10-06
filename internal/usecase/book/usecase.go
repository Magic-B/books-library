package book

import (
	"context"

	"github.com/Magic-B/books-library/internal/domain"
	"github.com/Magic-B/books-library/pkg/apperr"
	"github.com/Magic-B/books-library/pkg/op"
)

var operation = op.Namespace("usecase.book")

type Usecase struct {
	storage Storage
}

func New(storage Storage) *Usecase {
	return &Usecase{
		storage: storage,
	}
}

func (u *Usecase) Create(ctx context.Context, req CreateReq) (CreateRes, error) {
	op := operation("Create")

	book, err := domain.NewBook(req.Title, req.Description)
	if err != nil {
		return CreateRes{}, apperr.Wrap(op, err, "validation failed")
	}

	if err := u.storage.Create(ctx, &book); err != nil {
		return CreateRes{}, apperr.Wrap(op, err, "failed to create book")
	}

	return CreateRes{Book: book}, nil
}

func (u *Usecase) GetByID(ctx context.Context, id int) (GetByIDRes, error) {
	op := operation("GetByID")

	book, err := u.storage.GetByID(ctx, id)
	if err != nil {
		return GetByIDRes{}, apperr.Wrap(op, err, "failed to get book")
	}

	return GetByIDRes{Book: book}, nil
}
