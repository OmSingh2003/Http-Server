# Professional Go HTTP Server

This project is a hands-on implementation of professional best practices for building robust, maintainable, and production-ready HTTP services in Go. The architecture and patterns used here are inspired by the insightful Grafana engineering blog post, ["How I write HTTP services in Go after 13 years"](https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/).

The goal of this repository is to serve as a practical example of moving beyond basic tutorials to build a service that incorporates industry-standard techniques for structure, reliability, and testability.

## Key Features Implemented

This server is more than just a simple "Hello, World." It's built with a focus on professional software engineering principles:

-   **Structured Application:** Uses a central `server` struct for clean dependency injection, holding the router, logger, and data store.
-   **Advanced Routing:** Leverages the `chi` router for powerful and flexible routing, including dynamic URL parameters.
-   **Graceful Shutdown:** Implements a graceful shutdown mechanism to ensure the server finishes active requests before stopping, preventing data loss and client errors.
-   **Middleware:** Features a logging middleware that automatically logs the details of every incoming request, keeping handler logic clean and focused.
-   **RESTful API:** Provides a RESTful API for managing "items" with full CRUD (Create, Read, Update) functionality (POST, GET, PUT).
-   **Automated Testing:** Includes an initial test suite using Go's built-in `httptest` package to programmatically verify API endpoint functionality.

---

## Getting Started

### Prerequisites

-   Go (version 1.18 or later) installed on your system.

### Installation & Setup

1.  **Clone the repository:**
    ```sh
    git clone https://github.com/OmSingh2003/Http-Server.git
    cd Http-Server
    ```

2.  **Install dependencies:**
    The project uses Go Modules. The `chi` router will be downloaded automatically when you build or run the project. You can also fetch it manually:
    ```sh
    go mod tidy
    ```

### Running the Server

To start the API server, run the following command from the root of the project directory:

```sh
go run .
```

You should see a log message indicating that the server has started on port 8080:

```
API: 2025/06/24 12:00:00 Server starting on port :8080...
```

## API Endpoints

The server exposes the following endpoints for managing items. You can use a tool like curl to interact with them.

### 1. Create a New Item

**Method:** POST

**Endpoint:** /items

**Body:** JSON payload representing the item.

**Example curl command:**

```sh
curl -X POST -H "Content-Type: application/json" -d '{"id": 101, "name": "Alice", "age": 30}' http://localhost:8080/items
```

### 2. Get a Specific Item

**Method:** GET

**Endpoint:** /items/{id}

**Example curl command:**

```sh
curl http://localhost:8080/items/101
```

### 3. Update an Existing Item

**Method:** PUT

**Endpoint:** /items/{id}

**Body:** JSON payload with the updated item details.

**Example curl command:**

```sh
curl -X PUT -H "Content-Type: application/json" -d '{"id": 101, "name": "Alice Smith", "age": 31}' http://localhost:8080/items/101
```

## Running Tests

This project includes an automated test suite. To run the tests, use the standard go test command. The -v flag provides verbose output.

```sh
go test -v
```

You should see an output indicating that the tests have passed:

```
=== RUN   TestHandleCreateItem
--- PASS: TestHandleCreateItem (0.00s)
PASS
ok      github.com/OmSingh2003/Http-Server    0.006s
```
