package createbook

import "github.com/Magic-B/books-library/internal/domain"

type Input struct {
	Book InputParams `json:"book"`
}

type InputParams struct {
	Title string `json:"title"`
	Description string `json:"description"`
}

type Output struct {
	Book domain.Book `json:"book"`
}

