package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/stacvirus/log_input/internal/reader"
)

type Path struct {
	value string
}

func (p *Path) handler(w http.ResponseWriter, r *http.Request) {
	content := reader.ReadFileContent(p.value)
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

	path := os.Getenv("HASH_FILE_PATH")
	if path == "" {
		path = "../log_output/hashes.log"
	}
	p := &Path{value: path}
	log.Printf("Reading hash from file: %s", p.value)

	http.HandleFunc("/hashes", p.handler)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
