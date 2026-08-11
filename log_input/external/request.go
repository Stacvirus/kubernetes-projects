package external

import (
	"fmt"
	"io"
	"net/http"
)

// Forward a GET request to the specified URL and return the response body as a string.
func GetRequest(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error occured during http get request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response string: %w", err)
	}

	return string(body), nil
}
