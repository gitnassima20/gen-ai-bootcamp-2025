package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"lang-portal/internal/models"
)

// GroupRepository defines methods for group-related database operations
type GroupRepository interface {
	// GetGroups retrieves a list of groups with pagination and sorting
	GetGroups(ctx context.Context, params models.GroupQueryParams) (*models.GroupsResponse, error)

	// GetGroupByID retrieves a specific group by its ID
	GetGroupByID(ctx context.Context, groupID int64) (*models.Group, error)

	// GetGroupWords retrieves words in a group with pagination and sorting
	GetGroupWords(ctx context.Context, groupID int64, params models.GroupQueryParams) (*models.GroupWordsResponse, error)

	// GetGroupWordsRaw retrieves raw word details for a group
	GetGroupWordsRaw(ctx context.Context, groupID int64) (*models.RawGroupWordsResponse, error)
}

// SQLGroupRepository implements GroupRepository for SQL databases
type SQLGroupRepository struct {
	db *sql.DB
}

// NewGroupRepository creates a new instance of SQLGroupRepository
func NewGroupRepository(db *sql.DB) *SQLGroupRepository {
	return &SQLGroupRepository{db: db}
}

// GetGroups retrieves a list of groups with pagination and sorting
func (r *SQLGroupRepository) GetGroups(ctx context.Context, params models.GroupQueryParams) (*models.GroupsResponse, error) {
	// Validate and sanitize sorting parameters
	validSortColumns := map[string]bool{
		"name":        true,
		"words_count": true,
	}
	if !validSortColumns[params.SortBy] {
		params.SortBy = "name"
	}
	if params.Order != "asc" && params.Order != "desc" {
		params.Order = "asc"
	}

	// Construct the order by clause
	orderBy := fmt.Sprintf("%s %s", params.SortBy, params.Order)

	// Calculate offset
	offset := (params.Page - 1) * params.PerPage

	// Query to get groups
	query := fmt.Sprintf(`
		SELECT id, name, words_count
		FROM groups
		ORDER BY %s
		LIMIT ? OFFSET ?
	`, orderBy)

	rows, err := r.db.QueryContext(ctx, query, params.PerPage, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query groups: %w", err)
	}
	defer rows.Close()

	var groups []models.Group
	for rows.Next() {
		var group models.Group
		if err := rows.Scan(&group.ID, &group.Name, &group.WordsCount); err != nil {
			return nil, fmt.Errorf("failed to scan group: %w", err)
		}
		groups = append(groups, group)
	}

	// Get total groups count for pagination
	var totalGroups int
	err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM groups").Scan(&totalGroups)
	if err != nil {
		return nil, fmt.Errorf("failed to count total groups: %w", err)
	}
	totalPages := (totalGroups + params.PerPage - 1) / params.PerPage

	return &models.GroupsResponse{
		Groups:      groups,
		TotalPages:  totalPages,
		CurrentPage: params.Page,
	}, nil
}

// GetGroupByID retrieves a specific group by its ID
func (r *SQLGroupRepository) GetGroupByID(ctx context.Context, groupID int64) (*models.Group, error) {
	query := `
		SELECT id, name, words_count
		FROM groups
		WHERE id = ?
	`
	var group models.Group
	err := r.db.QueryRowContext(ctx, query, groupID).Scan(&group.ID, &group.Name, &group.WordsCount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("group not found: %w", err)
		}
		return nil, fmt.Errorf("failed to retrieve group: %w", err)
	}

	return &group, nil
}

// GetGroupWords retrieves words in a group with pagination and sorting
func (r *SQLGroupRepository) GetGroupWords(ctx context.Context, groupID int64, params models.GroupQueryParams) (*models.GroupWordsResponse, error) {
	// Validate and sanitize sorting parameters
	validSortColumns := map[string]string{
		"kanji":         "w.kanji",
		"romaji":        "w.romaji",
		"english":       "w.english",
		"correct_count": "COALESCE(wr.correct_count, 0)",
		"wrong_count":   "COALESCE(wr.wrong_count, 0)",
	}
	sqlColumn, ok := validSortColumns[params.SortBy]
	if !ok {
		params.SortBy = "kanji"
		sqlColumn = validSortColumns["kanji"]
	}
	if params.Order != "asc" && params.Order != "desc" {
		params.Order = "asc"
	}

	// Calculate offset
	offset := (params.Page - 1) * params.PerPage

	// First, verify the group exists
	_, err := r.GetGroupByID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	// Query to get words in the group
	query := fmt.Sprintf(`
		SELECT 
			w.id, 
			w.kanji, 
			w.romaji, 
			w.english, 
			COALESCE(wr.correct_count, 0) as correct_count,
			COALESCE(wr.wrong_count, 0) as wrong_count
		FROM words w
		JOIN word_groups wg ON w.id = wg.word_id
		LEFT JOIN word_reviews wr ON w.id = wr.word_id
		WHERE wg.group_id = ?
		ORDER BY %s %s
		LIMIT ? OFFSET ?
	`, sqlColumn, strings.ToUpper(params.Order))

	rows, err := r.db.QueryContext(ctx, query, groupID, params.PerPage, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query group words: %w", err)
	}
	defer rows.Close()

	var words []models.GroupWord
	for rows.Next() {
		var word models.GroupWord
		if err := rows.Scan(
			&word.ID, 
			&word.Kanji, 
			&word.Romaji, 
			&word.English, 
			&word.CorrectCount, 
			&word.WrongCount,
		); err != nil {
			return nil, fmt.Errorf("failed to scan word: %w", err)
		}
		words = append(words, word)
	}

	// Get total words count for pagination
	var totalWords int
	err = r.db.QueryRowContext(ctx, 
		"SELECT COUNT(*) FROM word_groups WHERE group_id = ?", 
		groupID,
	).Scan(&totalWords)
	if err != nil {
		return nil, fmt.Errorf("failed to count total words: %w", err)
	}
	totalPages := (totalWords + params.PerPage - 1) / params.PerPage

	return &models.GroupWordsResponse{
		Words:       words,
		TotalPages:  totalPages,
		CurrentPage: params.Page,
	}, nil
}

// GetGroupWordsRaw retrieves raw word details for a group
func (r *SQLGroupRepository) GetGroupWordsRaw(ctx context.Context, groupID int64) (*models.RawGroupWordsResponse, error) {
	// First, verify the group exists and get its name
	group, err := r.GetGroupByID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	// Query to get raw word details
	query := `
		SELECT 
			w.id, 
			w.kanji, 
			w.romaji, 
			w.english, 
			w.parts
		FROM words w
		JOIN word_groups wg ON w.id = wg.word_id
		WHERE wg.group_id = ?
	`

	rows, err := r.db.QueryContext(ctx, query, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to query raw group words: %w", err)
	}
	defer rows.Close()

	var words []models.RawWord
	for rows.Next() {
		var word models.RawWord
		var partsJSON []byte
		if err := rows.Scan(
			&word.ID, 
			&word.Kanji, 
			&word.Romaji, 
			&word.English, 
			&partsJSON,
		); err != nil {
			return nil, fmt.Errorf("failed to scan raw word: %w", err)
		}
		
		// Parse parts JSON
		word.Parts = json.RawMessage(partsJSON)
		words = append(words, word)
	}

	return &models.RawGroupWordsResponse{
		GroupID:   groupID,
		GroupName: group.Name,
		Words:     words,
	}, nil
}
