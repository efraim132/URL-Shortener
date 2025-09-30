# URL Shortener (Go)

A minimal in-memory URL shortener written in Go. This service exposes a simple HTTP API to create short URL mappings and redirect requests to the original (long) URL.

## Quick overview / Checklist

- Run the server locally
- Create a short URL mapping (POST /urls)
- Use the short URL (GET /{short})
- Run the test suite

## Prerequisites

- Go 1.18+ installed (ensure `go` is on your PATH)
- (Optional) curl for testing from the command line

## Run the server

From the project directory (where `main.go` lives) run:

```sh
go run main.go
```

The server listens on :8080 by default.

## API

1) Create or update a short URL

- Endpoint: POST /urls
- Content-Type: application/json
- Body JSON shape:
  {
    "short": "gh",
    "long": "https://github.com"
  }

Notes:
- The `short` value may be provided with or without a leading `/`. The server ensures the stored key begins with `/`.
- On success the server returns HTTP 201 Created with an empty body.
- If the request method is not POST, the server returns 405 Method Not Allowed.
- If the JSON is invalid or required fields are missing, the server returns 400 Bad Request.

Example (curl):

```sh
curl -X POST http://localhost:8080/urls \
  -H "Content-Type: application/json" \
  -d "{\"short\": \"gh\", \"long\": \"https://github.com\"}"
```

2) Follow (redirect) a short URL

- Endpoint: GET /{short}
- Example: GET /gh
- If the mapping exists the server responds with a 302 Found and a `Location` header pointing to the long URL.
- If not found, the server returns 404 Not Found.

Example (follow redirects with curl):

```sh
curl -L http://localhost:8080/gh
```

## How it works (implementation notes)

- Storage: an in-memory map stored in the package-level `theWholeStore` variable.
  - Keys are normalized to include a leading slash, e.g. `/gh`.
  - Values are the long/original URL strings.
- Concurrency: access to the map is protected by a `sync.RWMutex`.
  - Reads use `RLock`/`RUnlock` for concurrent read performance.
  - Writes use `Lock`/`Unlock`.
- Handlers:
  - `handlePostURL` parses and validates JSON from the request body and calls `addToWholeStore` to store the mapping.
  - `handleRedirect` looks up the requested path in the map and issues an HTTP 302 redirect if found.
- The server uses the standard library `net/http` (no external dependencies).

## Testing

There is a test file `main_test.go` in the project. To run tests from the command line:

```sh
go test -v
```

In GoLand:
- Open the test file or the package in the Project view.
- Click the green run icon next to a test function to run that single test, or right-click the file/package and choose "Run 'go test'".
- You can also run with coverage or debug tests using the test runner icons.

## Limitations & Next steps

- Persistence: the store is in-memory. Restarting the server loses all mappings. Add a file or database-backed store to persist mappings.
- Validation: currently the server does minimal validation of the `long` URL. Consider validating and normalizing URLs (ensure scheme present, etc.).
- Collisions: currently `addToWholeStore` will overwrite an existing mapping for the same short path. If you want immutability, reject duplicates.
- Security/rate limiting: consider adding rate limits and basic auth if exposing publicly.

## Troubleshooting

- Port 8080 already in use: either stop the other service or change the port in `main.go`.

