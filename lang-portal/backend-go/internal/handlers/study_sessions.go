package handlers

import (
	"net/http"
	"strconv"
	"time"

	"lang-portal/internal/models"
	"lang-portal/internal/repository"

	"github.com/gin-gonic/gin"
)

type StudySessionHandler struct {
	studySessionRepo repository.StudySessionRepository
}

// NewStudySessionHandler creates a new handler for study sessions
func NewStudySessionHandler(repo repository.StudySessionRepository) *StudySessionHandler {
	return &StudySessionHandler{
		studySessionRepo: repo,
	}
}

// ListStudySessions handles GET /api/v1/study-sessions
func (h *StudySessionHandler) ListStudySessions(c *gin.Context) {
	// Parse query parameters
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("sessions_per_page", "100")
	activityIDStr := c.Query("activity_id")
	groupIDStr := c.Query("group_id")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 100
	}

	// Parse optional activity and group IDs
	var activityID, groupID int64
	if activityIDStr != "" {
		activityID, err = strconv.ParseInt(activityIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid activity ID",
				"details": "Activity ID must be a valid integer",
			})
			return
		}
	}

	if groupIDStr != "" {
		groupID, err = strconv.ParseInt(groupIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid group ID",
				"details": "Group ID must be a valid integer",
			})
			return
		}
	}

	// Fetch study sessions
	sessions, totalCount, err := h.studySessionRepo.List(c.Request.Context(), activityID, groupID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve study sessions",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":        sessions,
		"current_page": page,
		"total_pages":  (totalCount + pageSize - 1) / pageSize,
		"total_count":  totalCount,
	})
}

// CreateStudySession handles POST /api/v1/study-sessions
func (h *StudySessionHandler) CreateStudySession(c *gin.Context) {
	// Define request body struct
	type CreateSessionRequest struct {
		GroupID         int64 `json:"group_id" binding:"required"`
		StudyActivityID int64 `json:"study_activity_id" binding:"required"`
	}

	// Bind and validate request body
	var req CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Create new study session
	session := &models.StudySession{
		GroupID:         req.GroupID,
		StudyActivityID: req.StudyActivityID,
		CreatedAt:       time.Now().UTC(),
	}

	// Save to database
	if err := h.studySessionRepo.Create(c.Request.Context(), session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create study session",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, session)
}

// GetStudySessionDetails handles GET /api/v1/study-sessions/:id
func (h *StudySessionHandler) GetStudySessionDetails(c *gin.Context) {
	// Parse session ID from URL
	sessionIDStr := c.Param("id")
	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid study session ID",
			"details": "ID must be a valid integer",
		})
		return
	}

	// Fetch specific study session details
	sessionDetails, err := h.studySessionRepo.GetStudySessionDetails(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Study session not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, sessionDetails)
}

// CreateWordReview handles POST /study-sessions/:id/words/:word-id/review
func (h *StudySessionHandler) CreateWordReview(c *gin.Context) {
	// Parse study session ID from URL
	sessionIDStr := c.Param("id")
	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid study session ID",
			"details": "Session ID must be a valid integer",
		})
		return
	}

	// Parse word ID from URL
	wordIDStr := c.Param("word-id")
	wordID, err := strconv.ParseInt(wordIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid word ID",
			"details": "Word ID must be a valid integer",
		})
		return
	}

	// Define request body struct
	type WordReviewRequest struct {
		Correct bool `json:"correct"`
	}

	// Bind and validate request body
	var req WordReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Create new word review
	review := &models.WordReviewItem{
		WordID:         wordID,
		StudySessionID: sessionID,
		Correct:        req.Correct,
		CreatedAt:      time.Now().UTC(),
	}

	// Save to database
	if err := h.studySessionRepo.CreateWordReview(c.Request.Context(), review); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create word review",
			"details": err.Error(),
		})
		return
	}

	// Prepare response
	c.JSON(http.StatusCreated, gin.H{
		"success":          true,
		"word_id":          review.WordID,
		"study_session_id": review.StudySessionID,
		"correct":          review.Correct,
		"created_at":       review.CreatedAt.Format(time.RFC3339),
	})
}

// ListStudySessionWords handles GET /api/v1/study-sessions/:id/words
func (h *StudySessionHandler) ListStudySessionWords(c *gin.Context) {
	// Parse study session ID from URL
	sessionIDStr := c.Param("id")
	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid study session ID",
			"details": "Session ID must be a valid integer",
		})
		return
	}

	// Parse pagination parameters
	pageStr := c.DefaultQuery("page", "1")
	wordsPerPageStr := c.DefaultQuery("words_per_page", "100")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	wordsPerPage, err := strconv.Atoi(wordsPerPageStr)
	if err != nil || wordsPerPage < 1 {
		wordsPerPage = 100
	}

	// Fetch words for the study session
	words, totalWords, err := h.studySessionRepo.ListWordsByStudySession(c.Request.Context(), sessionID, page, wordsPerPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve study session words",
			"details": err.Error(),
		})
		return
	}

	// Calculate total pages
	totalPages := (totalWords + wordsPerPage - 1) / wordsPerPage

	// Prepare response
	c.JSON(http.StatusOK, gin.H{
		"items":        words,
		"total_words":  totalWords,
		"current_page": page,
		"total_pages":  totalPages,
	})
}
