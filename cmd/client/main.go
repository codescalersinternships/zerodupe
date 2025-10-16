package main

import (
	"github.com/codescalersinternships/zerodupe/pkg/client"
)

func main() {
	_, err := client.Upload("./testdata/files/test1.txt")
	if err != nil {
		panic(err)
	}
}
