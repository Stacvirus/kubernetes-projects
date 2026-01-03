package routes

import (
	"net/http"
	"todo-backend/internal/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, handler *handlers.TodoHandler) {
	r.HandleFunc("/todos", handler.GetAllTodo).Methods(http.MethodGet)
	r.HandleFunc("/todos", handler.CreateTodo).Methods(http.MethodPost)
	r.HandleFunc("/todos/{id}", handler.UpdateTodo).Methods(http.MethodPut)

	r.HandleFunc("/healthz", handler.Health).Methods(http.MethodGet)
}
