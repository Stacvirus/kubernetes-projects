package routes

import (
	"net/http"
	"todo-backend/internal/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, handler *handlers.MemoryStoreHandler) {
	r.HandleFunc("/todos", handler.GetAllTodo).Methods(http.MethodGet)
	r.HandleFunc("/todos", handler.CreateTodo).Methods(http.MethodPost)
}
