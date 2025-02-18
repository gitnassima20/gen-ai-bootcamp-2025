package repository

import (
	"context"
	"database/sql"
	"fmt"
	"lang-portal/internal/models"
	"strings"
)

// WordRepository defines the interface for word-related database operations
type WordRepository interface {
	// Create adds a new word to the database
	Create(ctx context.Context, word *models.Word) error

	// GetByID retrieves a word by its ID
	GetByID(ctx context.Context, id int64) (*models.Word, error)

	// Update modifies an existing word
	Update(ctx context.Context, word *models.Word) error

	// Delete removes a word from the database
	Delete(ctx context.Context, id int64) error

	// List retrieves words with optional filtering and pagination
	List(ctx context.Context, filter WordFilter, page, pageSize int) ([]models.Word, int, error)

	// AddToGroup adds a word to a group
	AddToGroup(ctx context.Context, wordID, groupID int64) error

	// RemoveFromGroup removes a word from a group
	RemoveFromGroup(ctx context.Context, wordID, groupID int64) error
}

// WordFilter represents filtering options for word queries
type WordFilter struct {
	Kanji   string
	Romaji  string
	English string
	GroupID int64
}

// SQLWordRepository implements WordRepository using SQLite
type SQLWordRepository struct {
	db *sql.DB
}

// NewWordRepository creates a new instance of SQLWordRepository
func NewWordRepository(db *sql.DB) *SQLWordRepository {
	return &SQLWordRepository{db: db}
}

// Create adds a new word to the database
func (r *SQLWordRepository) Create(ctx context.Context, word *models.Word) error {
	query := `
		INSERT INTO words (kanji, romaji, english, parts)
		VALUES (?, ?, ?, ?)
	`
	result, err := r.db.ExecContext(ctx, query, word.Kanji, word.Romaji, word.English, word.Parts)
	if err != nil {
		return fmt.Errorf("failed to create word: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}
	word.ID = id
	return nil
}

// GetByID retrieves a word by its ID
func (r *SQLWordRepository) GetByID(ctx context.Context, id int64) (*models.Word, error) {
	query := `
		SELECT id, kanji, romaji, english, parts
		FROM words
		WHERE id = ?
	`
	var word models.Word
	var partsData interface{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&word.ID, &word.Kanji, &word.Romaji, &word.English, &partsData,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get word: %w", err)
	}

	// Use the new UnmarshalParts method
	if err := word.UnmarshalParts(partsData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal parts: %w", err)
	}

	return &word, nil
}

// Update modifies an existing word
func (r *SQLWordRepository) Update(ctx context.Context, word *models.Word) error {
	query := `
		UPDATE words
		SET kanji = ?, romaji = ?, english = ?, parts = ?
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query,
		word.Kanji, word.Romaji, word.English, word.Parts, word.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update word: %w", err)
	}
	return nil
}

// Delete removes a word from the database
func (r *SQLWordRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM words WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete word: %w", err)
	}
	return nil
}

// List retrieves words with optional filtering and pagination
func (r *SQLWordRepository) List(ctx context.Context, filter WordFilter, page, pageSize int) ([]models.Word, int, error) {
	// Build dynamic query based on filter
	baseQuery := `SELECT id, kanji, romaji, english, parts FROM words`
	countQuery := `SELECT COUNT(*) FROM words`
	var conditions []string
	var args []interface{}

	if filter.Kanji != "" {
		conditions = append(conditions, "kanji LIKE ?")
		args = append(args, "%"+filter.Kanji+"%")
	}
	if filter.Romaji != "" {
		conditions = append(conditions, "romaji LIKE ?")
		args = append(args, "%"+filter.Romaji+"%")
	}
	if filter.English != "" {
		conditions = append(conditions, "english LIKE ?")
		args = append(args, "%"+filter.English+"%")
	}
	if filter.GroupID > 0 {
		baseQuery += ` JOIN word_groups wg ON words.id = wg.word_id`
		countQuery += ` JOIN word_groups wg ON words.id = wg.word_id`
		conditions = append(conditions, "wg.group_id = ?")
		args = append(args, filter.GroupID)
	}

	// Add WHERE clause if conditions exist
	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
		countQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Add pagination
	offset := (page - 1) * pageSize
	baseQuery += ` LIMIT ? OFFSET ?`
	args = append(args, pageSize, offset)

	// Count total records
	var totalCount int
	err := r.db.QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count words: %w", err)
	}

	// Fetch words
	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list words: %w", err)
	}
	defer rows.Close()

	var words []models.Word
	for rows.Next() {
		var word models.Word
		var partsData interface{}
		if err := rows.Scan(
			&word.ID,
			&word.Kanji,
			&word.Romaji,
			&word.English,
			&partsData,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan word: %w", err)
		}

		// Use the new UnmarshalParts method
		if err := word.UnmarshalParts(partsData); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal parts: %w", err)
		}

		words = append(words, word)
	}

	return words, totalCount, nil
}

// AddToGroup adds a word to a group
func (r *SQLWordRepository) AddToGroup(ctx context.Context, wordID, groupID int64) error {
	query := `
		INSERT OR IGNORE INTO word_groups (word_id, group_id)
		VALUES (?, ?)
	`
	_, err := r.db.ExecContext(ctx, query, wordID, groupID)
	if err != nil {
		return fmt.Errorf("failed to add word to group: %w", err)
	}
	return nil
}

// RemoveFromGroup removes a word from a group
func (r *SQLWordRepository) RemoveFromGroup(ctx context.Context, wordID, groupID int64) error {
	query := `
		DELETE FROM word_groups 
		WHERE word_id = ? AND group_id = ?
	`
	_, err := r.db.ExecContext(ctx, query, wordID, groupID)
	if err != nil {
		return fmt.Errorf("failed to remove word from group: %w", err)
	}
	return nil
}
