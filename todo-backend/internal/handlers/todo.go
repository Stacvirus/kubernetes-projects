package handlers

import (
	"encoding/json"
	"net/http"
	"todo-backend/internal/store"
)

type MemoryStoreHandler struct {
	Store *store.MemoryStore
}

// Creates a new Todo and save in the memory store
func (m *MemoryStoreHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todoReq struct {
		Task string `json:"task"`
	}

	if err := json.NewDecoder(r.Body).Decode(&todoReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	todo := m.Store.AddTodo(todoReq.Task)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// GetAllTodo handles the retrieval of all todo items from the memory store
func (m *MemoryStoreHandler) GetAllTodo(w http.ResponseWriter, r *http.Request) {
	todos := m.Store.GetTodos()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}
