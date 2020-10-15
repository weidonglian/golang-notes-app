package handlers

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/weidonglian/golang-notes-app/auth"
	"github.com/weidonglian/golang-notes-app/store"
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

func (h SessionHandler) Routes() chi.Router {
	// Routes for /session
	r := chi.NewRouter()

	return r
}

// login and auth a new session
type reqSession struct {
	Username string `db:"user_name" json:"username" `
	Password string `db:"user_password" json:"password"`
}

func (req *reqSession) Bind(r *http.Request) error {
	if req.Username == "" || req.Password == "" {
		return fmt.Errorf("missing required session fields")
	}
	return nil
}

// NewSession POST /session
func (h SessionHandler) NewSession(w http.ResponseWriter, r *http.Request) {
	data := &reqSession{}
	if err := ReceiveJson(r, data); err != nil {
		SendError(w, r, http.StatusBadRequest, err)
		return
	}

	user := h.s.Users.FindByName(data.Username)
	if user == nil {
		SendError(w, r, http.StatusUnauthorized, fmt.Errorf("invalid login details"))
		return
	}

	if token, err := h.a.CreateToken(user.ID); err != nil {
		SendError(w, r, http.StatusUnprocessableEntity, err)
		return
	} else {
		SendJson(w, r, struct {
			Token string `json:"token"`
		}{
			Token: token,
		})
		return
	}
}

// DeleteSession DELETE /session
func (h SessionHandler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("DeleteSession"))
}
