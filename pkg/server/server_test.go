package server

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"testing"
)

func TestUpload(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		code     int
		body     string
	}{
		{
			name:     "correct testcase",
			filePath: "../../testdata/files/test1.txt",
			code:     http.StatusOK,
			body:     "9d259d5d9057fec99f14f4025f76188bff1de029cffe020440a474e8739a7719",
		},
		{
			name:     "empty file",
			filePath: "../../testdata/files/empty.txt",
			code:     http.StatusOK,
			body:     "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := os.ReadFile(tt.filePath)
			if err != nil {
				t.Fatalf("failed to read file: %v", err)
			}

			resp, err := http.Post("http://localhost:8080/upload", "application/octet-stream", bytes.NewReader(data))
			if err != nil {
				t.Fatalf("failed to send request: %v", err)
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("failed to read response body: %v", err)
			}

			if resp.StatusCode != tt.code {
				t.Errorf("expected status %v, got %v", tt.code, resp.StatusCode)
			}

			if tt.body != "" && string(body) != tt.body {
				t.Errorf("unexpected response body: got %s, want %s", string(body), tt.body)
			}
		})
	}
}
