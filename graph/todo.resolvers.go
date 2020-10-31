package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/weidonglian/golang-notes-app/graph/gmodel"
)

func (r *mutationResolver) AddTodo(ctx context.Context, input *gmodel.AddTodoInput) (*gmodel.Todo, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateTodo(ctx context.Context, input *gmodel.UpdateTodoInput) (*gmodel.Todo, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteTodo(ctx context.Context, id int) (*gmodel.Todo, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ToggleTodo(ctx context.Context, id int) (*gmodel.Todo, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Todo(ctx context.Context, id string) (*gmodel.Todo, error) {
	panic(fmt.Errorf("not implemented"))
}
