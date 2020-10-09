package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/weidonglian/golang-notes-app/handlers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/weidonglian/golang-notes-app/logging"
)

// App is the main application.
type App struct {
	logger *logrus.Logger
}

// Serve is the core serve http
func (a *App) Serve() {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(logging.NewStructuredLogger(a.logger))
	r.Use(middleware.Recoverer)
	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Mount("/todos", handlers.NewTodos().Routes())
	r.Mount("/users", handlers.NewUsers().Routes())
	r.Mount("/notes", handlers.NewNotes().Routes())

	port := 3000
	addr := fmt.Sprintf(":%v", port)
	a.logger.Infof("Listening on addr %v", addr)
	http.ListenAndServe(addr, r)
}

// NewApp create the main application
func NewApp() (*App, error) {
	// Setup the logger backend using sirupsen/logrus and configure
	// it to use a custom JSONFormatter. See the logrus docs for how to
	// configure the backend at github.com/sirupsen/logrus
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		//FullTimestamp: true,
		DisableTimestamp: true,
	})
	/* logger.Formatter = &logrus.JSONFormatter{
		// disable, as we set our own
		DisableTimestamp: true,
	}*/
	return &App{logger}, nil
}
