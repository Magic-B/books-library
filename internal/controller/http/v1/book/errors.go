package book

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Magic-B/books-library/internal/domain"
	"github.com/jackc/pgx/v5"
)

func HandleError(err error) (int, string) {
	if errors.Is(err, domain.ErrEmptyTitle) ||
		errors.Is(err, domain.ErrTitleTooLong) ||
		errors.Is(err, domain.ErrDescriptionLong) {
		return http.StatusBadRequest, err.Error()
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return http.StatusNotFound, "resource not found"
	}

	if strings.Contains(err.Error(), "validation failed") {
		return http.StatusBadRequest, err.Error()
	}

	return http.StatusInternalServerError, "internal server error"
}
