package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"lang-portal/internal/repository"
)

// DashboardHandler handles HTTP requests related to dashboard
type DashboardHandler struct {
	dashboardRepo repository.DashboardRepository
}

// NewDashboardHandler creates a new handler for dashboard
func NewDashboardHandler(repo repository.DashboardRepository) *DashboardHandler {
	return &DashboardHandler{dashboardRepo: repo}
}

// GetLastStudySession handles GET /api/v1/dashboard/last-study-session
func (h *DashboardHandler) GetLastStudySession(c *gin.Context) {
	// Fetch the last study session
	lastSession, err := h.dashboardRepo.GetLastStudySession(c.Request.Context())
	if err != nil {
		if err.Error() == "no study sessions found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "No study sessions found",
				"details": err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to retrieve last study session",
				"details": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, lastSession)
}

// GetStudyProgress handles GET /api/v1/dashboard/study-progress
func (h *DashboardHandler) GetStudyProgress(c *gin.Context) {
	// Fetch study progress
	progress, err := h.dashboardRepo.GetStudyProgress(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve study progress",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, progress)
}

// GetQuickStats handles GET /api/v1/dashboard/quick-stats
func (h *DashboardHandler) GetQuickStats(c *gin.Context) {
	// Fetch quick stats
	stats, err := h.dashboardRepo.GetQuickStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve quick stats",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}
