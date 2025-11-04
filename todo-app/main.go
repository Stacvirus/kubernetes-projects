package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"todo-app/internal/file"
	"todo-app/internal/picsum"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	port := os.Getenv("PORT")

	path := os.Getenv("CACHE_FILE_PATH")
	if path == "" {
		path = "./image"
	}
	imagePath := filepath.Join(path, "image.jpg")

	log.Printf("Starting TODO app on :%s", port)

	const cacheDuration = 1 * time.Minute
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check if cached image is still fresh
		needNew := true
		if modTime, err := file.ReadFileModTime(imagePath); err == nil {
			if time.Since(modTime) < cacheDuration {
				needNew = false
			}
		}

		if needNew {
			log.Println("⏳ Fetching new image from Picsum...")
			img, err := picsum.DownloadRandomImage(1200)
			if err == nil {
				if err := file.SaveBytesToFile(imagePath, img); err != nil {
					log.Printf("Failed to save image: %v", err)
				}
			} else {
				log.Printf("Failed to fetch image: %v", err)
			}
		}
		// Tell the browser we’re sending HTML
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// Write simple HTML response
		html := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<title>TODO App</title>
		</head>
		<body>
			<h1>The project App</h1>
			<img src="/image" alt="Random image" width="200"/>
			<p>DevOps with Kubernetes 2025</p>
		</body>
		</html>
		`

		// Write it to the response
		w.Write([]byte(html))
	})

	http.HandleFunc("/image", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, imagePath)
	})

	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
