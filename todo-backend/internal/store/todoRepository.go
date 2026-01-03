package store

import (
	"context"
	"database/sql"
	"todo-backend/internal/models"
)

type TodoRepository struct {
	db *sql.DB
}

func (r *TodoRepository) TestDB(ctx context.Context) bool {
	err := r.db.PingContext(ctx)
	return err == nil
}

func (r *TodoRepository) AddTodo(ctx context.Context, todo *models.Todo) error {
	query := `
		INSERT INTO todos (task)
		VALUES ($1) RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query, todo.Task).Scan(&todo.ID, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *TodoRepository) UpdateTodo(ctx context.Context, todo *models.Todo) (*models.Todo, error) {
	query := `
		UPDATE todos
		SET task = $1, done = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
		RETURNING id, task, done, created_at, updated_at
	`
	updatedTodo := &models.Todo{}
	err := r.db.QueryRowContext(ctx, query, todo.Task, todo.Done, todo.ID).Scan(&updatedTodo.ID, &updatedTodo.Task, &updatedTodo.Done, &updatedTodo.CreatedAt, &updatedTodo.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return updatedTodo, nil
}

func (r *TodoRepository) GetTodos(ctx context.Context) ([]*models.Todo, error) {
	query := `
		SELECT id, task, done, created_at, updated_at FROM todos ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*models.Todo
	for rows.Next() {
		todo := &models.Todo{}
		if err := rows.Scan(&todo.ID, &todo.Task, &todo.Done, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}
