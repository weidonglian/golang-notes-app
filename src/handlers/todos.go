package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
)

type TodosHandler struct{}

func NewTodos() TodosHandler {
	return TodosHandler{}
}

// Routes creates a REST router for the todos resource
func (h TodosHandler) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..

	r.Get("/", h.List)    // GET /todos - read a list of todos
	r.Post("/", h.Create) // POST /todos - create a new todo and pehist it
	r.Put("/", h.Delete)

	r.Route("/{id}", func(r chi.Router) {
		// r.Use(h.TodoCtx) // lets have a todos map, and lets actually load/manipulate
		r.Get("/", h.Get)       // GET /todos/{id} - read a single todo by :id
		r.Put("/", h.Update)    // PUT /todos/{id} - update a single todo by :id
		r.Delete("/", h.Delete) // DELETE /todos/{id} - delete a single todo by :id
	})

	return r
}

func (h TodosHandler) List(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todos list of stuff.."))
}

func (h TodosHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todos create"))
}

func (h TodosHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todo get"))
}

func (h TodosHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todo update"))
}

func (h TodosHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todo delete"))
}
