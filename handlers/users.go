package handlers

import (
	"github.com/weidonglian/golang-notes-app/store"
	"net/http"

	"github.com/go-chi/chi"
)

type UsersHandler struct {
	s *store.Store
}

func NewUsersHandler(s *store.Store) UsersHandler {
	return UsersHandler{
		s: s,
	}
}

func (h UsersHandler) Routes() chi.Router {
	// Routes for /users
	r := chi.NewRouter()

	r.Post("/new", h.Create) // POST /users - create a new user
	r.Put("/password", h.ChangePassword)
	r.Get("/", h.List) // GET /users - read a list of users

	r.Route("/{id}", func(r chi.Router) {
		r.Use(h.UserCtx)            // lets have a users map, and lets actually load/manipulate
		r.Put("/", h.UpdateByID)    // PUT /users/{id} - update a single user by :id
		r.Delete("/", h.DeleteByID) // DELETE /users/{id} - delete a single user by :id
		r.Get("/", h.GetByID)       // GET /users/{id} - read a single user by :id
	})

	return r
}

func (h UsersHandler) UserCtx(next http.Handler) http.Handler {
	return next
}

func (h UsersHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("users create"))
}

func (h UsersHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("change password"))
}

func (h UsersHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("user update"))
}

func (h UsersHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("user delete"))
}

func (h UsersHandler) List(w http.ResponseWriter, r *http.Request) {

}

func (h UsersHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("user get"))
}
