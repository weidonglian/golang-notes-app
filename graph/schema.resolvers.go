package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/weidonglian/notes-app/graph/generated"
)

func (r *mutationResolver) PlaceHolder(ctx context.Context) (*bool, error) {
	dummy := true
	return &dummy, nil
}

func (r *queryResolver) PlaceHolder(ctx context.Context) (*bool, error) {
	dummy := true
	return &dummy, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
