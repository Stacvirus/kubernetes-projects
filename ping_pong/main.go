package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type HandleRequest struct {
	mu      sync.Mutex
	counter int
}

func (h *HandleRequest) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	h.counter++
	count := h.counter
	h.mu.Unlock()

	log.Printf("Ping-Pong received a request #%d", count)
	w.Write([]byte(fmt.Sprintf("pong %d", count)))
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	port := os.Getenv("PORT")

	log.Printf("Starting Ping Pong app on :%s", port)
	http.Handle("/", &HandleRequest{})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
