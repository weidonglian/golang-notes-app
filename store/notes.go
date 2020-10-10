package store

import (
	"github.com/weidonglian/golang-notes-app/model"
)

type Notes interface {
	Get(id int) (*model.Note, error)
	Create(note model.Note) (string, error)
	Update(note model.Note) error
}

type implNotes struct {
}

func (impl *implNotes) Get(id int) (*model.Note, error) {
	return nil, nil
}

func (impl *implNotes) Create(note model.Note) (string, error) {
	return "", nil
}

func (impl *implNotes) Update(note model.Note) error {
	return nil
}

var _ Notes = (*implNotes)(nil)

func NewNotes(ctx *StoreContext) Notes {
	return &implNotes{}
}
