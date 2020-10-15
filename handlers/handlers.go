package handlers

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
	"github.com/weidonglian/golang-notes-app/auth"
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

		// session handler
		session := NewSessionHandler(store, auth)
		r.Delete("/session", session.DeleteSession)
		// users handler
		users := NewUsersHandler(store)
		r.Put("/users/password", users.ChangePassword)
		r.Route("/users/{id}", func(r chi.Router) {
			r.Use(users.CtxID)
			r.Put("/", users.UpdateByID)
			r.Delete("/", users.DeleteByID)
			r.Get("/", users.GetByID)
		})
		// notes handler
		notes := NewNotesHandler(store)
		r.Post("/notes", notes.Create)
		r.Delete("/notes", notes.Delete)
		r.Get("/notes", notes.List)
		r.Route("/notes/{id}", func(r chi.Router) {
			r.Use(notes.CtxID)
			r.Put("/", notes.UpdateByID)
			r.Delete("/", notes.DeleteByID)
			r.Get("/", notes.GetByID)
		})
		// todos handler
		//todos := NewTodosHandler()

	})
	return r
}
