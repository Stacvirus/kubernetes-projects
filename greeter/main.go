package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("Hello from version %s\n", os.Getenv("VERSION"))))
}

// Health endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"status":"ok"}`)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	port := os.Getenv(("PORT"))
	log.Printf("Starting greeter app on: %s", port)

	http.HandleFunc("/", handler)
	http.HandleFunc("/health", healthHandler)

	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
