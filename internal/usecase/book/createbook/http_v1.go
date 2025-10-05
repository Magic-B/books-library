package createbook

import (
	"encoding/json"
	"net/http"

	"github.com/Magic-B/books-library/pkg/httpserver"
	"github.com/go-chi/render"
)

var usecase *Usecase

func HttpV1(w http.ResponseWriter, r *http.Request) {
	var input Input

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		render.JSON(w, r, httpserver.Error(err.Error()))
		return
	}

	output, err := usecase.Create(r.Context(), input)
	if err != nil {
		render.JSON(w, r, httpserver.Error(err.Error()))
		return
	}

	render.JSON(w, r, output)
}