package app

import (
	"github.com/weidonglian/golang-notes-app/config"
	"github.com/weidonglian/golang-notes-app/model"
)

func DataInit(a *App) {
	if config.IsDevMode() {
		devUsers := []model.User{
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
		for _, user := range devUsers {
			if a.store.Users.FindByName(user.Username) == nil {
				a.store.Users.Create(user)
			}
		}
	}
}
