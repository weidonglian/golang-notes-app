package handlers

import (
	"fmt"
	"github.com/weidonglian/golang-notes-app/model"
	"github.com/weidonglian/golang-notes-app/store"
	"net/http"
)

type UsersHandler struct {
	usersStore *store.UsersStore
}

func NewUsersHandler(s *store.Store) UsersHandler {
	return UsersHandler{
		usersStore: &s.Users,
	}
}

func (h UsersHandler) UserCtx(next http.Handler) http.Handler {
	return next
}

type reqNewUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (req reqNewUser) Bind(r *http.Request) error {
	if req.Username == "" || req.Password == "" {
		return fmt.Errorf("missing required fields")
	}
	return nil
}

func (h UsersHandler) Create(w http.ResponseWriter, r *http.Request) {
	data := &reqNewUser{}

	if err := ReceiveJson(r, data); err != nil {
		SendError(w, r, http.StatusBadRequest, err)
		return
	}

	if h.usersStore.FindByName(data.Username) != nil {
		SendError(w, r, http.StatusBadRequest, fmt.Errorf("username '%s' already exists", data.Username))
		return
	}

	newUser := model.User{
		Username: data.Username,
		Password: data.Password,
		Role:     model.UserRoleUser,
	}

	if _, err := h.usersStore.Create(newUser); err != nil {
		SendError(w, r, http.StatusInternalServerError, err)
		return
	} else {
		SendStatus(w, r, http.StatusOK)
		return
	}
}

func (h UsersHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("NotImplementedYet"))
}

func (h UsersHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("NotImplementedYet"))
}

func (h UsersHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("NotImplementedYet"))
}

func (h UsersHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("NotImplementedYet"))
}
