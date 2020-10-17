package payload

import (
	"github.com/weidonglian/golang-notes-app/handlers/util"
	"github.com/weidonglian/golang-notes-app/model"
	"net/http"
)

// Payload of Request Todo
type ReqTodo struct {
	Name   string `json:"name"`
	Done   *bool  `json:"done"`
	NoteID *int   `json:"noteId"`
}

func (req *ReqTodo) Bind(r *http.Request) error {
	if req.Name == "" || req.NoteID == nil {
		return util.ErrorMissingRequiredFields
	}
	return nil
}

func NewTodoFromReq(req *ReqTodo) model.Todo {
	// 'done' field is optional, we have to check as below
	done := false
	if req.Done != nil {
		done = *req.Done
	}
	return model.Todo{
		Name:   req.Name,
		Done:   done,
		NoteID: *req.NoteID,
	}
}

// Payload of Request List Todo
type RespTodos struct {
	RespNote
	Todos []model.Todo `json:"todos"`
}

func NewRespTodos(note *model.Note, todos []model.Todo) RespTodos {
	return RespTodos{
		RespNote: NewRespNote(note),
		Todos:    todos,
	}
}
