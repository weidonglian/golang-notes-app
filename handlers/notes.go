package handlers

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/weidonglian/notes-app/handlers/payload"
	"github.com/weidonglian/notes-app/handlers/util"
	"github.com/weidonglian/notes-app/model"
	"github.com/weidonglian/notes-app/store"
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
			util.SendErrorBadRequest(w, r, util.ErrorMissingRequiredParams)
			return
		} else {
			if id, err := strconv.Atoi(idValue); err != nil {
				util.SendErrorBadRequest(w, r, err)
				return
			} else {
				noteId = id
			}
		}
		userId := util.GetUserIDFromRequest(r)
		if note = h.notesStore.FindByID(noteId, userId); note == nil {
			util.SendErrorUnprocessableEntity(w, r, fmt.Errorf("unable to find note with given %d", noteId))
			return
		}

		ctx := context.WithValue(r.Context(), "CtxID", note)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h NotesHandler) List(w http.ResponseWriter, r *http.Request) {
	notes := h.notesStore.FindByUserID(util.GetUserIDFromRequest(r))
	util.SendJson(w, r, payload.NewRespNoteArray(notes, h.todosStore))
}

func (h NotesHandler) Create(w http.ResponseWriter, r *http.Request) {
	data := &payload.ReqNote{}
	if err := util.ReceiveJson(r, data); err != nil {
		util.SendErrorBadRequest(w, r, err)
		return
	}
	if note, err := h.notesStore.Create(model.Note{Name: data.Name, UserID: util.GetUserIDFromRequest(r)}); err != nil {
		util.SendErrorUnprocessableEntity(w, r, err)
	} else {
		util.SendJson(w, r, payload.NewRespNote(note, h.todosStore))
	}
}

func (h NotesHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctxNote := r.Context().Value("CtxID").(*model.Note)
	util.SendJson(w, r, payload.NewRespNote(ctxNote, h.todosStore))
}

func (h NotesHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	data := &payload.ReqNote{}
	if err := util.ReceiveJson(r, data); err != nil {
		util.SendErrorBadRequest(w, r, err)
		return
	}
	ctxNote := r.Context().Value("CtxID").(*model.Note)
	userId := util.GetUserIDFromRequest(r)
	if updatedNote, err := h.notesStore.Update(ctxNote.ID, data.Name, userId); err != nil {
		util.SendErrorUnprocessableEntity(w, r, err)
	} else {
		util.SendJson(w, r, payload.NewRespNote(updatedNote, h.todosStore))
	}
}

func (h NotesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.notesStore.DeleteAll(util.GetUserIDFromRequest(r)); err != nil {
		util.SendErrorUnprocessableEntity(w, r, err)
		return
	}

	util.SendStatus(w, r, http.StatusOK)
}

func (h NotesHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	note := r.Context().Value("CtxID").(*model.Note)
	if _, err := h.notesStore.Delete(note.ID, util.GetUserIDFromRequest(r)); err != nil {
		util.SendErrorUnprocessableEntity(w, r, err)
		return
	}

	util.SendJson(w, r, http.StatusOK)
}
