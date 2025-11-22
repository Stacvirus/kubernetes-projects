package store

import (
	"context"
	"database/sql"
	"todo-backend/internal/models"
)

type TodoRepository struct {
	db *sql.DB
}

func (r *TodoRepository) AddTodo(ctx context.Context, todo *models.Todo) error {
	query := `
		INSERT INTO todos (task)
		VALUES ($1) RETURNING id, created_at
	`
	err := r.db.QueryRowContext(ctx, query, todo.Task).Scan(&todo.ID, &todo.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *TodoRepository) GetTodos(ctx context.Context) ([]*models.Todo, error) {
	query := `
		SELECT id, task, created_at FROM todos ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*models.Todo
	for rows.Next() {
		todo := &models.Todo{}
		if err := rows.Scan(&todo.ID, &todo.Task, &todo.CreatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}
