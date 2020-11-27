package handlers

import (
	"github.com/weidonglian/notes-app/auth"
	"github.com/weidonglian/notes-app/handlers/payload"
	"github.com/weidonglian/notes-app/handlers/util"
	"github.com/weidonglian/notes-app/store"
	"net/http"
)

// SessionHandler keeps the dependency for handle session.
type SessionHandler struct {
	s *store.Store
	a *auth.Auth
}

func NewSessionHandler(s *store.Store, a *auth.Auth) SessionHandler {
	return SessionHandler{
		s: s,
		a: a,
	}
}

// NewSession POST /session
func (h SessionHandler) NewSession(w http.ResponseWriter, r *http.Request) {
	data := &payload.ReqSession{}
	if err := util.ReceiveJson(r, data); err != nil {
		util.SendErrorBadRequest(w, r, err)
		return
	}

	user := h.s.Users.FindByName(data.Username)
	if user == nil {
		util.SendErrorUnauthorized(w, r)
		return
	}

	if token, err := h.a.CreateToken(user.ID); err != nil {
		util.SendErrorUnprocessableEntity(w, r, err)
	} else {
		util.SendJson(w, r, struct {
			Token string `json:"token"`
		}{
			Token: token,
		})
	}
}

// DeleteSession DELETE /session
func (h SessionHandler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("DeleteSession"))
}
