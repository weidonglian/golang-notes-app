package store

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/weidonglian/golang-notes-app/model"
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

func (i NotesStore) Create(note model.Note) (int, error) {
	var id int
	stmt, err := i.db.PrepareNamed(`
		INSERT INTO notes (note_name, user_id)
		VALUES(:note_name, :user_id)
		RETURNING note_id
	`)
	if err != nil {
		return id, err
	}
	err = stmt.Get(&id, note)
	return id, err
}

func (i NotesStore) Update(id int, name string) (*model.Note, error) {
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

func (i NotesStore) Delete(id int) error {
	_, err := i.db.Exec("DELETE FROM notes WHERE note_id = $1", id)
	return err
}

func (i NotesStore) DeleteAll() error {
	_, err := i.db.Exec("TRUNCATE TABLE notes CASCADE")
	return err
}

func (i NotesStore) FindByUserID(userId int) []model.Note {
	var notes []model.Note
	err := i.db.Select(&notes, "SELECT * FROM notes WHERE user_id = $1", userId)
	if err != nil {
		return nil
	}
	return notes
}

func (i NotesStore) FindByID(id int) *model.Note {
	note := model.Note{}
	err := i.db.Get(&note, "SELECT * FROM notes WHERE note_id = $1", id)
	if err != nil {
		return nil
	}
	return &note
}

func (i NotesStore) FindByName(name string) []model.Note {
	var notes []model.Note
	err := i.db.Select(&notes, "SELECT * FROM notes WHERE note_name = $1", name)
	if err != nil {
		return nil
	}
	return notes
}
