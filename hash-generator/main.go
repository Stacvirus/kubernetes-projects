package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/stacvirus/hash-generator-app/internal/hash"
)

func handler(w http.ResponseWriter, r *http.Request) {
	hash := hash.Generate()
	log.Printf("%s", hash)
	w.Write([]byte(hash))
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	port := os.Getenv(("PORT"))
	log.Printf("Starting hash app on :%s", port)

	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for range ticker.C {
			h := hash.Generate()
			log.Printf("%s", h)
			hash.WriteToFile("hashes.log", h)
		}
	}()

	http.HandleFunc("/", handler)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
