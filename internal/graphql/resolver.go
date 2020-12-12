package graphql

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/weidonglian/notes-app/internal/graphql/generated"
	"github.com/weidonglian/notes-app/internal/middleware"
	"github.com/weidonglian/notes-app/internal/pubsub"
	"github.com/weidonglian/notes-app/internal/store"
	"net/http"
	"time"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	logger     *logrus.Logger
	store      *store.Store
	publisher  pubsub.Publisher
	subscriber pubsub.Subscriber
}

func NewGraphQLHandler(logger *logrus.Logger, store *store.Store, publisher pubsub.Publisher, subscriber pubsub.Subscriber) http.Handler {
	srv := handler.New(generated.NewExecutableSchema(generated.Config{
		Resolvers: &Resolver{
			logger:     logger,
			store:      store,
			publisher:  publisher,
			subscriber: subscriber,
		},
	}))

	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			HandshakeTimeout: 45 * time.Second,
			CheckOrigin:      middleware.CheckOriginFunc,
		},
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})
	srv.SetQueryCache(lru.New(1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})
	return srv
}
