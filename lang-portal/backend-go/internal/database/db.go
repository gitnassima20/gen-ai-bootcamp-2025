package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Common errors
var (
	ErrNotFound = errors.New("record not found")
	ErrDatabase = errors.New("database error")
)

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Path         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
}

// Database wraps the sql.DB connection pool
type Database struct {
	*sql.DB
}

// CreateDatabase creates a new Database instance with the given configuration
func CreateDatabase(cfg DatabaseConfig) (*Database, error) {
	// Validate configuration
	if cfg.Path == "" {
		return nil, fmt.Errorf("database path is required")
	}

	// Set default values if not provided
	if cfg.MaxOpenConns == 0 {
		cfg.MaxOpenConns = 25 // default max open connections
	}
	if cfg.MaxIdleConns == 0 {
		cfg.MaxIdleConns = 25 // default max idle connections
	}
	if cfg.MaxIdleTime == 0 {
		cfg.MaxIdleTime = 15 * time.Minute // default max idle time
	}

	// Open database connection with foreign key support
	db, err := sql.Open("sqlite3", cfg.Path+"?_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxIdleTime(cfg.MaxIdleTime)

	// Create ping context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Verify database connection
	if err := db.PingContext(ctx); err != nil {
		db.Close() // Close the connection if ping fails
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return &Database{DB: db}, nil
}

// Close closes the database connection pool
func (db *Database) Close() error {
	return db.DB.Close()
}

// WithTransaction executes a function within a database transaction
func (db *Database) WithTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	// Attempt to execute the function
	if err := fn(tx); err != nil {
		// Attempt to rollback if the function fails
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("error rolling back transaction: %v (original error: %w)", rbErr, err)
		}
		return err
	}

	// Attempt to commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

// ExecContext executes a query without returning any rows
func (db *Database) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	result, err := db.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v (query: %s)", ErrDatabase, err, query)
	}
	return result, nil
}

// QueryContext executes a query that returns rows
func (db *Database) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := db.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v (query: %s)", ErrDatabase, err, query)
	}
	return rows, nil
}

// QueryRowContext executes a query that returns a single row
func (db *Database) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return db.DB.QueryRowContext(ctx, query, args...)
}
