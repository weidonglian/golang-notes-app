package handlers

import (
	"net/http"

	"github.com/go-chi/render"
)

type DataResponse struct {
	data interface{}
}

func (e DataResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, http.StatusOK)
	return nil
}
