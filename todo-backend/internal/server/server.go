package server

import (
	"net/http"
	"todo-backend/internal/handlers"
	"todo-backend/internal/middleware"
	"todo-backend/internal/routes"
	"todo-backend/internal/store"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Server struct {
	router *mux.Router
}

func New() *Server {
	store := store.NewMemoryStore()
	handler := handlers.NewMemoryStoreHandler(store)

	r := mux.NewRouter()
	r.Use(middleware.Logging)

	routes.RegisterRoutes(r, handler)

	return &Server{
		router: r,
	}
}

func (s *Server) Start(addr string) error {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	handler := c.Handler(s.router)

	return http.ListenAndServe(addr, handler)
}
