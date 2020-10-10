package model

import "time"

type User struct {
	ID        int       `db:"user_id" json:"id"`
	Username  string    `db:"user_name" json:"username"`
	Password  string    `db:"user_password" json:"password"`
	Role      string    `db:"user_role" json:"role"`
	CreatedAt time.Time `db:"user_created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"user_updated_at" json:"updatedAt"`
}

func NewUser() User {
	return User{}
}
