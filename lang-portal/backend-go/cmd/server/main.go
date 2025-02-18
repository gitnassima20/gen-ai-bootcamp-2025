package main

import (
	"log"

	"lang-portal/config"
	"lang-portal/internal/database"
	"lang-portal/internal/handlers"
	"lang-portal/internal/repository"
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

	// Create repositories
	groupRepo := repository.NewGroupRepository(db.DB)
	studyActivityRepo := repository.NewStudyActivityRepository(db.DB)
	studySessionRepo := repository.NewStudySessionRepository(db.DB)
	dashboardRepo := repository.NewDashboardRepository(db.DB)

	// Create handlers
	groupHandler := handlers.NewGroupHandler(groupRepo)
	studyActivityHandler := handlers.NewStudyActivityHandler(studyActivityRepo)
	studySessionHandler := handlers.NewStudySessionHandler(studySessionRepo)
	dashboardHandler := handlers.NewDashboardHandler(dashboardRepo)

	// Setup routes
	router := routes.SetupRoutes(
		groupHandler,
		studyActivityHandler,
		studySessionHandler,
		dashboardHandler,
	)

	// Run server
	if err := routes.RunServer(router, cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
