package domain

import (
	"time"

	"github.com/Magic-B/books-library/pkg/apperr"
	"github.com/Magic-B/books-library/pkg/op"
	"github.com/go-playground/validator/v10"
)

type Book struct {
	ID int64 `json:"id"`
	Title string `json:"title" validate:"required"`
	Description string `json:"description"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var validate = validator.New(validator.WithRequiredStructEnabled())

var operation = op.Namespace("domain.book")

func NewBook(title, description string) (Book, error) {
	op := operation("NewBook")

	b := Book{
		Title: title,
		Description: description,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := b.Validate(); err != nil {
		return Book{}, apperr.Wrap(op, err)
	}

	return b, nil
}

func (b Book) Validate() error {
	err := validate.Struct(b)
	if err != nil {
		return apperr.Wrap("validate.Struct", err)
	}

	return nil
}