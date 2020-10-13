package handlers

import (
	"github.com/go-chi/render"
	"github.com/weidonglian/golang-notes-app/errors"
	"net/http"
)

func SendJson(w http.ResponseWriter, r *http.Request, v interface{}) {
	render.JSON(w, r, v)
}

func SendError(w http.ResponseWriter, r *http.Request, httpStatusCode int, err error) {
	render.Render(w, r, errors.NewErrorResponse(httpStatusCode, err))
}
