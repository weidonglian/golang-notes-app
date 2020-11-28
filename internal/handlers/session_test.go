package handlers_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/weidonglian/notes-app/internal/model"
	. "github.com/weidonglian/notes-app/internal/test"
	"net/http"
)

var _ = Describe("Session", func() {
	Describe("POST /session", func() {
		var testApp MockApp

		BeforeEach(func() {
			testApp = NewMockApp()
		})

		AfterEach(func() {
			testApp.Close()
		})

		It("Test users should be able to login", func() {
			for _, user := range model.TestUsers {
				testApp.RawAPI.POST("/session").
					WithJSON(map[string]string{"username": user.Username, "password": user.Password}).
					Expect().
					Status(http.StatusOK).JSON().Object().ContainsKey("token")
			}
		})

		It("Nonexistent user should not be able to login", func() {
			testApp.RawAPI.POST("/session").
				WithJSON(map[string]string{"username": "xxx", "password": "yyy"}).
				Expect().
				Status(http.StatusUnauthorized).JSON().Object().NotContainsKey("token")
		})

		It("Without username and password should be bad", func() {
			By("only username should not be OK")
			testApp.RawAPI.POST("/session").
				WithJSON(map[string]string{"username": "xxx"}).
				Expect().
				Status(http.StatusBadRequest)

			By("only password should not be OK")
			testApp.RawAPI.POST("/session").
				WithJSON(map[string]string{"password": "xxx"}).
				Expect().
				Status(http.StatusBadRequest)

			By("no username and password should not be OK")
			testApp.RawAPI.POST("/session").
				Expect().
				Status(http.StatusBadRequest)
		})
	})
})
