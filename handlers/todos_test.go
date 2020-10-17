package handlers_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/weidonglian/golang-notes-app/handlers/test"
	"github.com/weidonglian/golang-notes-app/model"
	"github.com/weidonglian/golang-notes-app/store"
	"net/http"
)

var _ = Describe("Todos", func() {

	var (
		testApp       test.HandlerTestApp
		usersStore    store.UsersStore
		notesStore    store.NotesStore
		todosStore    store.TodosStore
		testUserId    int
		testUserNotes []model.NoteWithTodos

		devUserId    int
		devUserNotes []model.NoteWithTodos
	)

	BeforeEach(func() {
		testApp = test.NewTestAppAndServe()
		usersStore = testApp.App.GetStore().Users
		notesStore = testApp.App.GetStore().Notes
		todosStore = testApp.App.GetStore().Todos

		// test user test data
		testUserId = usersStore.FindByName("test").ID
		testUserNotes = []model.NoteWithTodos{
			{
				Note: &model.Note{
					Name: "n1",
				},
				Todos: []model.Todo{
					{
						Name: "todo_1",
						Done: true,
					},
					{
						Name: "todo_2",
						Done: false,
					},
					{
						Name: "todo3",
						Done: false,
					},
				},
			},
			{
				Note: &model.Note{
					Name: "n2",
				},
				Todos: []model.Todo{
					{
						Name: "todo_1",
						Done: true,
					},
					{
						Name: "todo_2",
						Done: false,
					},
					{
						Name: "todo3",
						Done: false,
					},
				},
			},
			{
				Note: &model.Note{
					Name: "n3",
				},
				Todos: []model.Todo{
					{
						Name: "todo_1",
						Done: true,
					},
					{
						Name: "todo_2",
						Done: false,
					},
					{
						Name: "todo3",
						Done: false,
					},
				},
			},
		}

		test.FillDataToStore(testApp.App.GetStore(), "test", testUserNotes)

		// dev user test data
		devUserId = usersStore.FindByName("dev").ID
		testUserNotes = []model.NoteWithTodos{
			{
				Note: &model.Note{
					Name: "n4",
				},
				Todos: []model.Todo{
					{
						Name: "todo_1",
						Done: true,
					},
					{
						Name: "todo_2",
						Done: false,
					},
					{
						Name: "todo3",
						Done: false,
					},
				},
			},
			{
				Note: &model.Note{
					Name: "n5",
				},
				Todos: []model.Todo{
					{
						Name: "todo_1",
						Done: true,
					},
					{
						Name: "todo_2",
						Done: false,
					},
					{
						Name: "todo3",
						Done: false,
					},
				},
			},
			{
				Note: &model.Note{
					Name: "n6",
				},
				Todos: []model.Todo{
					{
						Name: "todo_1",
						Done: true,
					},
					{
						Name: "todo_2",
						Done: false,
					},
					{
						Name: "todo3",
						Done: false,
					},
				},
			},
		}

		test.FillDataToStore(testApp.App.GetStore(), "dev", devUserNotes)
	})

	AfterEach(func() {
		testApp.Close()
	})

	It("GET /todos", func() {
		Context("should fetch notes of test user, should not include notes of dev user", func() {
			fetchedNotes := testApp.API.GET("/notes").
				Expect().
				Status(http.StatusOK).JSON().Array()
			fetchedNotes.Length().Equal(3)
			for i := range testUserNotes {
				fetchedNotes.Element(i).Object().Keys().Contains("id", "name").NotContains("userId")
				fetchedNotes.Element(i).Object().Values().Contains(testUserNotes[i].ID, testUserNotes[i].Name)
				fetchedNotes.Element(i).Object().Values().NotContains(devUserNotes[i].ID, devUserNotes[i].Name)
			}
		})

	})

	It("GET /todos/{id}", func() {
		Context("we should be able to get the test users notes by id", func() {
			for i := range testUserNotes {
				fetchedNote := testApp.API.GET("/notes/{id}", testUserNotes[i].ID).
					Expect().
					Status(http.StatusOK).JSON().Object()
				fetchedNote.Value("id").Equal(testUserNotes[i].ID)
				fetchedNote.Value("name").Equal(testUserNotes[i].Name)
				fetchedNote.NotContainsKey("userId")
			}
		})

		Context("we are not allowed to fetch another user's resources", func() {
			// testApp.API is authenticated for 'test' user should not get notes of 'dev' user even the note id is valid
			for i := range devUserNotes {
				testApp.API.GET("/notes/{id}", devUserNotes[i].ID).
					Expect().
					Status(http.StatusUnprocessableEntity).Body().Contains("unable to find note")
			}
		})
	})

	It("POST /todos", func() {
		Context("we should be able to create name with any non-empty string", func() {
			noteNames := []string{"pn1", "pn2", "pn3"}
			for _, noteName := range noteNames {
				testApp.API.POST("/notes").WithJSON(map[string]string{"name": noteName}).
					Expect().
					Status(http.StatusOK).JSON().Object().
					ContainsKey("id").NotContainsKey("userId").
					Value("name").Equal(noteName)
			}
		})

		Context("empty note name is not allowed to create", func() {
			testApp.API.POST("/notes").WithJSON(map[string]string{"name": ""}).
				Expect().
				Status(http.StatusBadRequest).JSON().Object().Value("error").String().Contains("missing required fields")
			testApp.API.POST("/notes").
				Expect().
				Status(http.StatusBadRequest).Body().Contains("unable to automatically decode the request content type")
		})
	})

	It("PUT /todos/{id}", func() {
		Context("we should be able to update note name", func() {
			newNames := []string{"xy1", "xy2", "xy3"}
			if len(newNames) != len(testUserNotes) {
				panic("new notes name count should match that of testUserNotes")
			}
			for i := range newNames {
				testApp.API.PUT("/notes/{id}", testUserNotes[i].ID).WithJSON(map[string]string{"name": newNames[i]}).
					Expect().
					Status(http.StatusOK).JSON().Object().
					ContainsKey("id").NotContainsKey("userId").
					Value("name").Equal(newNames[i])
			}
		})
	})

	It("DELETE /todos/{id}", func() {
		Context("we should be able to delete the test user's notes by id", func() {
			for i := range testUserNotes {
				testApp.API.DELETE("/notes/{id}", testUserNotes[i].ID).
					Expect().
					Status(http.StatusOK).Body().Empty()
			}
		})

		Context("we should not be able to delete valid id of another user", func() {
			for i := range devUserNotes {
				testApp.API.DELETE("/notes/{id}", devUserNotes[i].ID).
					Expect().
					Status(http.StatusUnprocessableEntity).Body().Contains("unable to find note")
			}
		})
	})
})
