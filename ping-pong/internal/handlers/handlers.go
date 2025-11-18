package handlers

import (
	"fmt"
	"log"
	"net/http"
	"ping-pong/internal/hash"
	"ping-pong/internal/store"
	"sync"
	"time"
)

type HandleRequest struct {
	mu      sync.Mutex
	counter int
	path    string
	store   store.Storage
}

// constructor for HandleRequest
func NewHandleRequest(path string, store store.Storage) *HandleRequest {
	return &HandleRequest{path: path, store: store}
}

func (h *HandleRequest) RouteHandler(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	h.counter++
	count := h.counter
	h.mu.Unlock()

	log.Printf("Ping-Pong received a request #%d", count)

	hash := hash.Generate()
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	line := &store.Line{
		Content: fmt.Sprintf("%s %s\n Ping / Pong: %d\n", timestamp, hash, count),
	}

	if err := h.store.Lines.Create(r.Context(), line); err != nil {
		log.Panic(err)
		http.Error(w, "Failed to store line", http.StatusInternalServerError)
		return
	}
}

func (h *HandleRequest) CountHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	line, err := h.store.Lines.ReadLatest(ctx)
	if err != nil {
		log.Panic(err)
		http.Error(w, "Failed to read latest line", http.StatusInternalServerError)
		return
	}

	if line == nil {
		fmt.Fprint(w, "No entries yet")
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, line.Content)
}
