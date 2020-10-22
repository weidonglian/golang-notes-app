package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/weidonglian/golang-notes-app/handlers/payload"
	"github.com/weidonglian/golang-notes-app/handlers/util"
	"github.com/weidonglian/golang-notes-app/model"
	"github.com/weidonglian/golang-notes-app/store"
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
			util.SendErrorBadRequest(w, r, util.ErrorMissingRequiredParams)
			return
		} else {
			if id, err := strconv.Atoi(idValue); err != nil {
				util.SendErrorBadRequest(w, r, err)
				return
			} else {
				todoId = id
			}
		}

		userId := util.GetUserIDFromRequest(r)
		if todo = h.todosStore.FindByID(todoId); todo == nil {
			util.SendErrorUnprocessableEntity(w, r, fmt.Errorf("unable to find todo"))
			return
		}

		// restrict access for current user only.
		if h.notesStore.FindByID(todo.NoteID, userId) == nil {
			util.SendErrorUnprocessableEntity(w, r, fmt.Errorf("unable to find todo for current user"))
			return
		}

		ctx := context.WithValue(r.Context(), "CtxID", todo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h TodosHandler) Create(w http.ResponseWriter, r *http.Request) {
	data := &payload.ReqTodo{}
	if err := util.ReceiveJson(r, data); err != nil {
		util.SendErrorBadRequest(w, r, err)
		return
	}

	if data.NoteID == nil {
		util.SendErrorBadRequest(w, r, util.ErrorMissingRequiredFields)
		return
	}

	if h.notesStore.FindByID(*data.NoteID, util.GetUserIDFromRequest(r)) == nil {
		util.SendErrorUnprocessableEntity(w, r, errors.New("provide note with id does not exist for current user"))
		return
	}

	if todo, err := h.todosStore.Create(payload.NewTodoFromReq(data)); err != nil {
		util.SendErrorUnprocessableEntity(w, r, err)
	} else {
		util.SendJson(w, r, todo)
	}
}

func (h TodosHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctxTodo := r.Context().Value("CtxID").(*model.Todo)
	util.SendJson(w, r, ctxTodo)
}

func (h TodosHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	data := &payload.ReqTodo{}
	if err := util.ReceiveJson(r, data); err != nil {
		util.SendErrorBadRequest(w, r, err)
		return
	}
	ctxTodo := r.Context().Value("CtxID").(*model.Todo)

	if updatedTodo, err := h.todosStore.Update(ctxTodo.ID, data.Name, data.Done); err != nil {
		util.SendErrorUnprocessableEntity(w, r, err)
	} else {
		util.SendJson(w, r, updatedTodo)
	}
}

func (h TodosHandler) ToggleByID(w http.ResponseWriter, r *http.Request) {
	ctxTodo := r.Context().Value("CtxID").(*model.Todo)
	if todo, err := h.todosStore.Toggle(ctxTodo.ID); err != nil {
		util.SendErrorUnprocessableEntity(w, r, err)
	} else {
		util.SendJson(w, r, todo)
	}
}

func (h TodosHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	todo := r.Context().Value("CtxID").(*model.Todo)
	if err := h.todosStore.Delete(todo.ID); err != nil {
		util.SendErrorUnprocessableEntity(w, r, err)
	} else {
		util.SendStatus(w, r, http.StatusOK)
	}
}
