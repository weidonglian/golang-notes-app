package graphql

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/sirupsen/logrus"
	"github.com/weidonglian/notes-app/internal/graphql/generated"
	"github.com/weidonglian/notes-app/internal/pubsub"
	"github.com/weidonglian/notes-app/internal/store"
	"net/http"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	logger    *logrus.Logger
	store     *store.Store
	publisher pubsub.Publisher
}

func NewGraphQLHandler(logger *logrus.Logger, store *store.Store, publisher pubsub.Publisher) http.Handler {
	return handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &Resolver{
			logger:    logger,
			store:     store,
			publisher: publisher,
		},
	}))
}
