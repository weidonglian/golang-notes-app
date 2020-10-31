package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/weidonglian/golang-notes-app/graph/gmodel"
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
		return nil, err
	}

	return &gmodel.AddTodoPayload{
		Name:   todo.Name,
		Done:   todo.Done,
		NoteID: todo.NoteID,
	}, nil
}

func (r *mutationResolver) UpdateTodo(ctx context.Context, input gmodel.UpdateTodoInput) (*gmodel.UpdateTodoPayload, error) {
	todo, err := r.store.Todos.Update(input.ID, input.Name, input.Done)

	if err != nil {
		return nil, err
	}

	return &gmodel.UpdateTodoPayload{
		ID:     todo.ID,
		Name:   todo.Name,
		Done:   todo.Done,
		NoteID: todo.NoteID,
	}, nil

}

func (r *mutationResolver) DeleteTodo(ctx context.Context, input gmodel.DeleteTodoInput) (*gmodel.DeleteTodoPayload, error) {
	err := r.store.Todos.Delete(input.ID)
	if err != nil {
		return nil, err
	}

	return &gmodel.DeleteTodoPayload{
		ID:     input.ID,
		NoteID: input.NoteID,
	}, nil
}

func (r *mutationResolver) ToggleTodo(ctx context.Context, input gmodel.ToggleTodoInput) (*gmodel.ToggleTodoPayload, error) {
	todo, err := r.store.Todos.Toggle(input.ID)

	if err != nil {
		return nil, err
	}

	return &gmodel.ToggleTodoPayload{
		ID:     todo.ID,
		Done:   todo.Done,
		NoteID: todo.NoteID,
	}, nil
}

func (r *queryResolver) Todo(ctx context.Context, id int) (*gmodel.Todo, error) {
	panic(fmt.Errorf("not implemented"))
}
