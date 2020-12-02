package handlers

import (
	"github.com/weidonglian/notes-app/internal/auth"
	"github.com/weidonglian/notes-app/internal/handlers/payload"
	"github.com/weidonglian/notes-app/internal/lib"
	"github.com/weidonglian/notes-app/internal/store"
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
	if err := lib.ReceiveJson(r, data); err != nil {
		lib.SendErrorBadRequest(w, r, err)
		return
	}

	user := h.s.Users.FindByName(data.Username)
	if user == nil {
		lib.SendErrorUnauthorized(w, r)
		return
	}

	if token, err := h.a.CreateToken(user.ID); err != nil {
		lib.SendErrorUnprocessableEntity(w, r, err)
	} else {
		lib.SendJson(w, r, struct {
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
