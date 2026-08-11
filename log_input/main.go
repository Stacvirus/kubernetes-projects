package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/stacvirus/log_input/external"
)

type Path struct {
	url string
}

func (p *Path) handler(w http.ResponseWriter, r *http.Request) {
	content, err := external.GetRequest(p.url)
	if err != nil {
		log.Printf("Error fetching content: %v", err)
		http.Error(w, "Error fetching content from pong service", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, content)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting log input app on :%s", port)

	pongService := os.Getenv("PONG_SERVICE_URL")
	if pongService == "" {
		log.Fatalf("PONG_SERVICE_URL environment variable is not set")
	}
	p := &Path{url: pongService}
	log.Printf("Getting` hash from service: %s", p.url)

	http.HandleFunc("/", p.handler)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
