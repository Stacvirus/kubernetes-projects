package store

import (
	"sync"
	"todo-backend/internal/models"
)

type MemoryStore struct {
	mu    sync.RWMutex
	todos []models.Todo
}

// memory store constructor
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

// AddTodo adds a new todo to the memory store
func (m *MemoryStore) AddTodo(task string) models.Todo {
	m.mu.Lock()
	defer m.mu.Unlock()

	todo := models.Todo{Task: task}
	m.todos = append(m.todos, todo)
	return todo
}

// GetTodos retrieves all todos from the memory store
func (m *MemoryStore) GetTodos() []models.Todo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return append([]models.Todo(nil), m.todos...)
}
