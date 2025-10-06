package book

import "github.com/Magic-B/books-library/internal/domain"

// Create operation
type CreateReq struct {
	Book BookInput `json:"book"`
}

type BookInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateRes struct {
	Book domain.Book `json:"book"`
}

// GetByID operation
type GetByIDReq struct {
	ID int `json:"id"`
}

type GetByIDRes struct {
	Book domain.Book `json:"book"`
}
