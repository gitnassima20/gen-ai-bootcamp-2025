package routes

import (
	"database/sql"
	"lang-portal/internal/handlers"
	"lang-portal/internal/middleware"
	"lang-portal/internal/repository"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures and returns a Gin router with all API routes
func SetupRoutes(db *sql.DB) *gin.Engine {
	// Create repositories
	wordRepo := repository.NewWordRepository(db)
	groupRepo := repository.NewGroupRepository(db)
	studyActivityRepo := repository.NewStudyActivityRepository(db)
	studySessionRepo := repository.NewStudySessionRepository(db)

	// Create handlers
	wordHandler := handlers.NewWordHandler(wordRepo)
	groupHandler := handlers.NewGroupHandler(groupRepo)
	studyActivityHandler := handlers.NewStudyActivityHandler(studyActivityRepo)
	studySessionHandler := handlers.NewStudySessionHandler(studySessionRepo)

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

			// Word group management
			words.POST("/groups", wordHandler.AddWordToGroup)
			words.DELETE("/groups", wordHandler.RemoveWordFromGroup)
		}

		// Group routes
		groups := v1.Group("/groups")
		{
			groups.GET("", groupHandler.GetGroups)
			groups.GET("/:id", groupHandler.GetGroup)
			//TODO: Retest this endpoint
			groups.GET("/:id/words", groupHandler.GetGroupWords)
			groups.GET("/:id/words/raw", groupHandler.GetGroupWordsRaw)
		}

		// Study Activities routes
		studyActivities := v1.Group("/study-activities")
		{
			studyActivities.GET("", studyActivityHandler.ListStudyActivities)
			studyActivities.GET("/:id", studyActivityHandler.GetStudyActivity)
		}

		// Study Sessions routes
		v1.GET("/study-sessions", studySessionHandler.ListStudySessions)
		v1.POST("/study-sessions", studySessionHandler.CreateStudySession)
		v1.GET("/study-sessions/:id", studySessionHandler.GetStudySession)
		v1.POST("/study-sessions/:id/words/:word-id/review", studySessionHandler.CreateWordReview)

		// TODO: Add routes for other resources
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
