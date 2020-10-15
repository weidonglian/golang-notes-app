package util

import (
	"github.com/go-chi/render"
	"net/http"
)

// ReceiveJson will extract the required struct that implements the binder interface
// returns a error if failed
func ReceiveJson(r *http.Request, v render.Binder) error {
	return render.Bind(r, v)
}

// SendJson sends a given json struct as response
func SendJson(w http.ResponseWriter, r *http.Request, v interface{}) {
	render.JSON(w, r, v)
}

// SendStatus only sends back the status without any message.
func SendStatus(w http.ResponseWriter, r *http.Request, httpStatusCode int) {
	w.WriteHeader(httpStatusCode)
}

// SendError sends the http code with given error message
func SendError(w http.ResponseWriter, r *http.Request, httpStatusCode int, err error) {
	render.Render(w, r, NewErrorResponse(httpStatusCode, err))
}
