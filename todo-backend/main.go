package main

import (
	"log"
	"todo-backend/internal/config"
	"todo-backend/internal/db"
	"todo-backend/internal/server"
	"todo-backend/internal/store"
)

func main() {
	cfg := config.Load()

	db, err := db.New(cfg)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	log.Printf("database connection pool established")

	repository := store.NewRepository(db)

	s := server.New(repository)
	log.Printf("🚀 Starting server on :%s", cfg.Port)
	if err := s.Start(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
