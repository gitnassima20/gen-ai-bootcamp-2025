package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"lang-portal/internal/models"
	"lang-portal/internal/repository"
)

// GroupHandler manages HTTP handlers for group-related operations
type GroupHandler struct {
	groupRepo repository.GroupRepository
}

// NewGroupHandler creates a new instance of GroupHandler
func NewGroupHandler(groupRepo repository.GroupRepository) *GroupHandler {
	return &GroupHandler{
		groupRepo: groupRepo,
	}
}

// GetGroups handles retrieving a list of groups
func (h *GroupHandler) GetGroups(c *gin.Context) {
	// Parse query parameters with defaults
	params := models.DefaultGroupQueryParams()

	// Override with request parameters if provided
	if page, exists := c.GetQuery("page"); exists {
		if pageNum, err := strconv.Atoi(page); err == nil && pageNum > 0 {
			params.Page = pageNum
		}
	}

	if perPage, exists := c.GetQuery("groups_per_page"); exists {
		if num, err := strconv.Atoi(perPage); err == nil && num > 0 {
			params.PerPage = num
		}
	}

	if sortBy, exists := c.GetQuery("sort_by"); exists {
		params.SortBy = sortBy
	}

	if order, exists := c.GetQuery("order"); exists {
		params.Order = order
	}

	// Retrieve groups
	groupsResponse, err := h.groupRepo.GetGroups(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve groups",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, groupsResponse)
}

// GetGroup handles retrieving a specific group by ID
func (h *GroupHandler) GetGroup(c *gin.Context) {
	// Parse group ID from URL
	groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid group ID",
		})
		return
	}

	// Retrieve group
	group, err := h.groupRepo.GetGroupByID(c.Request.Context(), groupID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Group not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, group)
}

// GetGroupWords handles retrieving words in a group
func (h *GroupHandler) GetGroupWords(c *gin.Context) {
	// Parse group ID from URL
	groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid group ID",
		})
		return
	}

	// Parse query parameters with defaults
	params := models.DefaultGroupQueryParams()

	// Override with request parameters if provided
	if page, exists := c.GetQuery("page"); exists {
		if pageNum, err := strconv.Atoi(page); err == nil && pageNum > 0 {
			params.Page = pageNum
		}
	}

	if perPage, exists := c.GetQuery("words_per_page"); exists {
		if num, err := strconv.Atoi(perPage); err == nil && num > 0 {
			params.PerPage = num
		}
	}

	if sortBy, exists := c.GetQuery("sort_by"); exists {
		params.SortBy = sortBy
	}

	if order, exists := c.GetQuery("order"); exists {
		params.Order = order
	}

	// Retrieve group words
	wordsResponse, err := h.groupRepo.GetGroupWords(c.Request.Context(), groupID, params)
	if err != nil {
		// Check if it's a "not found" error
		if err.Error() == "group not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Group not found",
				"details": err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to retrieve group words",
				"details": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, wordsResponse)
}

// GetGroupWordsRaw handles retrieving raw word details for a group
func (h *GroupHandler) GetGroupWordsRaw(c *gin.Context) {
	// Parse group ID from URL
	groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid group ID",
		})
		return
	}

	// Retrieve raw group words
	rawWordsResponse, err := h.groupRepo.GetGroupWordsRaw(c.Request.Context(), groupID)
	if err != nil {
		// Check if it's a "not found" error
		if err.Error() == "group not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Group not found",
				"details": err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to retrieve raw group words",
				"details": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, rawWordsResponse)
}
