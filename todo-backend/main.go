package main

import (
	"log"
	"todo-backend/internal/config"
	"todo-backend/internal/server"
)

func main() {
	cfg := config.Load()
	s := server.New()
	log.Printf("🚀 Starting server on :%s", cfg.Port)
	if err := s.Start(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
