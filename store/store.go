package store

import (
	"github.com/sirupsen/logrus"
	"github.com/weidonglian/golang-notes-app/config"
	"github.com/weidonglian/golang-notes-app/db"
)

type StoreContext struct {
	db *db.Session
}

type Store struct {
	ctx   *StoreContext
	Users Users
	Notes Notes
	Todos Todos
}

func NewStore(cfg config.Config, logger *logrus.Logger) (*Store, error) {
	sess, err := db.NewSession(logger, cfg)
	if err != nil {
		return nil, err
	}

	ctx := StoreContext{sess}

	return &Store{
		ctx:   &ctx,
		Users: NewUsers(&ctx),
		Notes: NewNotes(&ctx),
		Todos: NewTodos(&ctx),
	}, nil
}
