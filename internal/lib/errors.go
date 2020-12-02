package lib

import (
	"errors"
	"github.com/go-chi/render"
	"net/http"
)

// List of the error messages.
var (
	ErrorMissingRequiredFields = errors.New("missing required fields")
	ErrorMissingRequiredParams = errors.New("missing required params")
	ErrorUnauthorized          = errors.New("missing required credentials")
	ErrorUnprocessableEntity   = errors.New("UnprocessableEntity")
)

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
