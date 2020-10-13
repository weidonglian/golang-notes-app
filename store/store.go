package store

import (
	"github.com/weidonglian/golang-notes-app/config"
	"github.com/weidonglian/golang-notes-app/db"
	"github.com/weidonglian/golang-notes-app/model"
)

type StoreContext struct {
	Session *db.Session
}

type Store struct {
	ctx   *StoreContext
	Users UsersStore
	Notes NotesStore
	Todos TodosStore
}

func NewStore(sess *db.Session) (*Store, error) {
	ctx := StoreContext{
		Session: sess,
	}

	s := Store{
		ctx:   &ctx,
		Users: NewUsersStore(&ctx),
		Notes: NewNotesStore(&ctx),
		Todos: NewTodosStore(&ctx),
	}
	if !config.IsProdMode() {
		s.addTestUsers()
	}
	return &s, nil
}

func (s Store) addTestUsers() {
	if config.IsProdMode() {
		panic("TestUsers should never be used in production mode")
	}

	testUsers := []model.User{
		{
			Username: "dev",
			Password: "dev",
			Role:     model.UserRoleUser,
		},
		{
			Username: "admin",
			Password: "admin",
			Role:     model.UserRoleAdmin,
		},
		{
			Username: "test",
			Password: "test",
			Role:     model.UserRoleUser,
		},
	}

	for _, user := range testUsers {
		if s.Users.FindByName(user.Username) == nil {
			s.Users.Create(user)
		}
	}
}
