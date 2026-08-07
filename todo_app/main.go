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

	const cacheDuration = 10 * time.Minute
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		needNew := true
		if modTime, err := file.ReadFileModTime(imagePath); err == nil {
			if time.Since(modTime) < cacheDuration {
				needNew = false
			}
		}

		if needNew {
			log.Println("Fetching new image from Picsum...")
			img, err := picsum.DownloadRandomImage(1200)
			if err == nil {
				if err := file.SaveBytesToFile(imagePath, img); err != nil {
					log.Printf("Error saving image: %v", err)
				}
			} else {
				log.Printf("Error downloading image: %v", err)
			}
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// Write simple HTML response
		html := `
		<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Todo App</title>

	<style>
		* {
			box-sizing: border-box;
		}

		body {
			margin: 0;
			padding: 40px 16px;
			background-color: #f7f7f7;
			font-family: Arial, Helvetica, sans-serif;
			color: #292929;
		}

		.todo-app {
			width: 100%;
			max-width: 760px;
			margin: 0 auto;
			text-align: center;
		}

		h1 {
			margin: 0 0 20px;
			font-size: 30px;
		}

		.todo-image {
			display: block;
			width: 180px;
			height: 180px;
			margin: 0 auto 48px;
			object-fit: cover;
			border-radius: 7px;
			box-shadow: 0 2px 5px rgba(0, 0, 0, 0.25);
		}

		.todo-form {
			display: flex;
			justify-content: center;
			gap: 10px;
			width: 100%;
			max-width: 465px;
			margin: 0 auto 26px;
		}

		.todo-input {
			flex: 1;
			min-width: 0;
			padding: 10px;
			border: 2px solid #4caf50;
			border-radius: 4px;
			font-size: 15px;
			outline: none;
		}

		.todo-input:focus {
			border-color: #388e3c;
			box-shadow: 0 0 0 2px rgba(76, 175, 80, 0.15);
		}

		.todo-button {
			padding: 0 20px;
			border: none;
			border-radius: 3px;
			background-color: #4caf50;
			color: white;
			font-size: 15px;
			font-weight: bold;
			cursor: pointer;
		}

		.todo-button:hover {
			background-color: #43a047;
		}

		h2 {
			margin: 0 0 16px;
			font-size: 24px;
		}

		.todo-list {
			margin: 0;
			padding: 0;
			list-style: none;
			text-align: left;
		}

		.todo-item {
			position: relative;
			margin-bottom: 8px;
			padding: 14px 16px;
			background-color: white;
			border-radius: 4px;
			box-shadow: 0 1px 5px rgba(0, 0, 0, 0.1);
			font-size: 15px;
		}

		.todo-item::before {
			content: "";
			position: absolute;
			top: 0;
			bottom: 0;
			left: 0;
			width: 4px;
			background-color: #4caf50;
			border-radius: 4px 0 0 4px;
		}

		@media (max-width: 520px) {
			.todo-form {
				flex-direction: column;
			}

			.todo-button {
				padding: 11px 20px;
			}
		}
	</style>
</head>

<body>
	<main class="todo-app">
		<h1>Todo App</h1>

		<img
			class="todo-image"
			src="/image"
			alt="Random landscape"
		>

		<form class="todo-form">
			<input
				class="todo-input"
				type="text"
				maxlength="140"
				placeholder="Enter a new todo (max 140 characters)"
			>

			<button class="todo-button" type="submit">
				Send
			</button>
		</form>

		<h2>Todos</h2>

		<ul class="todo-list">
			<li class="todo-item">Learn Kubernetes basics</li>
			<li class="todo-item">Deploy application to cluster</li>
			<li class="todo-item">Configure persistent volumes</li>
		</ul>
	</main>
</body>
</html>
		`
		w.Write([]byte(html))
	})

	http.HandleFunc("/image", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, imagePath)
	})
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
