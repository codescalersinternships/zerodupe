package server

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	BlockSize   = 4 * 1024 // 4 kB
	StorageRoot = "./storage"
	BlocksDir   = "blocks"
	MetadataDir = "meta"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to zerodupe!")
}

func concatHashes(bloackHashes []string) []byte {
	var result string
	for _, hash := range bloackHashes {
		result += hash + "\n"
	}
	return []byte(result)
}

func Upload(w http.ResponseWriter, r *http.Request) {
	//create dirs if not exist
	os.MkdirAll(filepath.Join(StorageRoot, BlocksDir), 0755)
	os.MkdirAll(filepath.Join(StorageRoot, MetadataDir), 0755)

	var blockHashes []string

	for {
		block := make([]byte, BlockSize)
		n, err := io.ReadFull(r.Body, block)

		if err == io.EOF {
			break
		}

		if err != nil {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}
		block = block[:n] //in case of last block being smaller than BlockSize

		blockHash := sha256.Sum256(block)
		blockHashStr := fmt.Sprintf("%x", blockHash)
		blockHashes = append(blockHashes, blockHashStr)

		blockPath := filepath.Join(StorageRoot, BlocksDir, blockHashStr[0:4], blockHashStr)
		if _, err := os.Stat(blockPath); os.IsNotExist(err) {
			os.MkdirAll(filepath.Dir(blockPath), 0755)

			err := os.WriteFile(blockPath, block, 0644)
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to store block: %v", err), http.StatusInternalServerError)
				return
			}
		}
	}

	fileHash := sha256.New()
	for _, hash := range blockHashes {
		fileHash.Write([]byte(hash))
	}
	fileHashStr := fmt.Sprintf("%x", fileHash.Sum(nil))

	metadataPath := filepath.Join(StorageRoot, MetadataDir, fileHashStr[0:4], fileHashStr)

	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(metadataPath), 0755)

		filemeta := concatHashes(blockHashes)
		err := os.WriteFile(metadataPath, filemeta, 0644)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to store file metadata: %v", err), http.StatusInternalServerError)
			return
		}
	}

	//return file hash in response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fileHashStr))

}

func Download(w http.ResponseWriter, r *http.Request) {
	requestedFileHash := strings.TrimPrefix(r.URL.Path, "/download/")

	if requestedFileHash == "" {
		http.Error(w, "file hash not provided", http.StatusBadRequest)
		return
	}

	metadataPath := filepath.Join(StorageRoot, MetadataDir, requestedFileHash[0:4], requestedFileHash)
	metadata, err := os.ReadFile(metadataPath)
	if err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	blockHashes := strings.Split(string(metadata), "\n")

	for _, blockHash := range blockHashes {
		if blockHash == "" {
			continue
		}
		blockPath := filepath.Join(StorageRoot, BlocksDir, blockHash[0:4], blockHash)
		blockData, err := os.ReadFile(blockPath)
		if err != nil {
			http.Error(w, "failed to read block", http.StatusInternalServerError)
			return
		}
		w.Write(blockData)
	}
}
