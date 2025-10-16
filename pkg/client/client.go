package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

func Upload(filepath string) (*http.Response, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	resp, err := http.Post("http://localhost:8080/upload", "application/octet-stream", bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("Response status:", resp.Status)
	fmt.Printf("Response Body:\n%s\n", body)

	return resp, nil
}
