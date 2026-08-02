package hash

import (
	"fmt"
	"log"
	"os"
	"time"
)

// WriteToFile appends a hash with a timestamp to the given file.
func WriteToFile(filename, hash string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return
	}
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	line := fmt.Sprintf("%s - %s\n", timestamp, hash)
	if _, err := f.WriteString(line); err != nil {
		log.Printf("Error writing to file: %v", err)
	}
}
