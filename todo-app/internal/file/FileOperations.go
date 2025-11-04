package file

import (
	"io"
	"log"
	"os"
	"time"
)

// SaveBytesToFile completely overwrites the file with given bytes.
func SaveBytesToFile(filename string, data []byte) error {
	f, err := os.OpenFile(filename, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return err
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		log.Printf("Error writing to file: %v", err)
	}
	return err
}

// ReadFileModTime returns the last modified time of a file.
func ReadFileModTime(filename string) (time.Time, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}

// CopyFile copies contents of src to dst (useful for replacing cached files).
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
