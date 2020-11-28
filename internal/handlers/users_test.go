package handlers_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/weidonglian/notes-app/internal/model"
	. "github.com/weidonglian/notes-app/internal/test"
	"net/http"
)

var _ = Describe("Users", func() {

	var testApp MockApp

	BeforeEach(func() {
		testApp = NewMockApp()
	})

	AfterEach(func() {
		testApp.Close()
	})

	It("POST /users/new", func() {
		By("should not be able to create any existing users")
		for _, user := range model.TestUsers {
			testApp.RawAPI.POST("/users/new").WithJSON(map[string]string{"username": user.Username, "password": user.Password}).
				Expect().
				Status(http.StatusBadRequest).Body().Contains("already exists")
		}

		By("should be able to create users")
		testApp.RawAPI.POST("/users/new").WithJSON(map[string]string{"username": "u1", "password": "p1"}).
			Expect().
			Status(http.StatusOK).JSON().Object().ContainsKey("username").ContainsKey("role").
			ContainsKey("id").NotContainsKey("password")
		testApp.RawAPI.POST("/users/new").WithJSON(map[string]string{"username": "u2", "password": "p2"}).
			Expect().
			Status(http.StatusOK).JSON().Object().ContainsKey("username").ContainsKey("role").
			ContainsKey("id").NotContainsKey("password")
	})
})
