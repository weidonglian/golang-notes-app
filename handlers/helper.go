package handlers

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

//--
// Error response payloads & renderers
//--

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	AppMessage string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func NewErrorResponse(httpStatusCode int, err error) *ErrResponse {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: httpStatusCode,
		StatusText:     http.StatusText(httpStatusCode),
		AppMessage:     err.Error(),
	}
}

func ErrStatusUnauthorized(err error) *ErrResponse {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnauthorized,
		StatusText:     http.StatusText(http.StatusUnauthorized),
		AppMessage:     err.Error(),
	}
}

func ErrBadRequest(err error) *ErrResponse {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     http.StatusText(http.StatusBadRequest),
		AppMessage:     err.Error(),
	}
}

func ErrUnprocessableEntity(err error) *ErrResponse {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnprocessableEntity,
		StatusText:     http.StatusText(http.StatusUnprocessableEntity),
		AppMessage:     err.Error(),
	}
}

func ErrNotFound() *ErrResponse {
	return &ErrResponse{
		HTTPStatusCode: http.StatusNotFound,
		StatusText:     http.StatusText(http.StatusNotFound),
	}
}
