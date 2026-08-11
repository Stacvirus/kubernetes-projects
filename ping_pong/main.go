package main

import (
	"log"

	"github.com/Stacvirus/hash-generator-app/internal/config"
	"github.com/Stacvirus/hash-generator-app/internal/server"
)

func main() {
	cfg := config.Load()
	s := server.New(cfg)
	log.Printf("🚀 Starting server on :%s", cfg.Port)
	if err := s.Start(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server %v", err)
	}
}
