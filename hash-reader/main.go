package main

import (
	"fmt"
	"hash-reader/external"
	"hash-reader/internal/reader"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Path struct {
	url      string
	message  string
	filePath string
}

func (p *Path) handler(w http.ResponseWriter, r *http.Request) {
	content, err := external.GetRequest(p.url)
	if err != nil {
		log.Printf("Error fetching content: %v", err)
		http.Error(w, "Error fetching content from pong service", http.StatusInternalServerError)
		return
	}
	fileContent := reader.ReadFileContent(p.filePath)
	line := fmt.Sprintf("file content: %s\nenv variable: %s\n%s", fileContent, p.message, content)

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, line)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	port := os.Getenv("PORT")
	pongURL := os.Getenv("PONG_SERVICE_URL")
	message := os.Getenv("MESSAGE")
	filePath := os.Getenv("FILE_PATH")

	// path := os.Getenv("HASH_FILE_PATH")
	// if path == "" {
	// 	path = "logs.log"
	// }
	p := &Path{url: pongURL, message: message, filePath: filePath}

	log.Printf("Starting hash reader app on :%s", port)

	// log.Printf("Reading pings count from file: %s", path)

	http.HandleFunc("/", p.handler)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
