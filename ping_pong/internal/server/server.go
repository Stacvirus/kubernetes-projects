package server

import (
	"net/http"

	"github.com/Stacvirus/hash-generator-app/internal/config"
	"github.com/Stacvirus/hash-generator-app/internal/handler"
	"github.com/Stacvirus/hash-generator-app/internal/middleware"
	"github.com/Stacvirus/hash-generator-app/internal/routes"
)

type Server struct {
	router http.Handler
}

func New(config *config.Config) *Server {
	handler := handler.NewHandleRequest(config.Path)
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
