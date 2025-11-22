package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"todo-backend/internal/models"
	"todo-backend/internal/store"
)

type TodoHandler struct {
	Store store.Repository
}

// Constructor for MemoryStoreHandler
func NewTodoHandler(store store.Repository) *TodoHandler {
	return &TodoHandler{Store: store}
}

// Creates a new Todo and save in the memory store
func (m *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todoReq struct {
		Task string `json:"task"`
	}

	if err := json.NewDecoder(r.Body).Decode(&todoReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	todo := &models.Todo{
		Task: todoReq.Task,
	}

	err := m.Store.Todo.AddTodo(r.Context(), todo)
	if err != nil {
		log.Println("Error adding todo:", err)
		http.Error(w, "Failed to create todo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// GetAllTodo handles the retrieval of all todo items from the memory store
func (t *TodoHandler) GetAllTodo(w http.ResponseWriter, r *http.Request) {
	todos, err := t.Store.Todo.GetTodos(r.Context())
	if err != nil {
		log.Println("Error retrieving todos:", err)
		http.Error(w, "Failed to retrieve todos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}
