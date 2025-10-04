package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func generateHash() string {
	h := sha1.New()
	h.Write([]byte("Hello, World!"))
	sum := h.Sum(nil)

	return hex.EncodeToString(sum)[:10]
}

func handler(w http.ResponseWriter, r *http.Request) {
	hash := generateHash()
	log.Printf("%s", hash)
	w.Write([]byte(hash))
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	port := os.Getenv(("PORT"))
	log.Printf("Starting TODO app on :%s", port)

	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for range ticker.C {
			hash := generateHash()
			log.Printf("%s", hash)
		}
	}()

	http.HandleFunc("/", handler)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
