package handlers

import (
	"fmt"
	"log"
	"net/http"
	"ping-pong/internal/hash"
	"ping-pong/internal/input"
	"ping-pong/internal/reader"
	"sync"
	"time"
)

type HandleRequest struct {
	mu      sync.Mutex
	counter int
	path    string
}

// constructor for HandleRequest
func NewHandleRequest(path string) *HandleRequest {
	return &HandleRequest{path: path}
}

func (h *HandleRequest) RouteHandler(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	h.counter++
	count := h.counter
	h.mu.Unlock()

	log.Printf("Ping-Pong received a request #%d", count)

	hash := hash.Generate()
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	line := fmt.Sprintf("%s %s\n Ping / Pong: %d\n", timestamp, hash, count)

	input.WriteToFile(h.path, line)
	res := fmt.Sprintf("pong %d", count)
	fmt.Fprint(w, res)
}

func (h *HandleRequest) CountHandler(w http.ResponseWriter, r *http.Request) {
	content := reader.ReadFileContent(h.path)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, content)
}
