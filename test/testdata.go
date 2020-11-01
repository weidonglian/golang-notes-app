package test

import "github.com/weidonglian/golang-notes-app/model"

func NewTestUserNotesData(testApp *HandlerTestApp) []model.NoteWithTodos {
	// test user test data
	testUserNotes := []model.NoteWithTodos{
		{
			Note: &model.Note{
				Name: "n0",
			},
			Todos: []model.Todo{}, // without any todos
		},
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

	FillDataToStore(testApp.App.GetStore(), "test", testUserNotes)

	return testUserNotes
}

func NewDevUserNotesData(testApp *HandlerTestApp) []model.NoteWithTodos {
	// dev user test data
	devUserNotes := []model.NoteWithTodos{
		{
			Note: &model.Note{
				Name: "n4",
			},
			Todos: []model.Todo{
				{
					Name: "todo_dev_1",
					Done: true,
				},
				{
					Name: "todo_dev_2",
					Done: false,
				},
				{
					Name: "todo_dev_3",
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
					Name: "todo_dev_1",
					Done: true,
				},
				{
					Name: "todo_dev_2",
					Done: false,
				},
				{
					Name: "todo_dev_3",
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
					Name: "todo_dev_1",
					Done: true,
				},
				{
					Name: "todo_dev_2",
					Done: false,
				},
				{
					Name: "todo_dev_3",
					Done: false,
				},
			},
		},
	}

	FillDataToStore(testApp.App.GetStore(), "dev", devUserNotes)

	return devUserNotes
}
