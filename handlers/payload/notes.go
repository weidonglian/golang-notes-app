package payload

import (
	"github.com/weidonglian/golang-notes-app/handlers/util"
	"github.com/weidonglian/golang-notes-app/model"
	"net/http"
)

// Payload of Response Note
type RespNote struct {
	*model.Note
	UseId util.OmitField `json:"userId,omitempty"`
}

func NewRespNote(note *model.Note) RespNote {
	return RespNote{Note: note}
}

func NewRespNoteArray(notes []model.Note) []RespNote {
	respNotes := make([]RespNote, len(notes))
	for i := 0; i < len(notes); i++ {
		respNotes[i] = RespNote{Note: &notes[i]}
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
