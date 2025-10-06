package book

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Magic-B/books-library/internal/usecase/book"
	"github.com/Magic-B/books-library/pkg/httpserver"
	"github.com/Magic-B/books-library/pkg/logger/slg"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Handler struct {
	bookUsecase *book.Usecase
	logger      *slog.Logger
}

func NewHandler(bookUsecase *book.Usecase, logger *slog.Logger) *Handler {
	return &Handler{
		bookUsecase: bookUsecase,
		logger:      logger,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req book.CreateReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, httpserver.Error("invalid request body"))
		return
	}

	res, err := h.bookUsecase.Create(r.Context(), req)
	if err != nil {
		h.logger.Error("failed to create book", slg.Error(err))
		status, msg := HandleError(err)
		w.WriteHeader(status)
		render.JSON(w, r, httpserver.Error(msg))
		return
	}

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, res)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, httpserver.Error("invalid book id"))
		return
	}

	res, err := h.bookUsecase.GetByID(r.Context(), id)
	if err != nil {
		h.logger.Error("failed to get book", slg.Error(err))
		status, msg := HandleError(err)
		w.WriteHeader(status)
		render.JSON(w, r, httpserver.Error(msg))
		return
	}

	render.JSON(w, r, res)
}
