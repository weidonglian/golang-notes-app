package handlers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rs/xid"
	"github.com/weidonglian/notes-app/internal/store"
	"github.com/weidonglian/notes-app/pkg/model"
	"github.com/weidonglian/notes-app/pkg/test"
	"net/http"
)

var _ = Describe("Notes", func() {

	var (
		testApp       test.TestApp
		notesStore    store.NotesStore
		testUserId    int
		devUserId     int
		testUserNotes []model.NoteWithTodos
		devUserNotes  []model.NoteWithTodos
	)

	BeforeEach(func() {
		testApp = test.NewTestAppAndServe()
		notesStore = testApp.App.GetStore().Notes
		testUserId = testApp.App.GetStore().Users.FindByName("test").ID
		devUserId = testApp.App.GetStore().Users.FindByName("dev").ID
		testUserNotes = test.NewTestUserNotesData(&testApp)
		devUserNotes = test.NewDevUserNotesData(&testApp)
	})

	AfterEach(func() {
		testApp.Close()
	})

	It("GET /notes", func() {
		Context("should fetch notes of test user, should not include notes of dev user", func() {
			fetchedNotes := testApp.API.GET("/notes").
				Expect().
				Status(http.StatusOK).JSON().Array()
			fetchedNotes.Length().Equal(len(testUserNotes))
			for i := range testUserNotes {
				fetchedNote := fetchedNotes.Element(i).Object()
				fetchedNote.Keys().Contains("id", "name", "todos").NotContains("userId")
				fetchedNote.Values().Contains(testUserNotes[i].ID, testUserNotes[i].Name)
				for _, devNote := range devUserNotes {
					fetchedNote.Values().NotContains(devNote.ID, devNote.Name)
				}
				fetchedTodos := fetchedNote.Value("todos").Array()
				fetchedTodos.Length().Equal(len(testUserNotes[i].Todos))
				for j := range testUserNotes[i].Todos {
					fetchedTodos.Element(j).Object().Keys().Contains("id", "name", "done", "noteId")
					fetchedTodos.Element(j).Object().Equal(testUserNotes[i].Todos[j])
				}
			}
		})

	})

	It("GET /notes/{id}", func() {
		Context("we should be able to get the test users notes by id", func() {
			for i := range testUserNotes {
				fetchedNote := testApp.API.GET("/notes/{id}", testUserNotes[i].ID).
					Expect().
					Status(http.StatusOK).JSON().Object()
				fetchedNote.Value("id").Equal(testUserNotes[i].ID)
				fetchedNote.Value("name").Equal(testUserNotes[i].Name)
				fetchedNote.Value("todos").Equal(testUserNotes[i].Todos)
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

	It("POST /notes", func() {
		Context("we should be able to create name with any non-empty string", func() {
			noteNames := []string{"pn1", "pn2", "pn3"}
			for _, noteName := range noteNames {
				newNote := testApp.API.POST("/notes").WithJSON(map[string]string{"name": noteName}).
					Expect().
					Status(http.StatusOK).JSON().Object()
				newNote.Keys().Contains("id", "name", "todos").NotContains("userId")
				newNote.Value("name").Equal(noteName)
				newNote.Value("todos").Array().Length().Equal(0)
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

	It("PUT /notes/{id}", func() {
		Context("we should be able to update note name", func() {
			for _, note := range testUserNotes {
				randomName := xid.New().String()
				updatedNote := testApp.API.PUT("/notes/{id}", note.ID).WithJSON(map[string]string{"name": randomName}).
					Expect().
					Status(http.StatusOK).JSON().Object()
				updatedNote.Keys().Contains("id", "name", "todos").NotContains("userId")
				updatedNote.Value("name").Equal(randomName)
				updatedNote.Value("todos").Equal(note.Todos)
			}
		})
	})

	It("DELETE /notes/{id}", func() {
		Context("we should be able to delete the test user's notes by id", func() {
			for _, note := range testUserNotes {
				testApp.API.DELETE("/notes/{id}", note.ID).
					Expect().
					Status(http.StatusOK)
			}
		})

		Context("we should not be able to delete valid id of another user", func() {
			for _, note := range devUserNotes {
				testApp.API.DELETE("/notes/{id}", note.ID).
					Expect().
					Status(http.StatusUnprocessableEntity).Body().Contains("unable to find note")
			}
		})
	})

	It("DELETE /notes", func() {
		Context("we should be able to delete the test user's notes", func() {
			testApp.API.DELETE("/notes").
				Expect().
				Status(http.StatusOK).Body().Empty()
			Expect(len(notesStore.FindByUserID(testUserId))).To(BeZero())
			Expect(len(notesStore.FindByUserID(devUserId))).NotTo(BeZero())
		})
	})
})
