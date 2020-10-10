package store

type Todos interface {
	//Get(id int) (model.Todo, error)
	//Create(note model.Todo) (string, error)
	//Update(note model.Todo) error
}

type implTodos struct {
}

var _ Todos = (*implTodos)(nil)

func NewTodos(ctx *StoreContext) Todos {
	return &implTodos{}
}
