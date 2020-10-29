package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/weidonglian/golang-notes-app/graph/generated"
	"github.com/weidonglian/golang-notes-app/graph/model"
)

func (r *mutationResolver) CreateNote(ctx context.Context, name string) (*model.Note, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateNote(ctx context.Context, id int, name string) (*model.Note, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteNote(ctx context.Context, id int) (*model.Note, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateTodo(ctx context.Context, name string, noteID int) (*model.Todo, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateTodo(ctx context.Context, id int, name string) (*model.Todo, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteTodo(ctx context.Context, id int) (*model.Todo, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ToggleTodo(ctx context.Context, id int) (*model.Todo, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Notes(ctx context.Context) ([]*model.Note, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
