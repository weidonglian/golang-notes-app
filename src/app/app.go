package app

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/weidonglian/golang-notes-app/src/logging"
)

// App is the main application.
type App struct {
}

// Serve is the core serve http
func (a *App) Serve() {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(logging.LogHandler)
	r.Use(middleware.Recoverer)
	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	r.Get("/wait", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		logging.LogEntrySetField(r, "wait", true)
		w.Write([]byte("hi"))
	})
	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("oops")
	})

	http.ListenAndServe(":3000", r)
}

// NewApp create the main application
func NewApp() (*App, error) {
	return &App{}, nil
}
