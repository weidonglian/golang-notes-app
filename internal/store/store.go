package store

import (
	"github.com/sirupsen/logrus"
	"github.com/weidonglian/notes-app/config"
	"github.com/weidonglian/notes-app/internal/db"
	"github.com/weidonglian/notes-app/internal/model"
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
	if !config.IsProdMode() {
		s.loadTestData()
	}
	return &s, nil
}

func (s Store) loadTestData() {
	if config.IsProdMode() {
		panic("TestUsers should never be used in production mode")
	}

	for _, user := range model.TestUsers {
		if s.Users.FindByName(user.Username) == nil {
			s.Users.Create(user)
		}
	}
}
