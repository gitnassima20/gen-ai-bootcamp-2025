package repository

import (
	"context"
	"database/sql"
	"fmt"
)

// GroupListItem represents a group in the list view
type GroupListItem struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	WordCount int    `json:"word_count"`
}

// GroupDetails represents detailed information about a group
type GroupDetails struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	TotalWordCount int    `json:"total_word_count"`
}

// GroupRepository defines the interface for group-related database operations
type GroupRepository interface {
	// List retrieves groups with optional pagination
	List(ctx context.Context, page, groupsPerPage int) ([]GroupListItem, int, error)

	// GetByID retrieves detailed information about a specific group
	GetByID(ctx context.Context, groupID int64) (*GroupDetails, error)
}

// SQLGroupRepository implements GroupRepository using SQLite
type SQLGroupRepository struct {
	db *sql.DB
}

// NewGroupRepository creates a new instance of SQLGroupRepository
func NewGroupRepository(db *sql.DB) *SQLGroupRepository {
	return &SQLGroupRepository{db: db}
}

// List retrieves groups with pagination
func (r *SQLGroupRepository) List(ctx context.Context, page, groupsPerPage int) ([]GroupListItem, int, error) {
	// Calculate pagination
	offset := (page - 1) * groupsPerPage

	// Count total groups
	countQuery := `SELECT COUNT(*) FROM groups`
	var totalGroups int
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&totalGroups)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count groups: %w", err)
	}

	// Query to fetch groups with word count
	query := `
		SELECT 
			g.id, 
			g.name, 
			g.words_count
		FROM groups g
		ORDER BY g.name
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, groupsPerPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list groups: %w", err)
	}
	defer rows.Close()

	var groups []GroupListItem
	for rows.Next() {
		var group GroupListItem
		if err := rows.Scan(
			&group.ID,
			&group.Name,
			&group.WordCount,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan group: %w", err)
		}
		groups = append(groups, group)
	}

	return groups, totalGroups, nil
}

// GetByID retrieves detailed information about a specific group
func (r *SQLGroupRepository) GetByID(ctx context.Context, groupID int64) (*GroupDetails, error) {
	query := `
		SELECT 
			id, 
			name, 
			words_count
		FROM groups
		WHERE id = ?
	`
	var group GroupDetails
	err := r.db.QueryRowContext(ctx, query, groupID).Scan(
		&group.ID,
		&group.Name,
		&group.TotalWordCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("group not found")
		}
		return nil, fmt.Errorf("failed to get group details: %w", err)
	}

	return &group, nil
}
