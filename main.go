package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type server struct {
	router *http.ServeMux
	logger *log.Logger
}

func newServer() *server {
	newLogger := log.New(os.Stdout, "WEB: ", log.LstdFlags)

	// Create the server instance.
	s := &server{
		router: http.NewServeMux(),
		logger: newLogger,
	}

	// Set up the routes for the server we just created.
	s.routes()

	return s
}

func (s *server) helloHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Printf("received request for %s from %s", r.URL.Path, r.RemoteAddr)
		fmt.Fprintf(w, "Hello from a server that logs things!")
	}
}

func (s *server) goodbyeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Printf("received request for %s from %s ", r.URL.Path, r.RemoteAddr)
		fmt.Fprintf(w, "Goodbye from a server that logs things!")
	}
}

func (s *server) itemHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Printf("received request for %s from %s ", r.URL.Path, r.RemoteAddr)
		switch r.Method {
		case http.MethodGet:
			fmt.Fprintf(w, "Here are all the items that u requested")
		case http.MethodPost:
			fmt.Fprintf(w, "new item is inserted")
		case http.MethodPut:
			fmt.Fprintf(w, "item data is changed")
		default:
			// It's good practice to handle other methods too.
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func (s *server) routes() {
	s.router.HandleFunc("/", s.helloHandler())
	s.router.HandleFunc("/goodbye", s.goodbyeHandler())
	s.router.HandleFunc("/items", s.itemHandler())
}

func main() {
	server := newServer()

	server.logger.Printf("Server is starting on port :8080")

	err := http.ListenAndServe(":8080", server.router)
	if err != nil {
		server.logger.Fatalf("Cannot start server: %v", err)
	}
}
