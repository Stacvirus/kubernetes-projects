package hash

import (
	"crypto/sha1"
	"encoding/hex"
)

// Generate creates a SHA1 hash of a fixed message (or could later accept input).
func Generate() string {
	h := sha1.New()
	h.Write([]byte("Hello, World!"))
	sum := h.Sum(nil)
	return hex.EncodeToString(sum)[:10]
}
