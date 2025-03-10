package apiserver

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
	"olbcloud.com/webapi/internal/handlers"
)

func NewServer(hd handlers.Handlers) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /posts", hd.GetPostsHandler)
	mux.HandleFunc("POST /posts", hd.CreatePostHandler)
	mux.HandleFunc("PUT /posts/{id}", hd.UpdatePostHandler)
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

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := LogMiddleware(c.Handler(mux))

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}
