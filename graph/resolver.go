package graph

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/sirupsen/logrus"
	"github.com/weidonglian/notes-app/graph/generated"
	"github.com/weidonglian/notes-app/store"
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
