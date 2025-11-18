package main

import (
	"log"

	"ping-pong/internal/config"
	"ping-pong/internal/db"
	"ping-pong/internal/server"
	"ping-pong/internal/store"
)

func main() {
	cfg := config.Load()

	db, err := db.New(cfg)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	log.Printf("database connection pool established")

	store := store.NewStorage(db)

	s := server.New(cfg, store)
	log.Printf("🚀 Starting server on :%s", cfg.Port)

	if err := s.Start(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
