package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"lang-portal/config"
	"lang-portal/internal/database"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Configure database
	dbConfig := database.DatabaseConfig{
		Path:         cfg.DatabasePath,
		MaxOpenConns: 25,
		MaxIdleConns: 25,
	}

	// Create database connection
	db, err := database.CreateDatabase(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := runMigrations(db.DB); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migrations completed successfully")
}

func runMigrations(db *sql.DB) error {
	// Get the directory of the current executable
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}
	migrationDir := filepath.Join(filepath.Dir(execPath), "..", "..", "migrations")

	// Read migration files
	files, err := os.ReadDir(migrationDir)
	if err != nil {
		return fmt.Errorf("failed to read migration directory: %w", err)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Begin a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Execute each migration file
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			migrationPath := filepath.Join(migrationDir, file.Name())
			
			// Read migration file
			migrationSQL, err := os.ReadFile(migrationPath)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to read migration file %s: %w", file.Name(), err)
			}

			// Execute migration
			_, err = tx.ExecContext(ctx, string(migrationSQL))
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to execute migration %s: %w", file.Name(), err)
			}

			log.Printf("Applied migration: %s", file.Name())
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migrations: %w", err)
	}

	return nil
}
