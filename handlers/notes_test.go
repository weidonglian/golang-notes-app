package handlers_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	"github.com/weidonglian/golang-notes-app/model"
	"github.com/weidonglian/golang-notes-app/store"
	"net/http"
)

var _ = Describe("Notes", func() {

	var (
		testApp    HandlerTestApp
		usersStore store.UsersStore
		notesStore store.NotesStore
		testUserId int
		notes      []model.Note
	)

	BeforeEach(func() {
		testApp = NewTestAppAndServe()
		usersStore = testApp.App.GetStore().Users
		notesStore = testApp.App.GetStore().Notes
		testUserId = usersStore.FindByName("test").ID
		notes = []model.Note{
			{
				Name:   "n1",
				UserID: testUserId,
			},
			{
				Name:   "n2",
				UserID: testUserId,
			},
			{
				Name:   "n3",
				UserID: testUserId,
			},
		}
		for i := range notes {
			if id, err := notesStore.Create(notes[i]); err != nil {
				panic(fmt.Sprintf("failed to create note: %v", notes[i]))
			} else {
				notes[i].ID = id
			}
		}
	})

	AfterEach(func() {
		testApp.Close()
	})

	It("GET /notes", func() {
		fetchedNotes := testApp.API.GET("/notes").
			Expect().
			Status(http.StatusOK).JSON().Array()
		fetchedNotes.Length().Equal(3)
		for i := range notes {
			fetchedNotes.Element(i).Object().Keys().Contains("id", "name").NotContains("userId")
			fetchedNotes.Element(i).Object().Values().Contains(notes[i].ID, notes[i].Name)
		}
	})

	It("GET /notes/{id}", func() {
		for i := range notes {
			fetchedNote := testApp.API.GET("/notes/{id}", notes[i].ID).
				Expect().
				Status(http.StatusOK).JSON().Object()
			fetchedNote.Value("id").Equal(notes[i].ID)
			fetchedNote.Value("name").Equal(notes[i].Name)
			fetchedNote.NotContainsKey("userId")
		}
	})
})
