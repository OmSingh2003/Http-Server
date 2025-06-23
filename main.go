package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// Item represents our data structure. Note the json tags.
type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// server holds all our application's dependencies.
type server struct {
	logger    *log.Logger
	router    chi.Router
	datastore map[int]Item // Our simple in-memory database. Key is the item ID.
}

// newServer is the constructor for our server. It sets everything up.
func newServer() *server {
	logger := log.New(os.Stdout, "API: ", log.LstdFlags)
	router := chi.NewRouter()

	s := &server{
		logger:    logger,
		router:    router,
		datastore: make(map[int]Item), // Initialize the map!
	}

	// Set up the routes after the server is created.
	s.routes()
	return s
}

// routes defines all the application's endpoints.
func (s *server) routes() {
	s.router.Post("/items", s.handleCreateItem())
	s.router.Get("/items/{id}", s.handleGetItem())
}

// handleCreateItem handles requests to create a new item.
func (s *server) handleCreateItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newItem Item
		// Decode the incoming JSON from the request body.
		err := json.NewDecoder(r.Body).Decode(&newItem)
		if err != nil {
			s.logger.Printf("ERROR decoding request body: %v", err)
			http.Error(w, "Bad request: invalid JSON", http.StatusBadRequest)
			return
		}

		// Store the new item in our datastore map using its ID as the key.
		s.datastore[newItem.ID] = newItem
		s.logger.Printf("Successfully created and stored item: %+v", newItem)

		// Respond to the client.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated) // 201 Created
		json.NewEncoder(w).Encode(newItem)
	}
}

// handleGetItem handles requests to retrieve a single item by its ID.
func (s *server) handleGetItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Use chi.URLParam to get the "id" from the URL path.
		idStr := chi.URLParam(r, "id")

		// Convert the ID from a string to an integer.
		id, err := strconv.Atoi(idStr)
		if err != nil {
			s.logger.Printf("ERROR converting ID to int: %v", err)
			http.Error(w, "Invalid item ID", http.StatusBadRequest)
			return
		}

		// Look up the item in our datastore.
		item, found := s.datastore[id]
		if !found {
			s.logger.Printf("Item with ID %d not found", id)
			http.Error(w, "Item not found", http.StatusNotFound) // 404 Not Found
			return
		}

		// Respond with the found item.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(item)
	}
}

func main() {
	server := newServer()
	server.logger.Println("Server starting on port :8080...")

	// Start the server using the chi router.
	err := http.ListenAndServe(":8080", server.router)
	if err != nil {
		server.logger.Fatalf("Cannot start server: %v", err)
	}
}
