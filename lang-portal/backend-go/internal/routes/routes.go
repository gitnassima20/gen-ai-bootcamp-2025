package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"lang-portal/internal/handlers"
)

// SetupRoutes configures and returns the main router with all API routes
func SetupRoutes(
	groupHandler *handlers.GroupHandler,
	studyActivityHandler *handlers.StudyActivityHandler,
	studySessionHandler *handlers.StudySessionHandler,
	dashboardHandler *handlers.DashboardHandler,
) *gin.Engine {
	router := gin.Default()

	// API versioning
	v1 := router.Group("/api/v1")
	{
		// Groups routes
		groups := v1.Group("/groups")
		{
			groups.GET("", groupHandler.GetGroups)
			groups.GET("/:id", groupHandler.GetGroup)
			groups.GET("/:id/words", groupHandler.GetGroupWords)
			groups.GET("/:id/words/raw", groupHandler.GetGroupWordsRaw)
			groups.GET("/:id/study-sessions", groupHandler.GetGroupStudySessions)
		}

		// Study Activities routes
		studyActivities := v1.Group("/study-activities")
		{
			studyActivities.GET("", studyActivityHandler.ListStudyActivities)
			studyActivities.GET("/:id", studyActivityHandler.GetStudyActivity)
		}

		// Study Sessions routes
		studySessions := v1.Group("/study-sessions")
		{
			studySessions.GET("", studySessionHandler.ListStudySessions)
			studySessions.POST("", studySessionHandler.CreateStudySession)
			studySessions.GET("/:id", studySessionHandler.GetStudySessionDetails)
			studySessions.GET("/:id/words", studySessionHandler.ListStudySessionWords)
			studySessions.POST("/:id/words/:word-id/review", studySessionHandler.CreateWordReview)
		}

		// Dashboard routes
		dashboard := v1.Group("/dashboard")
		{
			dashboard.GET("/last-study-session", dashboardHandler.GetLastStudySession)
			dashboard.GET("/study-progress", dashboardHandler.GetStudyProgress)
			dashboard.GET("/quick-stats", dashboardHandler.GetQuickStats)
		}
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
func RunServer(router *gin.Engine, port string) error {
	// Remove leading ':' if present and use localhost
	if len(port) > 0 && port[0] == ':' {
		port = port[1:]
	}

	// Set up the address to listen on
	address := "localhost:" + port

	fmt.Printf("Starting server on %s\n", address)
	return router.Run(address)
}
