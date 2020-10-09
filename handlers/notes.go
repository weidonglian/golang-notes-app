package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
)

type NotesHandler struct{}

func NewNotes() NotesHandler {
	return NotesHandler{}
}

// Routes creates a REST router for the notes resource
func (h NotesHandler) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..

	r.Get("/", h.List)    // GET /notes - read a list of notes
	r.Post("/", h.Create) // POST /notes - create a new note and pehist it
	r.Put("/", h.Delete)

	r.Route("/{id}", func(r chi.Router) {
		// r.Use(h.NoteCtx) // lets have a notes map, and lets actually load/manipulate
		r.Get("/", h.Get)       // GET /notes/{id} - read a single note by :id
		r.Put("/", h.Update)    // PUT /notes/{id} - update a single note by :id
		r.Delete("/", h.Delete) // DELETE /notes/{id} - delete a single note by :id
	})

	return r
}

func (h NotesHandler) List(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("notes list of stuff.."))
}

func (h NotesHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("notes create"))
}

func (h NotesHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("note get"))
}

func (h NotesHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("note update"))
}

func (h NotesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("note delete"))
}
