package handlers_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/rs/xid"
	"github.com/weidonglian/golang-notes-app/model"
	"github.com/weidonglian/golang-notes-app/test"
	"net/http"
)

var _ = Describe("Todos", func() {

	var (
		testApp       test.TestApp
		testUserNotes []model.NoteWithTodos
		devUserNotes  []model.NoteWithTodos
	)

	BeforeEach(func() {
		testApp = test.NewTestAppAndServe()
		testUserNotes = test.NewTestUserNotesData(&testApp)
		devUserNotes = test.NewDevUserNotesData(&testApp)
	})

	AfterEach(func() {
		testApp.Close()
	})

	It("GET /todos/{id}", func() {
		Context("we should be able to get the test users todo by id", func() {
			for i := range testUserNotes {
				for j := range testUserNotes[i].Todos {
					testApp.API.GET("/todos/{id}", testUserNotes[i].Todos[j].ID).
						Expect().
						Status(http.StatusOK).JSON().Object().Equal(testUserNotes[i].Todos[j])
				}
			}
		})

		Context("Non-existent ids should not work", func() {
			for _, todoID := range []int{10000, 20000, 50, 19000000} {
				testApp.API.GET("/todos/{id}", todoID).
					Expect().
					Status(http.StatusUnprocessableEntity).Body().Contains("unable to find todo")
			}
		})

		Context("Bad {id} format should not crash the app", func() {
			for _, todoID := range []interface{}{"not_a_id", "@@@@@", `////`, `___????^%`} {
				testApp.API.GET("/todos/{id}", todoID).
					Expect().
					Status(http.StatusBadRequest)
			}
		})

		Context("we are not allowed to fetch another user's resources", func() {
			// testApp.API is authenticated for 'test' user should not get notes of 'dev' user even the note id is valid
			for i := range devUserNotes {
				for j := range devUserNotes[i].Todos {
					testApp.API.GET("/todos/{id}", devUserNotes[i].Todos[j].ID).
						Expect().
						Status(http.StatusUnprocessableEntity).Body().Contains("unable to find todo for current user")
				}
			}
		})
	})

	It("POST /todos", func() {

		todoNames := []string{"pn_todo_1", "pn_todo_2", "pn_todo_3", "pn_todo_1"}

		Context("we should be able to create with any non-empty string", func() {
			for i := range testUserNotes {
				for _, todoName := range todoNames {
					fetchedObject := testApp.API.POST("/todos").WithJSON(map[string]interface{}{
						"name":   todoName,
						"done":   true,
						"noteId": testUserNotes[i].ID,
					}).Expect().
						Status(http.StatusOK).JSON().Object()
					fetchedObject.Keys().Contains("id", "name", "done", "noteId")
					fetchedObject.Values().Contains(todoName, testUserNotes[i].ID, true)
				}
			}
		})

		Context("Invalid noteId should not be able to create any todo", func() {
			for _, noteID := range []int{100, 200, 500, 999999} {
				for _, todoName := range todoNames {
					testApp.API.POST("/todos").WithJSON(map[string]interface{}{
						"name":   todoName,
						"todo":   true,
						"noteId": noteID,
					}).Expect().
						Status(http.StatusUnprocessableEntity).Body().Contains("provide note with id does not exist for current user")
				}
			}
		})
	})

	It("PUT /todos/{id}", func() {
		Context("we should be able to update todo", func() {
			for i := range testUserNotes {
				for j := range testUserNotes[i].Todos {
					randomName := xid.New().String()
					obj := testApp.API.PUT("/todos/{id}", testUserNotes[i].Todos[j].ID).WithJSON(map[string]interface{}{
						"name": randomName,
						"done": true,
					}).Expect().
						Status(http.StatusOK).JSON().Object()
					obj.Keys().Contains("id", "name", "done", "noteId")
					obj.Values().Contains(testUserNotes[i].Todos[j].ID, randomName, true, testUserNotes[i].ID)
				}
			}
		})
	})

	It("PUT /todos/{id}/toggle", func() {
		Context("we should be able to toggle done", func() {
			for i := range testUserNotes {
				for j := range testUserNotes[i].Todos {
					testApp.API.PUT("/todos/{id}/toggle", testUserNotes[i].Todos[j].ID).Expect().
						Status(http.StatusOK).JSON().Object().Value("done").NotEqual(testUserNotes[i].Todos[j].Done)
				}
			}
		})

		Context("we should not be able to toggle another user's resources even valid ids", func() {
			for i := range devUserNotes {
				for j := range devUserNotes[i].Todos {
					testApp.API.PUT("/todos/{id}/toggle", devUserNotes[i].Todos[j].ID).Expect().
						Status(http.StatusUnprocessableEntity)
				}
			}
		})
	})

	It("DELETE /todos/{id}", func() {
		Context("we should be able to delete the test user's todo by id", func() {
			for i := range testUserNotes {
				for j := range testUserNotes[i].Todos {
					testApp.API.DELETE("/todos/{id}", testUserNotes[i].Todos[j].ID).Expect().
						Status(http.StatusOK)
				}
			}
		})

		Context("we should not be able to delete valid id of another user", func() {
			for i := range devUserNotes {
				for j := range devUserNotes[i].Todos {
					testApp.API.DELETE("/todos/{id}", devUserNotes[i].Todos[j].ID).Expect().
						Status(http.StatusUnprocessableEntity)
				}
			}
		})
	})
})
