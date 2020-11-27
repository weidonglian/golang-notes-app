package payload

import (
	"github.com/weidonglian/notes-app/handlers/util"
	"net/http"
)

// login and auth a new session
type ReqSession struct {
	Username string `json:"username" `
	Password string `json:"password"`
}

func (req *ReqSession) Bind(r *http.Request) error {
	if req.Username == "" || req.Password == "" {
		return util.ErrorMissingRequiredFields
	}
	return nil
}
