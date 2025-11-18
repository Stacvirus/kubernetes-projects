package server

import (
	"net/http"
	"ping-pong/internal/config"
	"ping-pong/internal/handlers"
	"ping-pong/internal/middleware"
	"ping-pong/internal/routes"
	"ping-pong/internal/store"
)

type Server struct {
	router http.Handler
	store  store.Storage
}

func New(cfg *config.Config, store store.Storage) *Server {
	handler := handlers.NewHandleRequest(cfg.Path, store)
	mux := http.NewServeMux()
	routes.RegisterRoutes(mux, handler)

	middleware := middleware.Logging(mux)

	return &Server{
		router: middleware,
		store:  store,
	}
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.router)
}
