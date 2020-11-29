package model

type Note struct {
	ID     int    `db:"note_id" json:"id"`
	Name   string `db:"note_name" json:"name"`
	UserID int    `db:"user_id" json:"userId"`
}

type NoteWithTodos struct {
	*Note
	Todos []Todo
}
