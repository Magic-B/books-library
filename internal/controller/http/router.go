package http

import (
	"log/slog"
	"net/http"

	"github.com/Magic-B/books-library/internal/controller/http/middleware"
	bookHandler "github.com/Magic-B/books-library/internal/controller/http/v1/book"
	"github.com/Magic-B/books-library/internal/usecase/book"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

type RouterDeps struct {
	BookUsecase *book.Usecase
	Logger      *slog.Logger
}

func NewRouter(deps RouterDeps) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.Logger(deps.Logger))

	r.Get("/live", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	r.Get("/ready", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			bookH := bookHandler.NewHandler(deps.BookUsecase)
			r.Route("/books", func(r chi.Router) {
				r.Post("/", bookH.Create)
				r.Get("/{id}", bookH.GetByID)
			})
		})
	})

	return r
}
