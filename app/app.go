package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/weidonglian/golang-notes-app/config"
	"github.com/weidonglian/golang-notes-app/handlers"
	"github.com/weidonglian/golang-notes-app/store"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/weidonglian/golang-notes-app/logging"
)

// App is the main application.
type App struct {
	logger *logrus.Logger
	config config.Config
	store  store.Store
}

// Serve is the core serve http
func (a *App) Serve() {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(logging.NewStructuredLogger(a.logger))
	r.Use(middleware.Recoverer)
	//r.Use(render.SetContentType(render.ContentTypeJSON)) // force response type with json
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi, hello!"))
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	r.Mount("/todos", handlers.NewTodos().Routes())
	r.Mount("/users", handlers.NewUsers().Routes())
	r.Mount("/notes", handlers.NewNotes().Routes())

	addr := fmt.Sprintf(":%v", a.config.ServerPort)
	a.logger.Infof("Listening on addr %v", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		a.logger.Fatal(err)
	}
}

// NewApp create the main application
func NewApp(logger *logrus.Logger) (*App, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	store, err := store.NewStore(cfg, logger)
	if err != nil {
		return nil, err
	}
	return &App{
		logger: logger,
		config: cfg,
		store:  store,
	}, nil
}
