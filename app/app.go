package app

import (
	"fmt"
	"github.com/weidonglian/golang-notes-app/db"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/weidonglian/golang-notes-app/auth"
	"github.com/weidonglian/golang-notes-app/config"
	"github.com/weidonglian/golang-notes-app/handlers"
	"github.com/weidonglian/golang-notes-app/store"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/weidonglian/golang-notes-app/logging"
)

// App is the main application.
type App struct {
	logger *logrus.Logger
	config config.Config
	db     *db.Session
	store  *store.Store
	auth   *auth.Auth
}

// Serve is the core serve http
func (a *App) Serve() {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(logging.NewStructuredLogger(a.logger))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	sessionHandler := handlers.NewSessionHandler(a.store, a.auth)
	// public routes and no auth required
	r.Group(func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hi, hello!"))
		})
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("pong"))
		})
		r.Post("/session", sessionHandler.NewSession) // i.e. login
	})

	// Protected routes and auth required
	r.Group(func(r chi.Router) {
		// middlewares for protected routes
		r.Use(a.auth.Verifier())
		r.Use(a.auth.Authenticator())
		r.Use(render.SetContentType(render.ContentTypeJSON)) // force response type with json

		r.Delete("/session", sessionHandler.DeleteSession) // i.e. logout
		r.Mount("/todos", handlers.NewTodosHandler().Routes())
		r.Mount("/users", handlers.NewUsersHandler().Routes())
		r.Mount("/notes", handlers.NewNotesHandler().Routes())
	})

	addr := fmt.Sprintf(":%v", a.config.ServerPort)
	a.logger.Infof("Listening on addr %v", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		a.logger.Fatal(err)
	}
}

func (a *App) Shutdown() {
	a.db.Close()
}

// NewApp create the main application
func NewApp(logger *logrus.Logger) (*App, error) {
	cfg := config.GetConfig()

	var (
		dbSess *db.Session
		s      *store.Store
	)
	if sess, err := db.NewSession(logger, cfg); err != nil {
		return nil, err
	} else {
		dbSess = sess
	}

	if sto, err := store.NewStore(dbSess); err != nil {
		return nil, err
	} else {
		s = sto
	}

	a := &App{
		logger: logger,
		config: cfg,
		db:     dbSess,
		auth:   auth.NewAuth(cfg),
		store:  s,
	}
	return a, nil
}
