package main

import (
	"crypto/sha1"
	"encoding/hex"
	"log"
	"net/http"
	"time"
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
}

func main() {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for range ticker.C {
			hash := generateHash()
			log.Printf("%s", hash)
		}
	}()

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
