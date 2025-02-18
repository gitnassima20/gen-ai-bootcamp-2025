package repository

import (
	"context"
	"database/sql"
	"fmt"
	"lang-portal/internal/models"
	"strings"
	"time"
)

// StudySessionListItem represents a study session in the list view
type StudySessionListItem struct {
	ID                int64     `json:"id"`
	ActivityName      string    `json:"activity_name"`
	GroupName         string    `json:"group_name"`
	StartTime         time.Time `json:"start_time"`
	EndTime           time.Time `json:"end_time"`
	TotalWordsReviewed int      `json:"total_words_reviewed"`
}

// StudySessionRepository defines the interface for study session-related database operations
type StudySessionRepository interface {
	// List retrieves study sessions with optional filtering and pagination
	List(ctx context.Context, studyActivityID, groupID int64, page, pageSize int) ([]StudySessionListItem, int, error)

	// Create adds a new study session
	Create(ctx context.Context, session *models.StudySession) error

	// CreateWordReview adds a new word review item to a study session
	CreateWordReview(ctx context.Context, review *models.WordReviewItem) error

	// GetWordReviewsBySessionID retrieves all word reviews for a specific study session
	GetWordReviewsBySessionID(ctx context.Context, studySessionID int64) ([]models.WordReviewItem, error)

	// ListWordsByStudySession retrieves words studied in a specific session with performance statistics
	ListWordsByStudySession(ctx context.Context, sessionID int64, page, wordsPerPage int) ([]WordStats, int, error)

	// GetStudySessionDetails retrieves detailed information about a specific study session
	GetStudySessionDetails(ctx context.Context, sessionID int64) (*StudySessionDetails, error)
}

// SQLStudySessionRepository implements StudySessionRepository using SQLite
type SQLStudySessionRepository struct {
	db *sql.DB
}

// NewStudySessionRepository creates a new instance of SQLStudySessionRepository
func NewStudySessionRepository(db *sql.DB) *SQLStudySessionRepository {
	return &SQLStudySessionRepository{db: db}
}

// List retrieves study sessions with optional filtering and pagination
func (r *SQLStudySessionRepository) List(ctx context.Context, studyActivityID, groupID int64, page, pageSize int) ([]StudySessionListItem, int, error) {
	// Prepare filter conditions
	conditions := []string{}
	args := []interface{}{}

	if studyActivityID > 0 {
		conditions = append(conditions, "ss.study_activity_id = ?")
		args = append(args, studyActivityID)
	}

	if groupID > 0 {
		conditions = append(conditions, "ss.group_id = ?")
		args = append(args, groupID)
	}

	// Construct where clause
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total sessions
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM study_sessions ss %s", whereClause)
	var totalSessions int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalSessions)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count study sessions: %w", err)
	}

	// Calculate pagination
	offset := (page - 1) * pageSize

	// Query to fetch detailed session information
	query := fmt.Sprintf(`
		SELECT 
			ss.id,
			sa.name as activity_name,
			g.name as group_name,
			ss.created_at as start_time,
			ss.created_at as end_time,
			(SELECT COUNT(*) FROM word_review_items wri WHERE wri.study_session_id = ss.id) as total_words_reviewed
		FROM study_sessions ss
		JOIN study_activities sa ON ss.study_activity_id = sa.id
		JOIN groups g ON ss.group_id = g.id
		%s
		ORDER BY ss.created_at DESC
		LIMIT ? OFFSET ?
	`, whereClause)
	
	// Add pagination args
	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list study sessions: %w", err)
	}
	defer rows.Close()

	var sessions []StudySessionListItem
	for rows.Next() {
		var session StudySessionListItem
		if err := rows.Scan(
			&session.ID,
			&session.ActivityName,
			&session.GroupName,
			&session.StartTime,
			&session.EndTime,
			&session.TotalWordsReviewed,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan study session: %w", err)
		}
		sessions = append(sessions, session)
	}

	return sessions, totalSessions, nil
}

// Create adds a new study session
func (r *SQLStudySessionRepository) Create(ctx context.Context, session *models.StudySession) error {
	query := `
		INSERT INTO study_sessions 
		(group_id, study_activity_id, created_at) 
		VALUES (?, ?, ?)
	`
	result, err := r.db.ExecContext(ctx, query,
		session.GroupID,
		session.StudyActivityID,
		session.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create study session: %w", err)
	}

	// Set the ID of the newly created session
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}
	session.ID = id

	return nil
}

// CreateWordReview adds a new word review item to a study session
func (r *SQLStudySessionRepository) CreateWordReview(ctx context.Context, review *models.WordReviewItem) error {
	// Validate that the study session exists
	sessionQuery := `SELECT 1 FROM study_sessions WHERE id = ?`
	var exists int
	err := r.db.QueryRowContext(ctx, sessionQuery, review.StudySessionID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("study session not found")
		}
		return fmt.Errorf("failed to validate study session: %w", err)
	}

	// Validate that the word exists
	wordQuery := `SELECT 1 FROM words WHERE id = ?`
	err = r.db.QueryRowContext(ctx, wordQuery, review.WordID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("word not found")
		}
		return fmt.Errorf("failed to validate word: %w", err)
	}

	// Insert the word review
	query := `
		INSERT INTO word_review_items 
		(word_id, study_session_id, correct, created_at) 
		VALUES (?, ?, ?, ?)
	`
	result, err := r.db.ExecContext(ctx, query,
		review.WordID,
		review.StudySessionID,
		review.Correct,
		review.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create word review: %w", err)
	}

	// Set the ID of the newly created review
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}
	review.ID = id

	return nil
}

// GetWordReviewsBySessionID retrieves all word reviews for a specific study session
func (r *SQLStudySessionRepository) GetWordReviewsBySessionID(ctx context.Context, studySessionID int64) ([]models.WordReviewItem, error) {
	query := `
		SELECT 
			id, 
			word_id, 
			study_session_id, 
			correct, 
			created_at
		FROM word_review_items
		WHERE study_session_id = ?
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, studySessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch word reviews: %w", err)
	}
	defer rows.Close()

	var reviews []models.WordReviewItem
	for rows.Next() {
		var review models.WordReviewItem
		if err := rows.Scan(
			&review.ID,
			&review.WordID,
			&review.StudySessionID,
			&review.Correct,
			&review.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan word review: %w", err)
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}

// WordStats represents word performance statistics
type WordStats struct {
	ID           int64  `json:"id"`
	Kanji        string `json:"kanji"`
	Romaji       string `json:"romaji"`
	English      string `json:"english"`
	CorrectCount int    `json:"correct_count"`
	WrongCount   int    `json:"wrong_count"`
}

// ListWordsByStudySession retrieves words studied in a specific session with performance statistics
func (r *SQLStudySessionRepository) ListWordsByStudySession(ctx context.Context, sessionID int64, page, wordsPerPage int) ([]WordStats, int, error) {
	// Base query to count total words in the session
	countQuery := `
		SELECT COUNT(DISTINCT w.id)
		FROM words w
		JOIN word_review_items wri ON w.id = wri.word_id
		WHERE wri.study_session_id = ?
	`
	var totalWords int
	err := r.db.QueryRowContext(ctx, countQuery, sessionID).Scan(&totalWords)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count words: %w", err)
	}

	// Calculate pagination
	offset := (page - 1) * wordsPerPage

	// Query to fetch words with their review statistics
	query := `
		SELECT 
			w.id, 
			w.kanji, 
			w.romaji, 
			w.english,
			SUM(CASE WHEN wri.correct = 1 THEN 1 ELSE 0 END) as correct_count,
			SUM(CASE WHEN wri.correct = 0 THEN 1 ELSE 0 END) as wrong_count
		FROM words w
		JOIN word_review_items wri ON w.id = wri.word_id
		WHERE wri.study_session_id = ?
		GROUP BY w.id, w.kanji, w.romaji, w.english
		ORDER BY w.id
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, sessionID, wordsPerPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch words: %w", err)
	}
	defer rows.Close()

	var words []WordStats
	for rows.Next() {
		var word WordStats
		if err := rows.Scan(
			&word.ID,
			&word.Kanji,
			&word.Romaji,
			&word.English,
			&word.CorrectCount,
			&word.WrongCount,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan word: %w", err)
		}
		words = append(words, word)
	}

	return words, totalWords, nil
}

// StudySessionDetails represents detailed information about a study session
type StudySessionDetails struct {
	ID                 int64     `json:"id"`
	ActivityName       string    `json:"activity_name"`
	GroupName          string    `json:"group_name"`
	StartTime          time.Time `json:"start_time"`
	TotalWordsReviewed int       `json:"total_words_reviewed"`
}

// GetStudySessionDetails retrieves detailed information about a specific study session
func (r *SQLStudySessionRepository) GetStudySessionDetails(ctx context.Context, sessionID int64) (*StudySessionDetails, error) {
	query := `
		SELECT 
			ss.id,
			sa.name as activity_name,
			g.name as group_name,
			ss.created_at as start_time,
			(SELECT COUNT(*) FROM word_review_items wri WHERE wri.study_session_id = ss.id) as total_words_reviewed
		FROM study_sessions ss
		JOIN study_activities sa ON ss.study_activity_id = sa.id
		JOIN groups g ON ss.group_id = g.id
		WHERE ss.id = ?
	`
	var details StudySessionDetails
	err := r.db.QueryRowContext(ctx, query, sessionID).Scan(
		&details.ID,
		&details.ActivityName,
		&details.GroupName,
		&details.StartTime,
		&details.TotalWordsReviewed,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("study session not found")
		}
		return nil, fmt.Errorf("failed to get study session details: %w", err)
	}

	return &details, nil
}
