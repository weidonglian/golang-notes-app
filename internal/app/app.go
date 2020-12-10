package app

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/weidonglian/notes-app/internal/db"
	"github.com/weidonglian/notes-app/internal/pubsub"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/weidonglian/notes-app/config"
	"github.com/weidonglian/notes-app/internal/auth"
	"github.com/weidonglian/notes-app/internal/handlers"
	"github.com/weidonglian/notes-app/internal/store"

	"github.com/go-chi/chi"
)

// App is the main application.
type App struct {
	logger     *logrus.Logger
	config     config.Config
	db         db.Session
	store      *store.Store
	auth       *auth.Auth
	publisher  pubsub.Publisher
	subscriber pubsub.Subscriber
}

func (a *App) Router() *chi.Mux {
	return handlers.NewRouter(a.logger, a.auth, a.store, a.publisher)
}

func (a *App) GetStore() *store.Store {
	return a.store
}

func (a *App) GetAuth() *auth.Auth {
	return a.auth
}

// Serve is the core serve http
func (a *App) Serve() error {
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

	shutdownChan := make(chan os.Signal)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)
	<-shutdownChan

	a.logger.Info("Got signal to shutdown, draining connections and quitting app")

	ctxWait, cancelWait := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelWait()

	if err := srv.Shutdown(ctxWait); err != nil {
		a.logger.Fatalf("server Shutdown Failed:%+s", err)
		return err
	}

	a.Close()

	a.logger.Info("server exited properly")

	return nil
}

func (a *App) Close() {
	a.publisher.Close()
	a.subscriber.Close()
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

func NewAppWith(logger *logrus.Logger, cfg config.Config, db db.Session) (*App, error) {
	var (
		sto        *store.Store
		publisher  pubsub.Publisher
		subscriber pubsub.Subscriber
		err        error
	)

	if publisher, err = pubsub.NewPubClient(logger, &cfg); err != nil {
		return nil, err
	}

	if subscriber, err = pubsub.NewSubClient(logger, &cfg); err != nil {
		return nil, err
	}

	if sto, err = store.NewStore(db, logger); err != nil {
		return nil, err
	}

	counter := 1
	subscriber.Subscribe(context.Background(), "app.entity.*", func(msg *nats.Msg) {
		counter += 1
		logger.Printf("[#%d] Received on [%s]: '%s'", counter, msg.Subject, string(msg.Data))
	})

	return &App{
		logger:     logger,
		config:     cfg,
		db:         db,
		auth:       auth.NewAuth(cfg),
		store:      sto,
		publisher:  publisher,
		subscriber: subscriber,
	}, nil
}
