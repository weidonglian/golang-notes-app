package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/weidonglian/notes-app/internal/graphql/generated"
)

func (r *mutationResolver) PlaceHolder(ctx context.Context) (*bool, error) {
	dummy := true
	return &dummy, nil
}

func (r *queryResolver) PlaceHolder(ctx context.Context) (*bool, error) {
	dummy := true
	return &dummy, nil
}

func (r *subscriptionResolver) PlaceHolder(ctx context.Context) (<-chan *bool, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
