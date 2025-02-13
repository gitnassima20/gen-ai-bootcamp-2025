package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"lang-portal/internal/database"
)

func main() {
	log.Println("Starting migration...")
	
	// Define flags
	dbPath := flag.String("db", "langportal.db", "Path to SQLite database")
	flag.Parse()

	// Read migration file
	log.Println("Reading migration file...")
	migrationSQL, err := os.ReadFile("migrations/001_initial_schema.sql")
	if err != nil {
		log.Fatalf("Failed to read migration file: %v", err)
	}
	log.Printf("Successfully read migration file (%d bytes)", len(migrationSQL))

	// Configure database
	log.Printf("Configuring database with path: %s", *dbPath)
	cfg := database.Config{
		Path:         *dbPath,
		MaxOpenConns: 1, // Only need one connection for migrations
		MaxIdleConns: 1,
	}

	// Connect to database
	log.Println("Attempting to connect to database...")
	db, err := database.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Successfully connected to database")
	defer db.Close()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test database connection
	log.Println("Testing database connection...")
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Database connection test successful")

	log.Println("Running migrations...")
	// Execute migration with context
	_, err = db.ExecContext(ctx, string(migrationSQL))
	if err != nil {
		log.Fatalf("Failed to execute migration: %v", err)
	}

	log.Println("Migrations completed successfully!")
}
