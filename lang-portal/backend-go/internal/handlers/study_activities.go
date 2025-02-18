package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lang-portal/internal/repository"
)

type StudyActivityHandler struct {
	studyActivityRepo repository.StudyActivityRepository
}

// NewStudyActivityHandler creates a new handler for study activities
func NewStudyActivityHandler(repo repository.StudyActivityRepository) *StudyActivityHandler {
	return &StudyActivityHandler{
		studyActivityRepo: repo,
	}
}

// ListStudyActivities handles GET /api/v1/study-activities
func (h *StudyActivityHandler) ListStudyActivities(c *gin.Context) {
	// Parse pagination parameters
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("activities_per_page", "100")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 100
	}

	// Fetch study activities
	activities, totalCount, err := h.studyActivityRepo.List(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve study activities",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":        activities,
		"current_page": page,
		"total_pages":  (totalCount + pageSize - 1) / pageSize,
		"total_count":  totalCount,
	})
}

// GetStudyActivity handles GET /api/v1/study-activities/:id
func (h *StudyActivityHandler) GetStudyActivity(c *gin.Context) {
	// Parse activity ID from URL
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid study activity ID",
			"details": "ID must be a valid integer",
		})
		return
	}

	// Fetch specific study activity
	activity, err := h.studyActivityRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Study activity not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, activity)
}
