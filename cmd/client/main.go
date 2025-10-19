package main

import (
	"github.com/codescalersinternships/zerodupe/pkg/client"
)

func main() {
	// _, err := client.Upload("./testdata/files/test1.txt")
	// if err != nil {
	// 	panic(err)
	// }

	_, err := client.Download("9d259d5d9057fec99f14f4025f76188bff1de029cffe020440a474e8739a7719")
	if err != nil {
		panic(err)
	}

}
