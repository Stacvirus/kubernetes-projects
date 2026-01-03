package dto

import (
	"encoding/json"
	"io"

	"github.com/go-playground/validator/v10"
)

type UpdateTodoRequest struct {
	Task string `json:"task" validate:"required,min=3,max=140"`
	Done bool   `json:"done"`
}

func (t *UpdateTodoRequest) FromJson(r io.Reader) error {
	return json.NewDecoder(r).Decode(t)
}

func (t *UpdateTodoRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(t)
}
