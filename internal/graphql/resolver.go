package graphql

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/go-chi/jwtauth"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/weidonglian/notes-app/internal/auth"
	"github.com/weidonglian/notes-app/internal/graphql/generated"
	"github.com/weidonglian/notes-app/internal/middleware"
	"github.com/weidonglian/notes-app/internal/pubsub"
	"github.com/weidonglian/notes-app/internal/store"
	"net/http"
	"strings"
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

func NewGraphQLHandler(logger *logrus.Logger, auth *auth.Auth, store *store.Store, publisher pubsub.Publisher, subscriber pubsub.Subscriber) http.Handler {
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
		InitFunc: func(ctx context.Context, initPayload transport.InitPayload) (context.Context, error) {
			if bearer := initPayload.Authorization(); bearer != "" {
				token := ""
				if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
					token = bearer[7:]
				} else {
					return nil, fmt.Errorf("invalid bearer token")
				}
				// verification
				jwtToken, err := auth.VerifyToken(token)
				if err != nil {
					return nil, err
				}
				ctx = jwtauth.NewContext(ctx, jwtToken, err)

				// authentication
				ctxToken, _, err := jwtauth.FromContext(ctx)
				if err != nil {
					return nil, err
				}
				if ctxToken == nil || !ctxToken.Valid {
					return nil, fmt.Errorf("invalid authorization")
				}
				return ctx, nil
			}
			return nil, fmt.Errorf("authentication required")
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
