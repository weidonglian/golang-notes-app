package store

import "github.com/weidonglian/golang-notes-app/model"

type Todos interface {
	Get(id int) (model.Todo, error)
	Create(note model.Todo) (string, error)
	Update(note model.Todo) error
}

type Impl struct {
}

var _ Todos = (*Impl)(nil)

func NewTodos(ctx *StoreContext) Todos {
	return &Impl{}
}
