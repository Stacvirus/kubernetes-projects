package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
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

	const cacheDuration = 10 * time.Minute

	// fs := http.FileServer(http.Dir("./static"))
	template := template.Must(template.ParseFiles("static/index.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := struct{ BackendURL string }{BackendURL: os.Getenv("BACKEND_URL")}
		template.Execute(w, data)
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/image", func(w http.ResponseWriter, r *http.Request) {
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
		http.ServeFile(w, r, imagePath)
	})

	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
