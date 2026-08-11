package reader

import (
	"log"
	"os"
)

func ReadFileContent(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("Error reading file: %v", err)
		return "Error: unable to read file"
	}
	return string(data)
}
