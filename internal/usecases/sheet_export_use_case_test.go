package usecases

import (
	"errors"
	"export-service/internal/adapters"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"testing"
)

func TestSheetExportUseCase_downloadData(t *testing.T) {
	s := SheetExportUseCase{downloader: &adapters.HTTPDownloader{}, logger: zap.NewExample()}
	t.Run("Should download data", func(t *testing.T) {
		r := ExportRequest{
			DataDownloadURL: getTestURL(t, serveJSON),
		}

		data, err := s.downloadData(r)
		require.NoError(t, err)

		assert.Lenf(t, data, 1, "Should have 1 company")

		assert.Equalf(t, "driva-tech", data[0]["public_id"], "Public ID should be driva-tech")
		assert.Lenf(t, data[0]["profiles"].([]any), 10, "Should have 10 profiles")
	})

	t.Run("Should fail if value is not an array", func(t *testing.T) {
		r := ExportRequest{
			DataDownloadURL: getTestURL(t, serveHTML),
		}

		_, err := s.downloadData(r)
		assert.Errorf(t, err, "Should fail if value is not an array")
	})

	t.Run("Should fail if status is not 200", func(t *testing.T) {
		r := ExportRequest{
			DataDownloadURL: getTestURL(t, serveErrorStatus),
		}

		_, err := s.downloadData(r)
		assert.Errorf(t, err, "Should fail if status is not 200")
	})
}

func serveJSON(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("test_data/export_request.json")
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func serveHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<!DOCTYPE html><html><body></body></html>"))
}

func serveErrorStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("[]"))
}

func getTestURL(t *testing.T, handler http.HandlerFunc) string {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Println("Error starting test server:", err)
		t.Fail()
	}

	port := listener.Addr().(*net.TCPAddr).Port
	server := &http.Server{Handler: handler}

	go func() {
		if err := server.Serve(listener); !errors.Is(err, http.ErrServerClosed) {
			log.Println("Test server closed with error:", err)
			t.Fail()
		}
	}()

	return fmt.Sprintf("http://0.0.0.0:%d", port)
}
