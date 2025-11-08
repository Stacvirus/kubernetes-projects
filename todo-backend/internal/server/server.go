package server

import (
	"net/http"
	"todo-backend/internal/handlers"
	"todo-backend/internal/middleware"
	"todo-backend/internal/routes"
	"todo-backend/internal/store"

	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
}

func New() *Server {
	store := store.NewMemoryStore()
	handler := &handlers.MemoryStoreHandler{Store: store}

	r := mux.NewRouter()
	r.Use(middleware.Logging)

	routes.RegisterRoutes(r, handler)

	return &Server{
		router: r,
	}
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.router)
}
