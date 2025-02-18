package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lang-portal/internal/repository"
)

// GroupHandler handles HTTP requests related to groups
type GroupHandler struct {
	groupRepo repository.GroupRepository
}

// NewGroupHandler creates a new handler for groups
func NewGroupHandler(repo repository.GroupRepository) *GroupHandler {
	return &GroupHandler{groupRepo: repo}
}

// GetGroups handles GET /api/v1/groups
func (h *GroupHandler) GetGroups(c *gin.Context) {
	// Parse pagination parameters
	pageStr := c.DefaultQuery("page", "1")
	groupsPerPageStr := c.DefaultQuery("groups_per_page", "100")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	groupsPerPage, err := strconv.Atoi(groupsPerPageStr)
	if err != nil || groupsPerPage < 1 {
		groupsPerPage = 100
	}

	// Fetch groups
	groups, totalGroups, err := h.groupRepo.List(c.Request.Context(), page, groupsPerPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve groups",
			"details": err.Error(),
		})
		return
	}

	// Calculate total pages
	totalPages := (totalGroups + groupsPerPage - 1) / groupsPerPage

	// Prepare response
	c.JSON(http.StatusOK, gin.H{
		"items":        groups,
		"total_count":  totalGroups,
		"current_page": page,
		"total_pages":  totalPages,
	})
}

// GetGroup handles GET /api/v1/groups/:id
func (h *GroupHandler) GetGroup(c *gin.Context) {
	// Parse group ID from URL
	groupIDStr := c.Param("id")
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid group ID",
			"details": "Group ID must be a valid integer",
		})
		return
	}

	// Fetch group details
	group, err := h.groupRepo.GetByID(c.Request.Context(), groupID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Group not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, group)
}

// GetGroupWords handles GET /api/v1/groups/:id/words
func (h *GroupHandler) GetGroupWords(c *gin.Context) {
	// TODO: Implement group words endpoint
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Not implemented",
	})
}

// GetGroupWordsRaw handles GET /api/v1/groups/:id/words/raw
func (h *GroupHandler) GetGroupWordsRaw(c *gin.Context) {
	// TODO: Implement raw group words endpoint
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Not implemented",
	})
}
