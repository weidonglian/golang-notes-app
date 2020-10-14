package store

import (
	"github.com/weidonglian/golang-notes-app/config"
	"github.com/weidonglian/golang-notes-app/db"
	"github.com/weidonglian/golang-notes-app/model"
)

type Context struct {
	Session *db.Session
}

type Store struct {
	ctx   *Context
	Users UsersStore
	Notes NotesStore
	Todos TodosStore
}

func NewStore(sess *db.Session) (*Store, error) {
	ctx := Context{
		Session: sess,
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
