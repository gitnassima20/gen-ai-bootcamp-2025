package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
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

// GroupWordItem represents a word in a group
type GroupWordItem struct {
	ID            int64  `json:"id"`
	Kanji         string `json:"kanji"`
	Romaji        string `json:"romaji"`
	English       string `json:"english"`
	CorrectCount  int    `json:"correct_count"`
	WrongCount    int    `json:"wrong_count"`
}

// GroupStudySessionItem represents a study session for a group
type GroupStudySessionItem struct {
	ID                 int64     `json:"id"`
	Name               string    `json:"name"`
	StartTime          time.Time `json:"start_time"`
	EndTime            time.Time `json:"end_time"`
	TotalWordsReviewed int       `json:"total_words_reviewed"`
}

// GroupRepository defines the interface for group-related database operations
type GroupRepository interface {
	// List retrieves groups with optional pagination
	List(ctx context.Context, page, groupsPerPage int) ([]GroupListItem, int, error)

	// GetByID retrieves detailed information about a specific group
	GetByID(ctx context.Context, groupID int64) (*GroupDetails, error)

	// GetGroupWords retrieves words in a group with pagination
	GetGroupWords(ctx context.Context, groupID int64, page, wordsPerPage int) ([]GroupWordItem, int, error)

	// GetGroupStudySessions retrieves study sessions for a group with pagination
	GetGroupStudySessions(ctx context.Context, groupID int64, page, sessionsPerPage int) ([]GroupStudySessionItem, int, error)
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

// GetGroupWords retrieves words in a group with pagination
func (r *SQLGroupRepository) GetGroupWords(ctx context.Context, groupID int64, page, wordsPerPage int) ([]GroupWordItem, int, error) {
	// First, verify the group exists
	_, err := r.GetByID(ctx, groupID)
	if err != nil {
		return nil, 0, err
	}

	// Count total words in the group
	countQuery := `
		SELECT COUNT(*) 
		FROM words w
		JOIN word_groups wg ON w.id = wg.word_id
		WHERE wg.group_id = ?
	`
	var totalWords int
	err = r.db.QueryRowContext(ctx, countQuery, groupID).Scan(&totalWords)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count group words: %w", err)
	}

	// Calculate pagination
	offset := (page - 1) * wordsPerPage

	// Query to fetch words with review statistics
	query := `
		SELECT 
			w.id, 
			w.kanji, 
			w.romaji, 
			w.english,
			COALESCE(SUM(CASE WHEN wri.correct = 1 THEN 1 ELSE 0 END), 0) as correct_count,
			COALESCE(SUM(CASE WHEN wri.correct = 0 THEN 1 ELSE 0 END), 0) as wrong_count
		FROM words w
		JOIN word_groups wg ON w.id = wg.word_id
		LEFT JOIN word_review_items wri ON w.id = wri.word_id
		WHERE wg.group_id = ?
		GROUP BY w.id, w.kanji, w.romaji, w.english
		ORDER BY w.id
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, groupID, wordsPerPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch group words: %w", err)
	}
	defer rows.Close()

	var words []GroupWordItem
	for rows.Next() {
		var word GroupWordItem
		if err := rows.Scan(
			&word.ID,
			&word.Kanji,
			&word.Romaji,
			&word.English,
			&word.CorrectCount,
			&word.WrongCount,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan group word: %w", err)
		}
		words = append(words, word)
	}

	return words, totalWords, nil
}

// GetGroupStudySessions retrieves study sessions for a group with pagination
func (r *SQLGroupRepository) GetGroupStudySessions(ctx context.Context, groupID int64, page, sessionsPerPage int) ([]GroupStudySessionItem, int, error) {
	// First, verify the group exists
	_, err := r.GetByID(ctx, groupID)
	if err != nil {
		return nil, 0, err
	}

	// Count total study sessions in the group
	countQuery := `
		SELECT COUNT(*) 
		FROM study_sessions ss
		WHERE ss.group_id = ?
	`
	var totalSessions int
	err = r.db.QueryRowContext(ctx, countQuery, groupID).Scan(&totalSessions)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count group study sessions: %w", err)
	}

	// Calculate pagination
	offset := (page - 1) * sessionsPerPage

	// Query to fetch study sessions with total words reviewed
	query := `
		SELECT 
			ss.id, 
			sa.name as session_name,
			ss.created_at as start_time,
			ss.created_at as end_time,
			(SELECT COUNT(*) FROM word_review_items wri WHERE wri.study_session_id = ss.id) as total_words_reviewed
		FROM study_sessions ss
		JOIN study_activities sa ON ss.study_activity_id = sa.id
		WHERE ss.group_id = ?
		ORDER BY ss.created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, groupID, sessionsPerPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch group study sessions: %w", err)
	}
	defer rows.Close()

	var sessions []GroupStudySessionItem
	for rows.Next() {
		var session GroupStudySessionItem
		if err := rows.Scan(
			&session.ID,
			&session.Name,
			&session.StartTime,
			&session.EndTime,
			&session.TotalWordsReviewed,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan group study session: %w", err)
		}
		sessions = append(sessions, session)
	}

	return sessions, totalSessions, nil
}
