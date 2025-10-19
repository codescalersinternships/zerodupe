package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

func Upload(filepath string) (string, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	resp, err := http.Post("http://localhost:8080/upload", "application/octet-stream", bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return string(body), fmt.Errorf("failed to read response body: %v", err)
	}

	return string(body), nil
}

func Download(fileHash string) (string, error) {

	endpoint := "http://localhost:8080/download" + "/" + fileHash

	resp, err := http.Get(endpoint)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return string(body), fmt.Errorf("failed to read response body: %v", err)
	}


	return string(body), nil
}
