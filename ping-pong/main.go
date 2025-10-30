package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"ping-pong/internal/input"
	"sync"

	"github.com/joho/godotenv"
)

type HandleRequest struct {
	mu      sync.Mutex
	counter int
	path    string
}

func (h *HandleRequest) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	h.counter++
	count := h.counter
	h.mu.Unlock()

	log.Printf("Ping-Pong received a request #%d", count)
	input.WriteToFile(h.path, fmt.Sprintf("%d", count))
	w.Write([]byte(fmt.Sprintf("pong %d", count)))
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	port := os.Getenv("PORT")
	path := os.Getenv("HASH_FILE_PATH")
	if path == "" {
		path = "hashes.log"
	}
	h := &HandleRequest{path: path}

	log.Printf("Starting Ping Pong app on :%s", port)
	http.Handle("/", h)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
