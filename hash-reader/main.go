package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"hash-reader/internal/reader"

	"github.com/joho/godotenv"
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

	path := os.Getenv("HASH_FILE_PATH")
	if path == "" {
		path = "../hash-generator/hashes.log"
	}
	p := &Path{value: path}

	log.Printf("Starting hash reader app on :%s", port)

	log.Printf("Reading hashes from file: %s", path)

	http.HandleFunc("/hashes", p.handler)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
