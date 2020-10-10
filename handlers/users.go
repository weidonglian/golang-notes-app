package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
)

type UsersHandler struct{}

func NewUsersHandler() UsersHandler {
	return UsersHandler{}
}

// Routes creates a REST router for the users resource
func (h UsersHandler) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..

	r.Get("/", h.List)    // GET /users - read a list of users
	r.Post("/", h.Create) // POST /users - create a new user and pehist it
	r.Put("/", h.Delete)
	r.Post("/password", h.ChangePassword)
	r.Post("/new", h.NewUser) // /users/new - signup

	r.Route("/{id}", func(r chi.Router) {
		// r.Use(h.UserCtx) // lets have a users map, and lets actually load/manipulate
		r.Get("/", h.Get)       // GET /users/{id} - read a single user by :id
		r.Put("/", h.Update)    // PUT /users/{id} - update a single user by :id
		r.Delete("/", h.Delete) // DELETE /users/{id} - delete a single user by :id
	})

	return r
}

func (h UsersHandler) List(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("users list of stuff.."))
}

func (h UsersHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("users create"))
}

func (h UsersHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("user get"))
}

func (h UsersHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("user update"))
}

func (h UsersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("user delete"))
}

func (h UsersHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("change password"))
}

func (h UsersHandler) NewUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("NewUser"))
}
