package handlers_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/weidonglian/golang-notes-app/model"
	"net/http"
)

var _ = Describe("Users", func() {
	It("POST /users/new", func() {
		By("should not be able to create any existing users")
		for _, user := range model.TestUsers {
			testApp.API.POST("/users/new").WithJSON(map[string]string{"username": user.Username, "password": user.Password}).
				Expect().
				Status(http.StatusBadRequest).Body().Contains("already exists")
		}

		By("should be able to create users")
		testApp.API.POST("/users/new").WithJSON(map[string]string{"username": "u1", "password": "p1"}).
			Expect().
			Status(http.StatusOK)
		testApp.API.POST("/users/new").WithJSON(map[string]string{"username": "u2", "password": "p2"}).
			Expect().
			Status(http.StatusOK)
	})
})
