package picsum

import (
	"fmt"
	"io"
	"net/http"
)

// DownloadRandomImage fetches a random image from Picsum with given width
// and returns its raw bytes.
func DownloadRandomImage(width int) ([]byte, error) {
	url := fmt.Sprintf("https://picsum.photos/%d", width)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading image: %w", err)
	}

	return data, nil
}
