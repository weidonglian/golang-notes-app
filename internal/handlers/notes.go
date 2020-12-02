package handlers

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/weidonglian/notes-app/internal/handlers/payload"
	"github.com/weidonglian/notes-app/internal/lib"
	"github.com/weidonglian/notes-app/internal/model"
	"github.com/weidonglian/notes-app/internal/store"
	"net/http"
	"strconv"
)

type NotesHandler struct {
	notesStore store.NotesStore
	todosStore store.TodosStore
}

func NewNotesHandler(store *store.Store) NotesHandler {
	return NotesHandler{store.Notes, store.Todos}
}

func (h NotesHandler) CtxID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var note *model.Note
		var noteId int

		// extract noteId from the URLParam
		if idValue := chi.URLParam(r, "id"); idValue == "" {
			lib.SendErrorBadRequest(w, r, lib.ErrorMissingRequiredParams)
			return
		} else {
			if id, err := strconv.Atoi(idValue); err != nil {
				lib.SendErrorBadRequest(w, r, err)
				return
			} else {
				noteId = id
			}
		}
		userId := lib.GetUserIDFromRequest(r)
		if note = h.notesStore.FindByID(noteId, userId); note == nil {
			lib.SendErrorUnprocessableEntity(w, r, fmt.Errorf("unable to find note with given %d", noteId))
			return
		}

		ctx := context.WithValue(r.Context(), "CtxID", note)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h NotesHandler) List(w http.ResponseWriter, r *http.Request) {
	notes := h.notesStore.FindByUserID(lib.GetUserIDFromRequest(r))
	lib.SendJson(w, r, payload.NewRespNoteArray(notes, h.todosStore))
}

func (h NotesHandler) Create(w http.ResponseWriter, r *http.Request) {
	data := &payload.ReqNote{}
	if err := lib.ReceiveJson(r, data); err != nil {
		lib.SendErrorBadRequest(w, r, err)
		return
	}
	if note, err := h.notesStore.Create(model.Note{Name: data.Name, UserID: lib.GetUserIDFromRequest(r)}); err != nil {
		lib.SendErrorUnprocessableEntity(w, r, err)
	} else {
		lib.SendJson(w, r, payload.NewRespNote(note, h.todosStore))
	}
}

func (h NotesHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctxNote := r.Context().Value("CtxID").(*model.Note)
	lib.SendJson(w, r, payload.NewRespNote(ctxNote, h.todosStore))
}

func (h NotesHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	data := &payload.ReqNote{}
	if err := lib.ReceiveJson(r, data); err != nil {
		lib.SendErrorBadRequest(w, r, err)
		return
	}
	ctxNote := r.Context().Value("CtxID").(*model.Note)
	userId := lib.GetUserIDFromRequest(r)
	if updatedNote, err := h.notesStore.Update(ctxNote.ID, data.Name, userId); err != nil {
		lib.SendErrorUnprocessableEntity(w, r, err)
	} else {
		lib.SendJson(w, r, payload.NewRespNote(updatedNote, h.todosStore))
	}
}

func (h NotesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.notesStore.DeleteAll(lib.GetUserIDFromRequest(r)); err != nil {
		lib.SendErrorUnprocessableEntity(w, r, err)
		return
	}

	lib.SendStatus(w, r, http.StatusOK)
}

func (h NotesHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	note := r.Context().Value("CtxID").(*model.Note)
	if _, err := h.notesStore.Delete(note.ID, lib.GetUserIDFromRequest(r)); err != nil {
		lib.SendErrorUnprocessableEntity(w, r, err)
		return
	}

	lib.SendJson(w, r, http.StatusOK)
}
