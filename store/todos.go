package store

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/weidonglian/golang-notes-app/model"
)

type TodosStore struct {
	db     *sqlx.DB
	logger *logrus.Logger
}

func NewTodosStore(ctx *StoreContext) TodosStore {
	return TodosStore{
		db:     ctx.Session.GetDB(),
		logger: ctx.Session.Logger,
	}
}

func (i TodosStore) Create(todo model.Todo) (int, error) {
	var id int
	stmt, err := i.db.PrepareNamed(`
		INSERT INTO todos (todo_name, todo_done, note_id)
		VALUES(:todo_name, :todo_done, :note_id)
		RETURNING todo_id
	`)
	if err != nil {
		return id, err
	}
	err = stmt.Get(&id, todo)
	return id, err
}

func (i TodosStore) UpdateName(id int, name string) (*model.Todo, error) {
	stmt, err := i.db.Preparex(`
		UPDATE todos
		SET todo_name = $1
		WHERE todo_id = $2
		RETURNING *
	`)
	if err != nil {
		return nil, err
	}
	todo := model.Todo{}
	err = stmt.Get(&todo, name, id)
	return &todo, err
}

// Tries to delete a user by id, and returns the number of records deleted;
func (i TodosStore) Delete(id int) error {
	_, err := i.db.Exec("DELETE FROM todos WHERE todo_id = $1", id)
	return err
}

// Removes all records from the table;
func (i TodosStore) DeleteAll() error {
	_, err := i.db.Exec("TRUNCATE TABLE todos CASCADE")
	return err
}

// Tries to find from id;
func (i TodosStore) FindByID(id int) *model.Todo {
	todo := model.Todo{}
	if err := i.db.Get(&todo, "SELECT * FROM todos WHERE todo_id = $1", id); err != nil {
		return nil
	}
	return &todo
}

// Tries to find from name;
func (i TodosStore) FindByName(name string) []model.Todo {
	var todos []model.Todo
	if err := i.db.Select(&todos, "SELECT * FROM todos WHERE todo_name = $1", name); err != nil {
		return []model.Todo{}
	}
	return todos
}

// Tries to find from note_id;
func (i TodosStore) FindByNoteID(noteId int) []model.Todo {
	var todos []model.Todo
	if err := i.db.Select(&todos, "SELECT * FROM todos WHERE note_id = $1", noteId); err != nil {
		return []model.Todo{}
	}
	return todos
}
