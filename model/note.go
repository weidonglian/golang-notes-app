package model

import (
	"time"
)

type Note struct {
	ID        int       `db:"note_id" json:"id"`
	Name      string    `db:"note_name" json:"name"`
	UserID    int       `db:"user_id" json:"userId"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

func NewNote() Note {
	return Note{}
}
