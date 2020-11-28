package store

import (
	"github.com/sirupsen/logrus"
	"github.com/weidonglian/notes-app/internal/db"
)

type Context struct {
	Session db.Session
	Logger  *logrus.Logger
}

type Store struct {
	ctx   *Context
	Users UsersStore
	Notes NotesStore
	Todos TodosStore
}

func NewStore(sess db.Session, logger *logrus.Logger) (*Store, error) {
	ctx := Context{
		Session: sess,
		Logger:  logger,
	}

	s := Store{
		ctx:   &ctx,
		Users: NewUsersStore(&ctx),
		Notes: NewNotesStore(&ctx),
		Todos: NewTodosStore(&ctx),
	}
	return &s, nil
}
