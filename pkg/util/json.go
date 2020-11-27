package util

import (
	"encoding/json"
)

type OmitField *struct{}

func ToJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}
	return string(b)
}
