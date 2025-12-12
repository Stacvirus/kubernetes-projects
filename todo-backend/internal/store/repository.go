package store

import (
	"context"
	"database/sql"
	"todo-backend/internal/models"
)

type Repository struct {
	Todo interface {
		AddTodo(context.Context, *models.Todo) error
		GetTodos(context.Context) ([]*models.Todo, error)
		TestDB(context.Context) bool
	}
}

func NewRepository(db *sql.DB) Repository {
	return Repository{
		Todo: &TodoRepository{db},
	}
}
