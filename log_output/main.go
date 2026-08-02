package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/stacvirus/log_output/internal/hash"
)

func handler(w http.ResponseWriter, r *http.Request) {
	h := hash.Generate()
	log.Printf("%s", h)
	w.Write([]byte(h))
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	port := os.Getenv("PORT")
	log.Printf("Starting log output app on :%s", port)

	path := os.Getenv("HASH_FILE_PATH")
	if path == "" {
		path = "hashes.log"
	}

	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for range ticker.C {
			h := hash.Generate()
			log.Printf("%s", h)
			hash.WriteToFile(path, h)
		}
	}()

	http.HandleFunc("/", handler)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
