package store

type Users interface {
	//Get(id int) (model.Todo, error)
	//Create(note model.Todo) (string, error)
	//Update(note model.Todo) error
}

type implUsers struct {
}

func NewUsers(ctx *StoreContext) Users {
	return &implUsers{}
}
