package reader

import (
	"io/ioutil"
	"log"
)

// ReadFileContent reads the entire content of a file and returns it as a string.
func ReadFileContent(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("Error reading file: %v", err)
		return "Error: unable to read hashes file"
	}
	return string(data)
}
