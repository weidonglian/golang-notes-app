package store

import (
	"github.com/weidonglian/golang-notes-app/model"
)

type NotesStore struct {
}

func (impl *NotesStore) Get(id int) (*model.Note, error) {
	return nil, nil
}

func (impl *NotesStore) Create(note model.Note) (string, error) {
	return "", nil
}

func (impl *NotesStore) Update(note model.Note) error {
	return nil
}

func NewNotesStore(ctx *StoreContext) NotesStore {
	return NotesStore{}
}
