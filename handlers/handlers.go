package handlers

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
	"github.com/weidonglian/golang-notes-app/auth"
	"github.com/weidonglian/golang-notes-app/graph"
	"github.com/weidonglian/golang-notes-app/graph/generated"
	"github.com/weidonglian/golang-notes-app/handlers/util"
	"github.com/weidonglian/golang-notes-app/logging"
	"github.com/weidonglian/golang-notes-app/store"
	"net/http"
	"time"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func NewRouter(logger *logrus.Logger, auth *auth.Auth, store *store.Store) *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(util.CorsOptions))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(logging.NewStructuredLogger(logger))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// public routes and no auth required
	r.Group(func(r chi.Router) {
		// root index
		r.Get("/", rootHandler)
		// ping vs pong
		r.Get("/ping", pingHandler)
		// session
		session := NewSessionHandler(store, auth)
		r.Post("/session", session.NewSession)
		// user
		users := NewUsersHandler(store)
		r.Post("/users/new", users.Create)
	})

	// Protected routes and auth required
	r.Group(func(r chi.Router) {
		// middlewares for protected routes
		r.Use(auth.Verifier())
		r.Use(auth.Authenticator())

		// playground for graphql api
		r.Handle("/playground", playground.Handler("GraphQL playground", "/graphql"))
		// Graphql handler
		r.Handle("/graphql", handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}})))

		// session handler
		session := NewSessionHandler(store, auth)
		r.Delete("/session", session.DeleteSession)
		// users handler
		users := NewUsersHandler(store)
		r.Put("/users/password", users.ChangePassword)
		r.Route("/users/{id}", func(r chi.Router) {
			r.Use(users.CtxID)
			r.Get("/", users.GetByID)
			r.Put("/", users.UpdateByID)
			r.Delete("/", users.DeleteByID)
		})

		// notes handler
		notes := NewNotesHandler(store)
		r.Get("/notes", notes.List)
		r.Post("/notes", notes.Create)
		r.Delete("/notes", notes.Delete)
		r.Route("/notes/{id}", func(r chi.Router) {
			r.Use(notes.CtxID)
			r.Get("/", notes.GetByID)
			r.Put("/", notes.UpdateByID)
			r.Delete("/", notes.DeleteByID)
		})
		// todos handler
		todos := NewTodosHandler(store)
		r.Post("/todos", todos.Create)
		r.Route("/todos/{id}", func(r chi.Router) {
			r.Use(todos.CtxID)
			r.Get("/", todos.GetByID)
			r.Put("/", todos.UpdateByID)
			r.Put("/toggle", todos.ToggleByID)
			r.Delete("/", todos.DeleteByID)
		})

	})
	return r
}
