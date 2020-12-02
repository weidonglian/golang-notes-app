package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/weidonglian/notes-app/internal/handlers/payload"
	"github.com/weidonglian/notes-app/internal/lib"
	"github.com/weidonglian/notes-app/internal/model"
	"github.com/weidonglian/notes-app/internal/store"
	"net/http"
	"strconv"
)

type TodosHandler struct {
	todosStore store.TodosStore
	notesStore store.NotesStore
}

func NewTodosHandler(store *store.Store) TodosHandler {
	return TodosHandler{
		todosStore: store.Todos,
		notesStore: store.Notes,
	}
}

func (h TodosHandler) CtxID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var todo *model.Todo
		var todoId int

		// extract noteId from the URLParam
		if idValue := chi.URLParam(r, "id"); idValue == "" {
			lib.SendErrorBadRequest(w, r, lib.ErrorMissingRequiredParams)
			return
		} else {
			if id, err := strconv.Atoi(idValue); err != nil {
				lib.SendErrorBadRequest(w, r, err)
				return
			} else {
				todoId = id
			}
		}

		userId := lib.GetUserIDFromRequest(r)
		if todo = h.todosStore.FindByID(todoId); todo == nil {
			lib.SendErrorUnprocessableEntity(w, r, fmt.Errorf("unable to find todo"))
			return
		}

		// restrict access for current user only.
		if h.notesStore.FindByID(todo.NoteID, userId) == nil {
			lib.SendErrorUnprocessableEntity(w, r, fmt.Errorf("unable to find todo for current user"))
			return
		}

		ctx := context.WithValue(r.Context(), "CtxID", todo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h TodosHandler) Create(w http.ResponseWriter, r *http.Request) {
	data := &payload.ReqTodo{}
	if err := lib.ReceiveJson(r, data); err != nil {
		lib.SendErrorBadRequest(w, r, err)
		return
	}

	if data.NoteID == nil {
		lib.SendErrorBadRequest(w, r, lib.ErrorMissingRequiredFields)
		return
	}

	if h.notesStore.FindByID(*data.NoteID, lib.GetUserIDFromRequest(r)) == nil {
		lib.SendErrorUnprocessableEntity(w, r, errors.New("provide note with id does not exist for current user"))
		return
	}

	if todo, err := h.todosStore.Create(payload.NewTodoFromReq(data)); err != nil {
		lib.SendErrorUnprocessableEntity(w, r, err)
	} else {
		lib.SendJson(w, r, todo)
	}
}

func (h TodosHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctxTodo := r.Context().Value("CtxID").(*model.Todo)
	lib.SendJson(w, r, ctxTodo)
}

func (h TodosHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	data := &payload.ReqTodo{}
	if err := lib.ReceiveJson(r, data); err != nil {
		lib.SendErrorBadRequest(w, r, err)
		return
	}
	ctxTodo := r.Context().Value("CtxID").(*model.Todo)

	if updatedTodo, err := h.todosStore.Update(ctxTodo.ID, data.Name, data.Done); err != nil {
		lib.SendErrorUnprocessableEntity(w, r, err)
	} else {
		lib.SendJson(w, r, updatedTodo)
	}
}

func (h TodosHandler) ToggleByID(w http.ResponseWriter, r *http.Request) {
	ctxTodo := r.Context().Value("CtxID").(*model.Todo)
	if todo, err := h.todosStore.Toggle(ctxTodo.ID); err != nil {
		lib.SendErrorUnprocessableEntity(w, r, err)
	} else {
		lib.SendJson(w, r, todo)
	}
}

func (h TodosHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	todo := r.Context().Value("CtxID").(*model.Todo)
	if _, err := h.todosStore.Delete(todo.ID, todo.NoteID); err != nil {
		lib.SendErrorUnprocessableEntity(w, r, err)
	} else {
		lib.SendStatus(w, r, http.StatusOK)
	}
}
