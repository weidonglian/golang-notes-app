package graphql

import (
	"github.com/weidonglian/notes-app/internal/graphql/gmodel"
	"github.com/weidonglian/notes-app/internal/model"
)

func NewGTodo(todo *model.Todo) *gmodel.Todo {
	return &gmodel.Todo{
		ID:     todo.ID,
		Name:   todo.Name,
		Done:   todo.Done,
		NoteID: todo.NoteID,
	}
}

func NewGNote(note *model.Note, todos []model.Todo) *gmodel.Note {
	gtodos := make([]*gmodel.Todo, len(todos))
	for i := range todos {
		gtodos[i] = NewGTodo(&todos[i])
	}
	return &gmodel.Note{
		ID:    note.ID,
		Name:  note.Name,
		Todos: gtodos,
	}
}
