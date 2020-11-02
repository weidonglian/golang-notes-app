package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/weidonglian/golang-notes-app/graph/gmodel"
	"github.com/weidonglian/golang-notes-app/graph/util"
	"github.com/weidonglian/golang-notes-app/model"
)

func (r *mutationResolver) AddTodo(ctx context.Context, input gmodel.AddTodoInput) (*gmodel.AddTodoPayload, error) {
	done := false
	if input.Done != nil {
		done = *input.Done
	}

	todo, err := r.store.Todos.Create(model.Todo{
		Name:   input.Name,
		Done:   done,
		NoteID: input.NoteID,
	})

	if err != nil {
		return nil, util.ErrorUnprocessableEntity
	}

	return &gmodel.AddTodoPayload{
		ID:     todo.ID,
		Name:   todo.Name,
		Done:   todo.Done,
		NoteID: todo.NoteID,
	}, nil
}

func (r *mutationResolver) UpdateTodo(ctx context.Context, input gmodel.UpdateTodoInput) (*gmodel.UpdateTodoPayload, error) {
	if r.store.Notes.FindByID(input.NoteID, util.GetUserId(ctx)) == nil {
		return nil, util.ErrorUnprocessableEntity
	}

	todo, err := r.store.Todos.Update(input.ID, input.Name, input.Done)

	if err != nil {
		return nil, util.ErrorUnprocessableEntity
	}

	return &gmodel.UpdateTodoPayload{
		ID:     todo.ID,
		Name:   todo.Name,
		Done:   todo.Done,
		NoteID: todo.NoteID,
	}, nil
}

func (r *mutationResolver) DeleteTodo(ctx context.Context, input gmodel.DeleteTodoInput) (*gmodel.DeleteTodoPayload, error) {
	if r.store.Notes.FindByID(input.NoteID, util.GetUserId(ctx)) == nil {
		return nil, util.ErrorUnprocessableEntity
	}

	id, err := r.store.Todos.Delete(input.ID, input.NoteID)
	if err != nil {
		return nil, util.ErrorUnprocessableEntity
	}

	return &gmodel.DeleteTodoPayload{
		ID:     id,
		NoteID: input.NoteID,
	}, nil
}

func (r *mutationResolver) ToggleTodo(ctx context.Context, input gmodel.ToggleTodoInput) (*gmodel.ToggleTodoPayload, error) {
	if r.store.Notes.FindByID(input.NoteID, util.GetUserId(ctx)) == nil {
		return nil, util.ErrorUnprocessableEntity
	}

	todo, err := r.store.Todos.Toggle(input.ID)

	if err != nil {
		return nil, util.ErrorUnprocessableEntity
	}

	return &gmodel.ToggleTodoPayload{
		ID:     todo.ID,
		Done:   todo.Done,
		NoteID: todo.NoteID,
	}, nil
}

func (r *queryResolver) Todos(ctx context.Context, noteID int) ([]*gmodel.Todo, error) {
	if r.store.Notes.FindByID(noteID, util.GetUserId(ctx)) == nil {
		return nil, util.ErrorUnprocessableEntity
	}

	todos := r.store.Todos.FindByNoteID(noteID)
	gtodos := make([]*gmodel.Todo, len(todos))
	for i := range todos {
		gtodos[i] = util.NewGTodo(&todos[i])
	}
	return gtodos, nil
}
