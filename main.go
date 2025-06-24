// Package main indicates that this is an executable program.
package main

// The import block lists all the external packages our code needs to function.
import (
	"context"
	"encoding/json" // Used for encoding and decoding JSON data.
	"fmt"           // Used for formatted I/O, like printing strings with variables.
	"log"           // Provides logging capabilities.
	"net/http"      // The core package for all HTTP functionality.
	"os"            // Used here to specify the output for our logger (standard output).
	"os/signal"     // Used here to check for interrupt
	"strconv"       // Provides functions to convert strings to other types, like integers.
	"time"          // Used for adding timeout over here.

	"github.com/go-chi/chi/v5" // The chi router we are using.
)

// Item represents the data structure for the items we will store.
// The `json:"..."` tags are called "struct tags". They tell the json package
// how to map the JSON keys to our Go struct fields when encoding and decoding.
type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// server is a struct that holds all the dependencies for our application.
// This is a form of dependency injection, making our app more modular and testable.
type server struct {
	logger    *log.Logger
	router    chi.Router
	datastore map[int]Item // Our simple in-memory database. The key is the item ID.
}

// newServer is the constructor function for our server. It's responsible for
// creating and initializing all the components of our application.
func newServer() *server {
	// Create a new logger that writes to the standard output, with a prefix and standard flags.
	logger := log.New(os.Stdout, "API: ", log.LstdFlags)
	// Create a new chi router instance.
	router := chi.NewRouter()

	// Create an instance of our server struct.
	s := &server{
		logger:    logger,
		router:    router,
		datastore: make(map[int]Item), // Initialize the map! Otherwise, it's nil and will cause a crash.
	}

	// Set up the application's routes.
	s.routes()
	return s
}

// routes defines all the application's API endpoints and maps them to their handlers.
func (s *server) routes() {
	// A POST request to /items will create a new item.
	s.router.Post("/items", s.handleCreateItem())
	// A GET request to /items/{id} will retrieve a specific item.
	s.router.Get("/items/{id}", s.handleGetItem())
	// A PUT request to /items/{id} will update a specific item.
	s.router.Put("/items/{id}", s.handleChangeItem())
	// A GET request to /slow for gracefull shutdown
	s.router.Get("/slow", s.handleSlow())
}

func (s *server) handleSlow() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Println("Starting slow request...")
		time.Sleep(10 * time.Second) // Simulate a long-running task
		s.logger.Println("Finished slow request.")
		fmt.Fprintf(w, "Finally, I am done.")
	}
}

// handleCreateItem handles requests to create a new item (e.g., POST /items).
func (s *server) handleCreateItem() http.HandlerFunc {
	// This is a closure. It returns the actual handler function,
	// which has access to the server `s` and its dependencies.
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a variable to store the JSON data from the request body.
		var newItem Item
		// Create a new JSON decoder that reads from the request body.
		err := json.NewDecoder(r.Body).Decode(&newItem)
		if err != nil {
			// If decoding fails, log the error and send a 400 Bad Request to the client.
			s.logger.Printf("ERROR decoding request body: %v", err)
			http.Error(w, "Bad request: invalid JSON", http.StatusBadRequest)
			return
		}

		// Check if an item with this ID already exists in our datastore.
		_, found := s.datastore[newItem.ID]
		if found {
			s.logger.Printf("Attempted to create item with duplicate ID: %d", newItem.ID)
			// Respond with a 409 Conflict error, which is more specific than 400.
			http.Error(w, fmt.Sprintf("Error: ID %d already in use", newItem.ID), http.StatusConflict)
			return
		}

		// If everything is okay, store the new item in our datastore map.
		s.datastore[newItem.ID] = newItem
		s.logger.Printf("Successfully created and stored item: %+v", newItem)

		// --- Respond to the client ---
		// Set the Content-Type header to inform the client we are sending JSON.
		w.Header().Set("Content-Type", "application/json")
		// Set the HTTP status code to 201 Created.
		w.WriteHeader(http.StatusCreated)
		// Encode the newly created item into JSON and write it to the response.
		json.NewEncoder(w).Encode(newItem)
	}
}

// handleGetItem handles requests to retrieve a single item by its ID (e.g., GET /items/101).
func (s *server) handleGetItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Use chi's URLParam function to get the value of "id" from the URL path.
		idStr := chi.URLParam(r, "id")

		// The ID from the URL is a string, so we need to convert it to an integer.
		id, err := strconv.Atoi(idStr)
		if err != nil {
			s.logger.Printf("ERROR converting ID string to int: %v", err)
			http.Error(w, "Invalid item ID", http.StatusBadRequest)
			return
		}

		// Look up the item in our datastore using the integer ID.
		// The "value, found" is a common Go idiom for checking if a key exists in a map.
		item, found := s.datastore[id]
		if !found {
			s.logger.Printf("Item with ID %d not found", id)
			// If the item doesn't exist, respond with a 404 Not Found error.
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}

		// If the item is found, respond with it.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(item)
	}
}

// handleChangeItem handles requests to update an existing item (e.g., PUT /items/101).
func (s *server) handleChangeItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// --- First, find the existing item just like in handleGetItem ---
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			s.logger.Printf("ERROR converting ID to int: %v", err)
			http.Error(w, "Invalid item ID", http.StatusBadRequest)
			return
		}

		// Check if the item we are trying to update actually exists.
		_, found := s.datastore[id]
		if !found {
			s.logger.Printf("Attempted to update non-existent item with ID %d", id)
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}

		// --- Now, decode the new data from the request body ---
		var updatedItem Item
		err = json.NewDecoder(r.Body).Decode(&updatedItem)
		if err != nil {
			s.logger.Printf("ERROR decoding request body: %v", err)
			http.Error(w, "Bad request: invalid JSON", http.StatusBadRequest)
			return
		}

		// --- Update the item in our datastore ---
		// Enforce the ID from the URL to prevent a mismatch with the body.
		updatedItem.ID = id
		s.datastore[id] = updatedItem // Replace the old item with the new one at the same ID.
		s.logger.Printf("Successfully updated item with ID: %d", id)

		// --- Respond with the updated item ---
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedItem)
	}
}

// main is the entry point for the application.
func main() {
	// Create a new instance of our server with all its dependencies.
	server := newServer()
	server.logger.Println("Server starting on port :8080...")

	// --- Graceful Shutdown Setup ---

	// We create a custom http.Server to have finer control over its behavior.
	srv := &http.Server{
		Addr:    ":8080",
		Handler: server.router, // Our chi router is the handler.
	}

	// Run the server in a goroutine so that it doesn't block the main thread.
	// This allows the main thread to listen for shutdown signals.
	go func() {
		// srv.ListenAndServe() starts the server. It's a blocking call.
		// We check for any error returned by ListenAndServe, ignoring ErrServerClosed,
		// which is the expected error when we gracefully shut down the server.
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			server.logger.Fatalf("Cannot start server: %v", err)
		}
	}() // The `()` immediately invokes the anonymous function.

	// Create a channel to receive OS signals. We buffer it with a size of 1.
	quit := make(chan os.Signal, 1)
	// signal.Notify redirects incoming os.Interrupt signals (like Ctrl+C) to our `quit` channel.
	signal.Notify(quit, os.Interrupt)

	// Block the main goroutine until a signal is received on the `quit` channel.
	<-quit
	server.logger.Println("Shutdown signal received, initiating graceful shutdown...")

	// Create a context with a 5-second timeout to give active connections time to finish.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// `defer cancel()` ensures the context is canceled to release its resources,
	// no matter how the function exits.
	defer cancel()

	// srv.Shutdown() gracefully shuts down the server.
	// It stops accepting new connections and waits for active connections to finish.
	if err := srv.Shutdown(ctx); err != nil {
		server.logger.Fatalf("Server forced to shutdown: %v", err)
	}

	server.logger.Println("Server exited gracefully")
}
