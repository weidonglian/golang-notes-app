package graph_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/rs/xid"
	"github.com/weidonglian/golang-notes-app/graph/gmodel"
	"github.com/weidonglian/golang-notes-app/graph/util"
	"github.com/weidonglian/golang-notes-app/model"
	"github.com/weidonglian/golang-notes-app/test"
)

var _ = Describe("Graph Todos", func() {

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

	It("query todo by noteId", func() {
		Context("fetch the todos of given noteId", func() {
			for _, note := range testUserNotes {
				testApp.GraphMustData(test.QueryTodos, map[string]interface{}{"noteId": note.ID}).
					ContainsKey("todos").Value("todos").Array().Length().Equal(len(note.Todos))
			}
		})

		Context("we are not allowed to fetch another user's resources", func() {
			// testApp.API is authenticated for 'test' user should not get notes of 'dev' user even the note id is valid
			for _, note := range devUserNotes {
				testApp.GraphMustError(test.QueryTodos, map[string]interface{}{"noteId": note.ID}).
					NotEmpty().Element(0).Object().Value("message").String().Contains(util.ErrorUnprocessableEntity.Error())
			}
		})
	})

	It("mutation addTodo", func() {

		todoNames := []string{"pn_todo_1", "pn_todo_2", "pn_todo_3", "pn_todo_1"}

		Context("we should be able to create with any non-empty string", func() {
			done := true
			for i := range testUserNotes {
				for _, todoName := range todoNames {
					fetchedObject := testApp.GraphMustData(test.MutationAddTodo, test.GraphWithInput(gmodel.AddTodoInput{
						Name:   todoName,
						Done:   &done,
						NoteID: testUserNotes[i].ID,
					})).ContainsKey("addTodo").Value("addTodo").Object()
					fetchedObject.Keys().Contains("id", "name", "done", "noteId")
					fetchedObject.Values().Contains(todoName, testUserNotes[i].ID, done).NotContains(0)
				}
			}
		})

		Context("Invalid noteId should not be able to create any todo", func() {
			for _, noteID := range []int{100, 200, 500, 999999} {
				for _, todoName := range todoNames {
					testApp.GraphMustError(test.MutationAddTodo, test.GraphWithInput(gmodel.AddTodoInput{
						Name:   todoName,
						NoteID: noteID,
					})).NotEmpty().Element(0).Object().Value("message").String().Contains(util.ErrorUnprocessableEntity.Error())
				}
			}
		})
	})

	It("mutation updateTodo", func() {
		Context("we should be able to update todo", func() {
			done := true
			for i := range testUserNotes {
				for j := range testUserNotes[i].Todos {
					randomName := xid.New().String()
					obj := testApp.GraphMustData(test.MutationUpdateTodo, test.GraphWithInput(gmodel.UpdateTodoInput{
						ID:     testUserNotes[i].Todos[j].ID,
						Name:   randomName,
						Done:   &done,
						NoteID: testUserNotes[i].ID,
					})).ContainsKey("updateTodo").Value("updateTodo").Object()

					obj.Keys().Contains("id", "name", "done", "noteId")
					obj.Values().Contains(testUserNotes[i].Todos[j].ID, randomName, done, testUserNotes[i].ID)
				}
			}
		})
	})

	It("mutation toggleTodo", func() {
		Context("we should be able to toggle done", func() {
			for i := range testUserNotes {
				for j := range testUserNotes[i].Todos {
					obj := testApp.GraphMustData(test.MutationToggleTodo, test.GraphWithInput(gmodel.ToggleTodoInput{
						ID:     testUserNotes[i].Todos[j].ID,
						NoteID: testUserNotes[i].ID,
					})).ContainsKey("toggleTodo").Value("toggleTodo").Object()

					obj.Keys().ContainsOnly("id", "done", "noteId", "name")
					obj.Values().Contains(testUserNotes[i].Todos[j].ID, !testUserNotes[i].Todos[j].Done, testUserNotes[i].ID)
				}
			}
		})

		Context("we should not be able to toggle another user's resources even valid ids", func() {
			for i := range devUserNotes {
				for j := range devUserNotes[i].Todos {
					testApp.GraphMustError(test.MutationToggleTodo, test.GraphWithInput(gmodel.ToggleTodoInput{
						ID:     devUserNotes[i].Todos[j].ID,
						NoteID: devUserNotes[i].ID,
					})).NotEmpty().Element(0).Object().Value("message").String().Contains(util.ErrorUnprocessableEntity.Error())
				}
			}
		})
	})

	It("mutation deleteTodo", func() {
		Context("we should be able to delete the test user's todo by id", func() {
			for i := range testUserNotes {
				for j := range testUserNotes[i].Todos {
					obj := testApp.GraphMustData(test.MutationDeleteTodo, test.GraphWithInput(gmodel.DeleteTodoInput{
						ID:     testUserNotes[i].Todos[j].ID,
						NoteID: testUserNotes[i].ID,
					})).ContainsKey("deleteTodo").Value("deleteTodo").Object()
					obj.Keys().ContainsOnly("id", "noteId")
					obj.Values().Contains(testUserNotes[i].Todos[j].ID, testUserNotes[i].ID)

					testApp.GraphMustError(test.MutationDeleteTodo, test.GraphWithInput(gmodel.DeleteTodoInput{
						ID:     testUserNotes[i].Todos[j].ID,
						NoteID: testUserNotes[i].ID,
					})).NotEmpty().Element(0).Object().Value("message").String().Contains(util.ErrorUnprocessableEntity.Error())
				}
			}
		})

		Context("we should not be able to delete valid id of another user", func() {
			for i := range devUserNotes {
				for j := range devUserNotes[i].Todos {
					testApp.GraphMustError(test.MutationDeleteTodo, test.GraphWithInput(gmodel.DeleteTodoInput{
						ID:     devUserNotes[i].Todos[j].ID,
						NoteID: devUserNotes[i].ID,
					})).NotEmpty().Element(0).Object().Value("message").String().Contains(util.ErrorUnprocessableEntity.Error())
				}
			}
		})
	})
})
