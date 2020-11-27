package store

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/weidonglian/notes-app/model"
)

type NotesStore struct {
	db     *sqlx.DB
	logger *logrus.Logger
}

func NewNotesStore(ctx *Context) NotesStore {
	return NotesStore{
		db:     ctx.Session.GetDB(),
		logger: ctx.Session.Logger,
	}
}

func (i NotesStore) Create(note model.Note) (*model.Note, error) {
	stmt, err := i.db.PrepareNamed(`
		INSERT INTO notes (note_name, user_id)
		VALUES(:note_name, :user_id)
		RETURNING *
	`)
	if err != nil {
		return nil, err
	}
	var retNote model.Note
	err = stmt.Get(&retNote, note)
	return &retNote, err
}

func (i NotesStore) Update(id int, name string, userId int) (*model.Note, error) {
	stmt, err := i.db.Preparex(`
		UPDATE notes
		SET note_name = $1
		WHERE note_id = $2
		RETURNING *
	`)
	if err != nil {
		return nil, err
	}
	note := model.Note{}
	err = stmt.Get(&note, name, id)
	return &note, err
}

func (i NotesStore) Delete(id int, userId int) (int, error) {
	stmt, err := i.db.Preparex("DELETE FROM notes WHERE note_id = $1 AND user_id = $2 RETURNING note_id")
	if err != nil {
		return 0, err
	}
	retID := 0
	err = stmt.Get(&retID, id, userId)
	return retID, err
}

func (i NotesStore) DropAll(userId int) error {
	_, err := i.db.Exec("TRUNCATE TABLE notes CASCADE")
	return err
}

func (i NotesStore) DeleteAll(userId int) error {
	_, err := i.db.Exec("DELETE FROM notes WHERE user_id = $1", userId)
	return err
}

func (i NotesStore) FindByUserID(userId int) []model.Note {
	notes := make([]model.Note, 0)
	err := i.db.Select(&notes, "SELECT * FROM notes WHERE user_id = $1", userId)
	if err != nil {
		return notes
	}
	return notes
}

func (i NotesStore) FindByID(id int, userId int) *model.Note {
	note := model.Note{}
	err := i.db.Get(&note, "SELECT * FROM notes WHERE note_id = $1 AND user_id = $2", id, userId)
	if err != nil {
		return nil
	}
	return &note
}

func (i NotesStore) FindByName(name string, userId int) []model.Note {
	notes := make([]model.Note, 0)
	err := i.db.Select(&notes, "SELECT * FROM notes WHERE note_name = $1 AND user_id = $2", name, userId)
	if err != nil {
		return notes
	}
	return notes
}
