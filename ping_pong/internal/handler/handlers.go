package handler

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Stacvirus/hash-generator-app/internal/hash"
	"github.com/Stacvirus/hash-generator-app/internal/input"
	"github.com/Stacvirus/hash-generator-app/internal/reader"
)

type HandleRequest struct {
	mu      sync.Mutex
	counter int
	path    string
}

func NewHandleRequest(path string) *HandleRequest {
	return &HandleRequest{path: path}
}

func (h *HandleRequest) RouteHandler(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	h.counter++
	count := h.counter
	h.mu.Unlock()

	log.Printf("Ping pong received a request #%d", count)

	hash := hash.Generate()
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	line := fmt.Sprintf("%s %s\n Ping / Pong: %d\n", timestamp, hash, count)

	input.WriteToFile(h.path, line)
	res := fmt.Sprintf("pong %d", count)
	fmt.Fprint(w, res)
}

func (h *HandleRequest) CounterHandler(w http.ResponseWriter, r *http.Request) {
	content := reader.ReadFileContent(h.path)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, content)
}
