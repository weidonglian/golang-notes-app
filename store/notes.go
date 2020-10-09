package store

import (
	"github.com/weidonglian/golang-notes-app/model"
)

type Notes interface {
	Get(id int) (model.Note, error)
	Create(note model.Note) (string, error)
	Update(note model.Note) error
}

type Impl struct {
}

func (impl *Impl) Get(id int) (model.Note, error) {

}

func (impl *Impl) Create(note model.Note) (string, error) {

}

func (impl *Impl) Update(note model.Note) error {

}

var _ Notes = (*Impl)(nil)

func NewNotes(ctx *StoreContext) Notes {
	return &Impl{}
}
