package output

import (
	"log"
	"os"
)

// ReadFile reads the contents of the given file and returns it as a string.
func ReadFile(filename string) string {
	f, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("Error opening file: %v", err)
	}
	defer f.Close()

	data := make([]byte, 1024)
	n, err := f.Read(data)
	if err != nil {
		log.Printf("Error reading file: %v", err)
	}
	return string(data[:n])
}
