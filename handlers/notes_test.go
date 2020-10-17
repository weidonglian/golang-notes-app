package handlers_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/weidonglian/golang-notes-app/handlers/test"
	"github.com/weidonglian/golang-notes-app/model"
	"github.com/weidonglian/golang-notes-app/store"
	"net/http"
)

var _ = Describe("Notes", func() {

	var (
		testApp       HandlerTestApp
		usersStore    store.UsersStore
		notesStore    store.NotesStore
		testUserId    int
		testUserNotes []model.Note
		devUserId     int
		devUserNotes  []model.Note
	)

	BeforeEach(func() {
		testApp = NewTestAppAndServe()
		usersStore = testApp.App.GetStore().Users
		notesStore = testApp.App.GetStore().Notes

		testUserId = usersStore.FindByName("test").ID
		testUserNotes = []model.Note{
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
		for i := range testUserNotes {
			if createdNote, err := notesStore.Create(testUserNotes[i]); err != nil {
				panic(fmt.Sprintf("failed to create note: %v with error:%v", testUserNotes[i], err))
			} else {
				testUserNotes[i].ID = createdNote.ID
			}
		}

		devUserId = usersStore.FindByName("dev").ID
		devUserNotes = []model.Note{
			{
				Name:   "n4",
				UserID: devUserId,
			},
			{
				Name:   "n5",
				UserID: devUserId,
			},
			{
				Name:   "n6",
				UserID: devUserId,
			},
		}
		for i := range devUserNotes {
			if createdNote, err := notesStore.Create(devUserNotes[i]); err != nil {
				panic(fmt.Sprintf("failed to create note: %v with error:%v", devUserNotes[i], err))
			} else {
				devUserNotes[i].ID = createdNote.ID
			}
		}
	})

	AfterEach(func() {
		testApp.Close()
	})

	It("GET /notes", func() {
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

	It("GET /notes/{id}", func() {
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

	It("POST /notes", func() {
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

	It("PUT /notes/{id}", func() {
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

	It("DELETE /notes/{id}", func() {
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
