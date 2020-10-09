package store

import "github.com/weidonglian/golang-notes-app/model"

type Users interface {
	Get(id int) (model.Todo, error)
	Create(note model.Todo) (string, error)
	Update(note model.Todo) error
}

var _ Users = (*Impl)(nil)

func NewUsers(ctx *StoreContext) Users {
	return &Impl{}
}
