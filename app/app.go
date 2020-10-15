package app

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/weidonglian/golang-notes-app/db"
	"net/http"

	"github.com/weidonglian/golang-notes-app/auth"
	"github.com/weidonglian/golang-notes-app/config"
	"github.com/weidonglian/golang-notes-app/handlers"
	"github.com/weidonglian/golang-notes-app/store"

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

func NewTestApp(logger *logrus.Logger, cfg config.Config) (*App, error) {
	if !config.IsTestMode() {
		panic("NewTestApp should only be used for test application")
	}

	dbSess := db.LoadSessionPool(logger, cfg).ForkNewSession()

	var s *store.Store
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
