package adapters

import (
	"export-service/internal/core/ports"
	"fmt"
	"io"
	"net/http"
)

type HTTPDownloader struct {
}

var _ ports.Downloader = (*HTTPDownloader)(nil)

func (h *HTTPDownloader) Download(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
