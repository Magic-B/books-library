package app

import (
	"github.com/Magic-B/books-library/internal/adapter/postgres"
	"github.com/Magic-B/books-library/internal/usecase/book"
)

type Usecases struct {
	Book *book.Usecase
}

type Application struct {
	Usecases *Usecases
}

func New(repos *postgres.Repos) *Application {
	return &Application{
		Usecases: &Usecases{
			Book: book.New(repos.Book),
		},
	}
}