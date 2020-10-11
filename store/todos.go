package store

type TodosStore interface {
	//Get(id int) (model.Todo, error)
	//Create(note model.Todo) (string, error)
	//Update(note model.Todo) error
}

type implTodosStore struct {
}

var _ TodosStore = (*implTodosStore)(nil)

func NewTodosStore(ctx *StoreContext) TodosStore {
	return &implTodosStore{}
}
