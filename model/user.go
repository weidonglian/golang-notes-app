package model

type User struct {
	ID       int    `db:"user_id" json:"id"`
	Username string `db:"user_name" json:"username"`
	Password string `db:"user_password" json:"password"`
	Role     string `db:"user_role" json:"role"`
}

func NewUser() User {
	return User{}
}
