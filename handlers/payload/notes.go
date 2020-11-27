package payload

import (
	"github.com/weidonglian/notes-app/handlers/util"
	"github.com/weidonglian/notes-app/model"
	"github.com/weidonglian/notes-app/store"
	"net/http"
)

// Payload of Response Note
type RespNote struct {
	*model.Note
	UserId util.OmitField `json:"userId,omitempty"`
	Todos  []model.Todo   `json:"todos"`
}

func NewRespNote(note *model.Note, todosStore store.TodosStore) RespNote {
	return RespNote{Note: note, Todos: todosStore.FindByNoteID(note.ID)}
}

func NewRespNoteArray(notes []model.Note, todosStore store.TodosStore) []RespNote {
	respNotes := make([]RespNote, len(notes))
	for i := 0; i < len(notes); i++ {
		respNotes[i] = NewRespNote(&notes[i], todosStore)
	}
	return respNotes
}

// Payload of Request Note
type ReqNote struct {
	Name string `json:"name"`
}

func (req *ReqNote) Bind(r *http.Request) error {
	if req.Name == "" {
		return util.ErrorMissingRequiredFields
	}
	return nil
}
