package server

import (
	"net/http"
	"ping-pong/internal/config"
	"ping-pong/internal/handlers"
	"ping-pong/internal/middleware"
	"ping-pong/internal/routes"
)

type Server struct {
	router http.Handler
}

func New(cfg *config.Config) *Server {
	handler := handlers.NewHandleRequest(cfg.Path)
	mux := http.NewServeMux()
	routes.RegisterRoutes(mux, handler)

	middleware := middleware.Logging(mux)

	return &Server{
		router: middleware,
	}
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.router)
}
