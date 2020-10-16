package handlers

import (
	"github.com/weidonglian/golang-notes-app/handlers/util"
	"github.com/weidonglian/golang-notes-app/store"
	"net/http"
)

type NotesHandler struct {
	notesStore store.NotesStore
}

func NewNotesHandler(store *store.Store) NotesHandler {
	return NotesHandler{store.Notes}
}

func (h NotesHandler) CtxID(next http.Handler) http.Handler {
	return next
}

func (h NotesHandler) List(w http.ResponseWriter, r *http.Request) {
	userId := util.GetUserIDFromRequest(r)
	util.SendJson(w, r, h.notesStore.FindByUserID(userId))
}

func (h NotesHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("notes create"))
}

func (h NotesHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userId := util.GetUserIDFromRequest(r)
	util.SendJson(w, r, h.notesStore.FindByID(userId))
}

func (h NotesHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("note update"))
}

func (h NotesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("note delete"))
}

func (h NotesHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("note delete"))
}
