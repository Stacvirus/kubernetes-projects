package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"todo-backend/internal/models"
	"todo-backend/internal/models/dto"
	"todo-backend/internal/nats"
	"todo-backend/internal/store"

	"github.com/gorilla/mux"
)

type TodoHandler struct {
	Store store.Repository
	Nats  *nats.NatClient
}

// Constructor for RepositoryStoreHandler
func NewTodoHandler(store store.Repository, nats *nats.NatClient) *TodoHandler {
	return &TodoHandler{Store: store, Nats: nats}
}

// Creates a new Todo and save in the memory store
func (t *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	todoReq := &dto.CreateTodoRequest{}
	if err := todoReq.FromJson(r.Body); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := todoReq.Validate(); err != nil {
		log.Println("error validating todo request:", err.Error())
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	todo := &models.Todo{
		Task: todoReq.Task,
	}

	err := t.Store.Todo.AddTodo(r.Context(), todo)
	if err != nil {
		log.Println("Error adding todo:", err)
		http.Error(w, "Failed to create todo", http.StatusInternalServerError)
		return
	}

	// publish message to nats server on another current thread to no block the current flow
	go t.Nats.Publish(todo.Task)

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

func (t *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "missing todo id", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid todo id", http.StatusBadRequest)
		return
	}

	req := &dto.UpdateTodoRequest{}
	if err := req.FromJson(r.Body); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		log.Println("error validating update todo request:", err.Error())
		http.Error(w, "validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Build model
	todo := &models.Todo{
		ID:   id,
		Task: req.Task,
		Done: req.Done,
	}

	// Update in DB
	updatedTodo, err := t.Store.Todo.UpdateTodo(r.Context(), todo)
	if err != nil {
		log.Println("Error updating todo:", err)
		http.Error(w, "failed to update todo", http.StatusInternalServerError)
		return
	}

	// publish message to nats server on another current thread to no block the current flow
	go t.Nats.Publish(todo.Task)

	// Respond
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTodo)
}

func (t *TodoHandler) Health(w http.ResponseWriter, r *http.Request) {
	if !t.Store.Todo.TestDB(r.Context()) {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("healthy"))
}
