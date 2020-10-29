package graph

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/sirupsen/logrus"
	"github.com/weidonglian/golang-notes-app/graph/generated"
	"github.com/weidonglian/golang-notes-app/store"
	"net/http"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	logger *logrus.Logger
	store  *store.Store
}

func NewGraphQLHandler(logger *logrus.Logger, store *store.Store) http.Handler {
	return handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &Resolver{logger, store}}))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
