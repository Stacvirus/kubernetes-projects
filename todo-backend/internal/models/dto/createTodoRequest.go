package dto

import (
	"encoding/json"
	"io"

	"github.com/go-playground/validator/v10"
)

type CreateTodoRequest struct {
	Task string `json:"task" validate:"required,min=3,max=140"`
}

func (t *CreateTodoRequest) FromJson(r io.Reader) error {
	return json.NewDecoder(r).Decode(t)
}

func (t *CreateTodoRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(t)
}
