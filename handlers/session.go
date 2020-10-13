package handlers

import (
	"fmt"
	"github.com/weidonglian/golang-notes-app/auth"
	"github.com/weidonglian/golang-notes-app/store"
	"net/http"

	"github.com/go-chi/render"
	"github.com/weidonglian/golang-notes-app/errors"
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

func (h SessionHandler) NewSession(w http.ResponseWriter, r *http.Request) {
	/*var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	//compare the user from the request, with the one we defined:
	if user.Username != u.Username || user.Password != u.Password {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}
	token, err := CreateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	c.JSON(http.StatusOK, token)*/

	data := &reqSession{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, errors.ErrBadRequest(err))
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

func (h SessionHandler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("DeleteSession"))
}
