package domain

import "errors"

var (
	ErrEmptyTitle      = errors.New("title cannot be empty")
	ErrTitleTooLong    = errors.New("title is too long (max 255 characters)")
	ErrDescriptionLong = errors.New("description is too long (max 1000 characters)")
)