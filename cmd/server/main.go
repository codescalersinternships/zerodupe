package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/codescalersinternships/zerodupe/pkg/server"
	router "github.com/go-chi/chi/v5"
)

func handleRequests() {
	r := router.NewRouter()

	r.Get("/", server.HomePage)
	r.Post("/upload", server.Upload)
	r.Get("/download/{fileHash}", server.Download)

	fmt.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func main() {
	handleRequests()
}
