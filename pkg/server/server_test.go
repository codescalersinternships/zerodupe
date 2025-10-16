package server

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"testing"
)

func TestUpload(t *testing.T) {
	data, err := os.ReadFile("../../testdata/files/test1.txt")
	if err != nil {
		t.Errorf("failed to read file: %v", err)
	}

	resp, err := http.Post("http://localhost:8080/upload", "application/octet-stream", bytes.NewReader(data))
	if err != nil {
		t.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %v", resp.Status)
	}
	if string(body) != "9d259d5d9057fec99f14f4025f76188bff1de029cffe020440a474e8739a7719" {
		t.Errorf("unexpected response body: %s", string(body))
	}

}
