package domain

import (
	"strings"
	"time"
)

type Book struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewBook(title, description string) (Book, error) {
	b := Book{
		Title:       title,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := b.Validate(); err != nil {
		return Book{}, err
	}

	return b, nil
}

func (b Book) Validate() error {
	title := strings.TrimSpace(b.Title)

	if title == "" {
		return ErrEmptyTitle
	}

	if len(title) > 255 {
		return ErrTitleTooLong
	}

	if len(b.Description) > 1000 {
		return ErrDescriptionLong
	}

	return nil
}
