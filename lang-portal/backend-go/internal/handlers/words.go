package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lang-portal/internal/models"
	"lang-portal/internal/repository"
)

type WordHandler struct {
	wordRepo repository.WordRepository
}

func NewWordHandler(wordRepo repository.WordRepository) *WordHandler {
	return &WordHandler{wordRepo: wordRepo}
}

// GetWords retrieves a list of words with optional filtering
func (h *WordHandler) GetWords(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// Parse filter parameters
	filter := repository.WordFilter{
		Kanji:   c.Query("kanji"),
		Romaji:  c.Query("romaji"),
		English: c.Query("english"),
	}

	// Parse optional group ID filter
	groupID, _ := strconv.ParseInt(c.Query("groupId"), 10, 64)
	filter.GroupID = groupID

	// Retrieve words
	words, totalCount, err := h.wordRepo.List(c.Request.Context(), filter, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve words",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"words":      words,
		"page":       page,
		"pageSize":   pageSize,
		"totalCount": totalCount,
	})
}

// CreateWord adds a new word to the database
func (h *WordHandler) CreateWord(c *gin.Context) {
	var word models.Word
	if err := c.ShouldBindJSON(&word); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validate required fields
	if word.Kanji == "" || word.Romaji == "" || word.English == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Kanji, Romaji, and English are required",
		})
		return
	}

	// Ensure parts is valid JSON
	if !json.Valid(word.Parts) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid parts JSON",
		})
		return
	}

	// Create word
	if err := h.wordRepo.Create(c.Request.Context(), &word); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create word",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, word)
}

// GetWord retrieves a specific word by ID
func (h *WordHandler) GetWord(c *gin.Context) {
	// Parse word ID from URL
	wordID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid word ID",
		})
		return
	}

	// Retrieve word
	word, err := h.wordRepo.GetByID(c.Request.Context(), wordID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Word not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve word",
				"details": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, word)
}

// UpdateWord modifies an existing word
func (h *WordHandler) UpdateWord(c *gin.Context) {
	// Parse word ID from URL
	wordID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid word ID",
		})
		return
	}

	// Bind request body
	var word models.Word
	if err := c.ShouldBindJSON(&word); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Set the ID from URL parameter
	word.ID = wordID

	// Validate required fields
	if word.Kanji == "" || word.Romaji == "" || word.English == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Kanji, Romaji, and English are required",
		})
		return
	}

	// Update word
	if err := h.wordRepo.Update(c.Request.Context(), &word); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update word",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, word)
}

// DeleteWord removes a word from the database
func (h *WordHandler) DeleteWord(c *gin.Context) {
	// Parse word ID from URL
	wordID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid word ID",
		})
		return
	}

	// Delete word
	if err := h.wordRepo.Delete(c.Request.Context(), wordID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete word",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// AddWordToGroup adds a word to a specific group
func (h *WordHandler) AddWordToGroup(c *gin.Context) {
	// Parse group request
	var groupRequest struct {
		WordID  int64 `json:"wordId"`
		GroupID int64 `json:"groupId"`
	}
	if err := c.ShouldBindJSON(&groupRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validate input
	if groupRequest.WordID == 0 || groupRequest.GroupID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "WordID and GroupID are required",
		})
		return
	}

	// Add word to group
	if err := h.wordRepo.AddToGroup(c.Request.Context(), groupRequest.WordID, groupRequest.GroupID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add word to group",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Word added to group successfully",
	})
}

// RemoveWordFromGroup removes a word from a specific group
func (h *WordHandler) RemoveWordFromGroup(c *gin.Context) {
	// Parse group request
	var groupRequest struct {
		WordID  int64 `json:"wordId"`
		GroupID int64 `json:"groupId"`
	}
	if err := c.ShouldBindJSON(&groupRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validate input
	if groupRequest.WordID == 0 || groupRequest.GroupID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "WordID and GroupID are required",
		})
		return
	}

	// Remove word from group
	if err := h.wordRepo.RemoveFromGroup(c.Request.Context(), groupRequest.WordID, groupRequest.GroupID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to remove word from group",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Word removed from group successfully",
	})
}
