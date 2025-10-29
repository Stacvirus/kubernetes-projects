package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"hash-reader/internal/reader"

	"github.com/joho/godotenv"
)

func handler(w http.ResponseWriter, r *http.Request) {
	content := reader.ReadFileContent("hashes.log")
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, content)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	port := os.Getenv("PORT")

	log.Printf("Starting hash reader app on :%s", port)

	http.HandleFunc("/hashes", handler)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
