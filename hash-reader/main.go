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
	pingPongUrl string
	greeterUrl  string
	message     string
	filePath    string
}

func (p *Path) handler(w http.ResponseWriter, r *http.Request) {
	pingContent, err := external.GetRequest(p.pingPongUrl)
	if err != nil {
		log.Printf("Error fetching pingContent: %v", err)
		http.Error(w, "Error fetching pingContent from pong service", http.StatusInternalServerError)
		return
	}

	greeterContent, err := external.GetRequest(p.greeterUrl)
	if err != nil {
		log.Printf("Error fetching greeter content: %v", err)
		http.Error(w, "Error fetching greeter content from greeter service", http.StatusInternalServerError)
		return
	}
	fileContent := reader.ReadFileContent(p.filePath)
	line := fmt.Sprintf("file content: %s\nenv variable: %s\n%sgreetings: %s", fileContent, p.message, pingContent, greeterContent)

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, line)
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
	port := os.Getenv("PORT")
	pongURL := os.Getenv("PONG_SERVICE_URL")
	greeterURL := os.Getenv("GREETER_SERVICE_URL")
	message := os.Getenv("MESSAGE")
	filePath := os.Getenv("FILE_PATH")

	// path := os.Getenv("HASH_FILE_PATH")
	// if path == "" {
	// 	path = "logs.log"
	// }
	p := &Path{pingPongUrl: pongURL, greeterUrl: greeterURL, message: message, filePath: filePath}

	log.Printf("Starting hash reader app on :%s", port)

	// log.Printf("Reading pings count from file: %s", path)

	http.HandleFunc("/", p.handler)
	http.HandleFunc("/health", healthHandler)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
