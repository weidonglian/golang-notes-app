package app

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/weidonglian/notes-app/internal/db"
	"github.com/weidonglian/notes-app/pkg/util"
	"net/http"
	"time"

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
func (a *App) Serve() error {
	ctx := util.NewShutdownContext()

	r := a.Router()
	addr := fmt.Sprintf(":%v", a.config.ServerPort)

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		a.logger.Infof("Listening on addr %v", addr)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			a.logger.Fatal(err)
		}
	}()

	<-ctx.Done()

	a.logger.Info("Server is stopping")

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctxTimeout); err != nil {
		a.logger.Fatalf("server Shutdown Failed:%+s", err)
		return err
	}

	a.logger.Info("server exited properly")

	a.Close()

	return nil
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
