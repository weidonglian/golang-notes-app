package payload

import (
	"github.com/weidonglian/notes-app/handlers/util"
	"github.com/weidonglian/notes-app/model"
	"net/http"
)

// Payload of request user
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

// Payload of Response User
type RespUser struct {
	*model.User
	Password util.OmitField `json:"password,omitempty"`
}

func NewRespUser(user *model.User) RespUser {
	return RespUser{User: user}
}
