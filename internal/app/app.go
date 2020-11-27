package app

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/weidonglian/notes-app/internal/db"
	"net/http"

	"github.com/weidonglian/notes-app/config"
	"github.com/weidonglian/notes-app/internal/auth"
	"github.com/weidonglian/notes-app/internal/handlers"
	"github.com/weidonglian/notes-app/internal/store"

	"github.com/go-chi/chi"
)

// App is the main application.
type App struct {
	logger *logrus.Logger
	config config.Config
	db     *db.Session
	store  *store.Store
	auth   *auth.Auth
}

func (a *App) Router() *chi.Mux {
	return handlers.NewRouter(a.logger, a.auth, a.store)
}

func (a *App) GetStore() *store.Store {
	return a.store
}

func (a *App) GetAuth() *auth.Auth {
	return a.auth
}

// Serve is the core serve http
func (a *App) Serve() {
	r := a.Router()
	addr := fmt.Sprintf(":%v", a.config.ServerPort)
	a.logger.Infof("Listening on addr %v", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		a.logger.Fatal(err)
	}
}

func (a *App) Close() {
	a.db.Close()
}

// NewApp create the main application
func NewApp(logger *logrus.Logger, cfg config.Config) (*App, error) {
	if sess, err := db.NewSession(logger, cfg); err != nil {
		return nil, err
	} else {
		return NewAppWith(logger, cfg, sess)
	}
}

func NewAppWith(logger *logrus.Logger, cfg config.Config, db *db.Session) (*App, error) {
	if sto, err := store.NewStore(db); err != nil {
		return nil, err
	} else {
		return &App{
			logger: logger,
			config: cfg,
			db:     db,
			auth:   auth.NewAuth(cfg),
			store:  sto,
		}, nil
	}

}
