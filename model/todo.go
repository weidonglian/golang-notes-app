package model

import "time"

type Todo struct {
	ID        int       `db:"todo_id" json:"id"`
	Name      string    `db:"todo_name" json:"name"`
	Done      bool      `db:"todo_done" json:"done"`
	NoteID    int       `db:"note_id" json:"noteId"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

func NewTodo() Todo {
	return Todo{}
}
