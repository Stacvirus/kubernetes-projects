package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/stacvirus/hash-generator-app/internal/hash"
	"github.com/stacvirus/hash-generator-app/internal/output"
)

type Path struct {
	value string
}

func (p *Path) handler(w http.ResponseWriter, r *http.Request) {
	output := output.ReadFile(p.value)

	hash := hash.Generate()
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	line := fmt.Sprintf("%s %s\n Ping / Pong: %s\n", timestamp, hash, output)

	log.Printf("%s", line)
	w.Write([]byte(line))
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	port := os.Getenv(("PORT"))
	log.Printf("Starting hash app on :%s", port)

	path := os.Getenv("HASH_FILE_PATH")
	if path == "" {
		path = "hashes.log"
	}

	// ticker := time.NewTicker(5 * time.Second)
	// go func() {
	// 	for range ticker.C {
	// 		h := hash.Generate()
	// 		log.Printf("%s", h)
	// 		hash.WriteToFile(path, h)
	// 	}
	// }()
	p := &Path{value: path}

	http.HandleFunc("/", p.handler)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
