package store

import (
	"github.com/weidonglian/golang-notes-app/model"
)

type NotesStore interface {
	Get(id int) (*model.Note, error)
	Create(note model.Note) (string, error)
	Update(note model.Note) error
}

type implNotesStore struct {
}

func (impl *implNotesStore) Get(id int) (*model.Note, error) {
	return nil, nil
}

func (impl *implNotesStore) Create(note model.Note) (string, error) {
	return "", nil
}

func (impl *implNotesStore) Update(note model.Note) error {
	return nil
}

var _ NotesStore = (*implNotesStore)(nil)

func NewNotesStore(ctx *StoreContext) NotesStore {
	return &implNotesStore{}
}
