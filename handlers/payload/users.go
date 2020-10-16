package payload

import (
	"github.com/weidonglian/golang-notes-app/handlers/util"
	"net/http"
)

type ReqUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (req ReqUser) Bind(r *http.Request) error {
	if req.Username == "" || req.Password == "" {
		return util.ErrorMissingRequiredFields
	}
	return nil
}
