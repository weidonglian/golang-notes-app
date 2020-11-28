package store

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/weidonglian/notes-app/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type UsersStore struct {
	db     *sqlx.DB
	logger *logrus.Logger
}

func NewUsersStore(ctx *Context) UsersStore {
	return UsersStore{
		db:     ctx.Session.GetDB(),
		logger: ctx.Session.Logger,
	}
}

func (i UsersStore) Create(user model.User) (*model.User, error) {
	var retUser model.User
	if hashedPassword, err := HashPassword(user.Password); err != nil {
		return nil, err
	} else {
		user.Password = hashedPassword
	}

	if user.Role == "" {
		user.Role = model.UserRoleUser
	}
	stmt, err := i.db.PrepareNamed(`
		INSERT INTO users (user_name, user_password, user_role)
		VALUES(:user_name, :user_password, :user_role)
		RETURNING *
	`)
	if err != nil {
		return nil, err
	}
	err = stmt.Get(&retUser, user)
	return &retUser, err
}

func (i UsersStore) UpdatePassword(user model.User) (int, error) {
	var id int
	if hashedPassword, err := HashPassword(user.Password); err != nil {
		return id, err
	} else {
		user.Password = hashedPassword
	}

	stmt, err := i.db.PrepareNamed(`
		UPDATE users
		SET user_password = :user_password
		WHERE user_id = :user_id
		RETURNING user_id
	`)
	if err != nil {
		return id, err
	}
	err = stmt.Get(&id, user)
	return id, err
}

// Tries to delete a user by id, and returns the number of records deleted;
func (i UsersStore) Delete(id int) error {
	_, err := i.db.Exec("DELETE FROM users WHERE user_id = $1", id)
	return err
}

// Removes all records from the table;
func (i UsersStore) DeleteAll() error {
	_, err := i.db.Exec("TRUNCATE TABLE users CASCADE")
	return err
}

// Tries to find a user from id;
func (i UsersStore) FindByID(id int) *model.User {
	user := model.User{}
	err := i.db.Get(&user, "SELECT * FROM users WHERE user_id = $1", id)
	if err != nil {
		return nil
	}
	return &user
}

// Tries to find a user from name;
func (i UsersStore) FindByName(name string) *model.User {
	user := model.User{}
	err := i.db.Get(&user, "SELECT * FROM users WHERE user_name = $1 LIMIT 1", name)
	if err != nil {
		return nil
	}
	return &user
}

func HashPassword(pwd string) (string, error) {

	// Use GenerateFromPassword to hash & salt pwd
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash), nil
}

func CheckPassword(hashedPwd string, plainPwd string) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	if err != nil {
		return false
	}
	return true
}
