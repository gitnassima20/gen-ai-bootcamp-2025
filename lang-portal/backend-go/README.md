# Lang Portal Backend (Go)

This is the Go implementation of the Lang Portal backend.

## Project Structure

```
backend-go/
├── cmd/
│   └── server/        # Main application entry point
├── internal/
│   ├── database/      # Database connection and queries
│   ├── handlers/      # HTTP request handlers
│   ├── models/        # Data models
│   └── middleware/    # Middleware components
├── migrations/        # Database migrations
├── config/           # Configuration management
└── go.mod            # Go module file
```

## Setup

1. Install Go 1.21 or later
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Run the server:
   ```bash
   go run cmd/server/main.go
   ```

## Development

- The server runs on port 8080 by default
- SQLite database is used for storage
- API endpoints follow RESTful conventions

## TODO

- Implement database operations
- Add authentication middleware
- Complete API endpoints
- Add tests
- Configure CORS
- Set up CI/CD
