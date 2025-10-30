package input

import (
	"log"
	"os"
)

// WriteToFile completely overwrites the file with the provided value.
func WriteToFile(filename, value string) {
	f, err := os.OpenFile(filename, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(value); err != nil {
		log.Printf("Error writing to file: %v", err)
	}
}
