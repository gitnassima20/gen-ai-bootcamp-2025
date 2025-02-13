package main

import (
	"log"

	"lang-portal/config"
	"lang-portal/internal/database"
	"lang-portal/internal/routes"
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

	// Setup routes
	router := routes.SetupRoutes(db.DB)

	// Run server
	routes.RunServer(router, cfg.ServerPort)
}
