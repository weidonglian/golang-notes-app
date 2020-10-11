package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/weidonglian/golang-notes-app/errors"
)

// SessionHandler keeps the dependency for handle session.
type SessionHandler struct {
}

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
		return
	}

	render.JSON(w, r, data)
}

func (h SessionHandler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("DeleteSession"))
}

func NewSessionHandler() SessionHandler {
	return SessionHandler{}
}
