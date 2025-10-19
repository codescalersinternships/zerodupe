package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/codescalersinternships/zerodupe/pkg/client"
)

func main() {
	uploadPath := flag.String("upload", "", "Path to the file to upload")
	downloadHash := flag.String("download", "", "Hash of the file to download")
	outputFile := flag.String("out", "downloaded", "output path for downloaded file")
	flag.Parse()

	if *uploadPath != "" {
		hash, err := client.Upload(*uploadPath)
		if err != nil {
			panic(err)
		}
		fmt.Println("Upload successful, file hash:", hash)

	} else if *downloadHash != "" {
		file, err := client.Download(*downloadHash)
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(*outputFile, []byte(file), 0644)
		if err != nil {
			panic(err)
		}
		fmt.Println("Download successful. file saved at ", *outputFile)
	}

}
