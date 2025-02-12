# Go Backend Migration Plan

This document outlines the plan to migrate the Flask backend to Go, suitable for junior developers.

## 1. Project Setup (Week 1)

### 1.1 Initial Setup
- Create a new directory `backend-go`
- Initialize Go module: `go mod init lang-portal`
- Set up the project structure:
```
backend-go/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── database/
│   ├── handlers/
│   ├── models/
│   └── middleware/
├── migrations/
├── config/
└── go.mod
```

### 1.2 Dependencies
Install necessary Go packages:
```bash
go get -u github.com/gin-gonic/gin        # Web framework
go get -u github.com/mattn/go-sqlite3     # SQLite driver
go get -u github.com/golang-migrate/migrate # Database migrations
go get -u github.com/rs/cors              # CORS middleware
```

## 2. Database Layer (Week 2)

### 2.1 Database Connection
- Create database package in `internal/database`
- Implement SQLite connection pool
- Create database models for:
  - Words
  - Word Reviews
  - Word Review Items
  - Groups
  - Study Activities

### 2.2 Migrations
- Port existing SQL schemas from `sql/setup/*.sql`
- Create migration files for each table
- Implement migration system using `golang-migrate`

## 3. API Routes Implementation (Week 3-4)

### 3.1 Core Functionality
Implement the following routes (matching current Flask implementation):

#### Words Routes
- GET `/words` - List words
- POST `/words` - Add new word
- GET `/words/:id` - Get word details
- PUT `/words/:id` - Update word
- DELETE `/words/:id` - Delete word

#### Groups Routes
- GET `/groups` - List groups
- POST `/groups` - Create group
- GET `/groups/:id` - Get group details
- PUT `/groups/:id` - Update group
- DELETE `/groups/:id` - Delete group

#### Study Sessions Routes
- GET `/study-sessions` - List sessions
- POST `/study-sessions` - Create session
- GET `/study-sessions/:id` - Get session details

### 3.2 Middleware
- Implement CORS middleware with dynamic origin configuration
- Add request logging
- Add error handling middleware

## 4. Testing (Week 5)

### 4.1 Unit Tests
- Write tests for database operations
- Test route handlers
- Test middleware functions

### 4.2 Integration Tests
- Test API endpoints
- Test database migrations
- Test error scenarios

## 5. Documentation (Throughout)

### 5.1 Code Documentation
- Add godoc comments to all exported functions
- Document database schema
- Document API endpoints

### 5.2 Setup Instructions
- Write clear setup instructions
- Document environment variables
- Add example API requests

## 6. Performance Considerations

### 6.1 Optimizations
- Implement connection pooling
- Add caching where appropriate
- Use goroutines for concurrent operations

## Best Practices

1. **Error Handling**
   - Use custom error types
   - Implement proper error wrapping
   - Return appropriate HTTP status codes

2. **Code Organization**
   - Follow standard Go project layout
   - Use interfaces for better testing
   - Implement dependency injection

3. **Security**
   - Sanitize user inputs
   - Implement proper CORS policies
   - Use prepared statements for SQL

## Learning Resources

1. [Go Documentation](https://golang.org/doc/)
2. [Gin Web Framework](https://github.com/gin-gonic/gin)
3. [Go Database Tutorial](https://golang.org/pkg/database/sql/)
4. [Go Project Layout](https://github.com/golang-standards/project-layout)

## Timeline

- Week 1: Project Setup and Environment Configuration
- Week 2: Database Layer Implementation
- Week 3-4: API Routes Implementation
- Week 5: Testing and Documentation
- Week 6: Review and Optimization

## Getting Started

1. Install Go (1.21 or later)
2. Clone the repository
3. Run `go mod tidy`
4. Copy configuration files
5. Run migrations
6. Start the server: `go run cmd/server/main.go`

Remember to commit changes frequently and write clear commit messages!
