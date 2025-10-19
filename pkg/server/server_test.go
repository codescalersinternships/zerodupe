package server

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func countFiles(dir string) int {
	files, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	return len(files)
}

func TestUpload_NoDuplication(t *testing.T) {
	storagePath := "../../storage"

	os.RemoveAll(storagePath)
	os.MkdirAll(filepath.Join(storagePath, BlocksDir), 0755)
	os.MkdirAll(filepath.Join(storagePath, MetadataDir), 0755)

	data, err := os.ReadFile("../../testdata/files/test1.txt")
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	// 1st upload
	resp1, err := http.Post("http://localhost:8080/upload", "application/octet-stream", bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to send 1st request: %v", err)
	}
	resp1.Body.Close()

	count1 := countFiles(filepath.Join(storagePath, BlocksDir))

	// 2nd upload
	resp2, err := http.Post("http://localhost:8080/upload", "application/octet-stream", bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to send 2nd request: %v", err)
	}
	resp2.Body.Close()

	count2 := countFiles(filepath.Join(storagePath, BlocksDir))

	if count2 != count1 {
		t.Error("duplicate blocks found")
	}
}
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

func TestDownload(t *testing.T) {
	tests := []struct {
		name         string
		expectedFile string
		code         int
		passedHash   string
	}{
		{
			name:         "correct testcase",
			expectedFile: "../../testdata/files/test1.txt",
			code:         http.StatusOK,
			passedHash:   "9d259d5d9057fec99f14f4025f76188bff1de029cffe020440a474e8739a7719",
		},
		{
			name:         "empty file",
			expectedFile: "../../testdata/files/empty.txt",
			code:         http.StatusOK,
			passedHash:   "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name:       "file hash not passed",
			code:       http.StatusNotFound,
			passedHash: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			endpoint := "http://localhost:8080/download" + "/" + tt.passedHash
			resp, err := http.Get(endpoint)
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

			if tt.code != http.StatusNotFound {
				expectedFileData, err := os.ReadFile(tt.expectedFile)
				if err != nil {
					t.Fatalf("failed to read file: %v", err)
				}

				if string(expectedFileData) != string(body) {
					t.Errorf("unexpected response body: got %s, want %s", string(body), string(expectedFileData))
				}
			}

		})
	}
}
