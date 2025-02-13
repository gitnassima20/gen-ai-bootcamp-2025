package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"lang-portal/internal/handlers"
	"lang-portal/internal/middleware"
	"lang-portal/internal/repository"
	"log"
	"strings"
)

// SetupRoutes configures and returns a Gin router with all API routes
func SetupRoutes(db *sql.DB) *gin.Engine {
	// Create repositories
	wordRepo := repository.NewWordRepository(db)

	// Create handlers
	wordHandler := handlers.NewWordHandler(wordRepo)

	// Create router
	router := gin.Default()

	// Apply CORS middleware
	router.Use(middleware.CORSMiddleware())

	// API group
	v1 := router.Group("/api/v1")
	{
		// Word routes
		words := v1.Group("/words")
		{
			words.GET("", wordHandler.GetWords)
			words.POST("", wordHandler.CreateWord)
			words.GET("/:id", wordHandler.GetWord)
			words.PUT("/:id", wordHandler.UpdateWord)
			words.DELETE("/:id", wordHandler.DeleteWord)

			// Word group management with different route structure
			words.POST("/groups", wordHandler.AddWordToGroup)
			words.DELETE("/groups", wordHandler.RemoveWordFromGroup)
		}

		// TODO: Add routes for other resources (groups, study activities, etc.)
	}

	// Health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})

	return router
}

// RunServer starts the HTTP server
func RunServer(router *gin.Engine, port string) {
	// Remove leading ':' if present and use localhost
	if strings.HasPrefix(port, ":") {
		port = port[1:]
	}

	// Set up the address to listen on
	address := "localhost:" + port

	log.Printf("Starting server on %s", address)
	if err := router.Run(address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
