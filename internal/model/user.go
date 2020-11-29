package model

const (
	UserRoleUser  = "USER"
	UserRoleAdmin = "ADMIN"
)

var TestUsers = []User{
	{
		Username: "dev",
		Password: "dev",
		Role:     UserRoleUser,
	},
	{
		Username: "admin",
		Password: "admin",
		Role:     UserRoleAdmin,
	},
	{
		Username: "test",
		Password: "test",
		Role:     UserRoleUser,
	},
}

type User struct {
	ID       int    `db:"user_id" json:"id"`
	Username string `db:"user_name" json:"username"`
	Password string `db:"user_password" json:"password"`
	Role     string `db:"user_role" json:"role"`
}
