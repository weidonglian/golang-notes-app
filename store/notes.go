package store

import (
	"github.com/weidonglian/golang-notes-app/model"
)

type Notes interface {
	Get(id int) (model.Note, error)
	Create(note model.Note) (string, error)
	Update(note model.Note) error
}
