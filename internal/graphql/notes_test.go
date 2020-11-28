package graphql_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/rs/xid"
	"github.com/weidonglian/notes-app/internal/graphql/gmodel"
	"github.com/weidonglian/notes-app/internal/model"
	"github.com/weidonglian/notes-app/internal/test"
)

var _ = Describe("GraphQL Notes", func() {

	var (
		testApp       test.MockApp
		testUserNotes []model.NoteWithTodos
		devUserNotes  []model.NoteWithTodos
	)

	BeforeEach(func() {
		testApp = test.NewMockApp()
		testUserNotes = test.NewTestUserNotesData(&testApp)
		devUserNotes = test.NewDevUserNotesData(&testApp)
	})

	AfterEach(func() {
		testApp.Close()
	})

	It("query notes", func() {
		Context("should fetch notes of test user, should not include notes of dev user", func() {
			fetchedNotes := testApp.GraphqlMustData(test.QueryNotes).ContainsKey("notes").Value("notes").Array()

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

	It("query note by id", func() {
		Context("we should be able to get the test users notes by id", func() {
			for i := range testUserNotes {
				fetchedNote := testApp.GraphqlMustData(test.QueryNote, map[string]interface{}{
					"id": testUserNotes[i].ID,
				}).ContainsKey("note").Value("note").Object()
				fetchedNote.Value("id").Equal(testUserNotes[i].ID)
				fetchedNote.Value("name").Equal(testUserNotes[i].Name)
				fetchedNote.Value("todos").Equal(testUserNotes[i].Todos)
				fetchedNote.NotContainsKey("userId")
			}
		})

		Context("we are not allowed to fetch another user's resources", func() {
			// testApp.API is authenticated for 'test' user should not get notes of 'dev' user even the note id is valid
			for i := range devUserNotes {
				testApp.GraphqlMustError(test.QueryNote, map[string]interface{}{
					"id": devUserNotes[i].ID,
				}).NotEmpty().Element(0).Object().Value("message").String().
					Contains("failed to find a note with id")
			}
		})
	})

	It("mutation addNote", func() {
		Context("we should be able to create name with any non-empty string", func() {
			noteNames := []string{"pn1", "pn2", "pn3"}
			for _, noteName := range noteNames {
				newNote := testApp.GraphqlMustData(test.MutationAddNote, test.GraphqlWithInput(gmodel.AddNoteInput{
					Name: noteName,
				})).ContainsKey("addNote").Value("addNote").Object()

				newNote.Keys().Contains("id", "name", "todos").NotContains("userId")
				newNote.Value("name").Equal(noteName)
			}
		})

		Context("empty note name is not allowed to create", func() {
			testApp.GraphqlMustError(test.MutationAddNote, test.GraphqlWithInput(gmodel.AddNoteInput{Name: ""})).
				NotEmpty().Element(0).Object().Value("message").String().
				Contains("'name' field can not be empty")
			testApp.GraphqlMustError(test.MutationAddNote).
				NotEmpty().Element(0).Object().Value("message").String().
				Contains("must be defined")
		})
	})

	It("mutation updateNote", func() {
		Context("we should be able to update note name", func() {
			for _, note := range testUserNotes {
				randomName := xid.New().String()

				updatedNote := testApp.GraphqlMustData(test.MutationUpdateNote, test.GraphqlWithInput(gmodel.UpdateNoteInput{
					ID:   note.ID,
					Name: randomName,
				})).ContainsKey("updateNote").Value("updateNote").Object()

				updatedNote.Keys().Contains("id", "name", "todos").NotContains("userId")
				updatedNote.Value("name").Equal(randomName)
				updatedNote.Value("id").Equal(note.ID)
			}
		})

		Context("empty name is not allowed to update", func() {
			testApp.GraphqlMustError(test.MutationUpdateNote, test.GraphqlWithInput(gmodel.UpdateNoteInput{
				ID:   testUserNotes[0].ID,
				Name: "",
			})).NotEmpty().Element(0).Object().Value("message").String().Contains("'name' field can not be empty")
		})
	})

	It("mutation deleteNote", func() {
		Context("we should be able to delete the test user's notes by id", func() {
			for _, note := range testUserNotes {
				testApp.GraphqlMustData(test.MutationDeleteNote, test.GraphqlWithInput(gmodel.DeleteNoteInput{
					ID: note.ID,
				})).ContainsKey("deleteNote").Value("deleteNote").
					Object().Keys().ContainsOnly("id")
			}
		})

		Context("we should not be able to delete valid id of another user", func() {
			for _, note := range devUserNotes {
				testApp.GraphqlMustError(test.MutationDeleteNote, test.GraphqlWithInput(gmodel.DeleteNoteInput{
					ID: note.ID,
				})).NotEmpty().Element(0).Object().Value("message").String().Contains("unprocessable entity with")
			}
		})
	})
})
