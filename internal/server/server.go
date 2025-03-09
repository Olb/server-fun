package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"olbcloud.com/webapi/internal/handlers"
)

func NewServer(hd handlers.Handlers) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /posts", hd.GetPostsHandler)
	mux.HandleFunc("GET /posts/{id}", hd.GetPostByIDHandler)

	return mux
}

func StartServer(mux *http.ServeMux) {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting server on %s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
