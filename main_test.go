// The `_test.go` suffix tells the `go test` command that this file contains tests.
package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHandleCreateItem is a test function for our handleCreateItem handler.
// Test functions in Go must start with `Test` and take a `*testing.T` argument.
func TestHandleCreateItem(t *testing.T) {
	// 1. Create a new instance of our server. This gives us a fresh, clean
	// datastore for each test run.
	server := newServer()

	// 2. Create the JSON payload for our request body.
	// We use a struct to ensure it's well-formed and then marshal it to bytes.
	itemPayload := Item{ID: 101, Name: "Test Item", Age: 999}
	body, _ := json.Marshal(itemPayload)

	// 3. Create a new HTTP request object.
	// We use `bytes.NewReader` to create an `io.Reader` from our byte slice.
	req, err := http.NewRequest("POST", "/items", bytes.NewReader(body))
	if err != nil {
		// If we can't even create the request, the test should fail immediately.
		t.Fatalf("could not create request: %v", err)
	}

	// 4. Create a "Response Recorder".
	// This is a special tool from httptest that acts like a ResponseWriter
	// but records the response's status code, headers, and body for us to check.
	rr := httptest.NewRecorder()

	// --- Act ---
	// 5. Call the handler.
	// Instead of starting a full server, we can directly serve the HTTP request
	// to our router. The router will then call the correct handler.
	server.router.ServeHTTP(rr, req)

	// --- Assert ---
	// 6. Check the status code.
	// We assert that the response recorder got a `201 Created` status code.
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// 7. Check the response body.
	// We decode the JSON from the response body into a struct.
	var responseItem Item
	err = json.NewDecoder(rr.Body).Decode(&responseItem)
	if err != nil {
		t.Fatalf("could not decode response body: %v", err)
	}

	// We assert that the item returned in the response is the same as what we sent.
	if responseItem.ID != itemPayload.ID || responseItem.Name != itemPayload.Name {
		t.Errorf("handler returned unexpected body: got %+v want %+v",
			responseItem, itemPayload)
	}
}
