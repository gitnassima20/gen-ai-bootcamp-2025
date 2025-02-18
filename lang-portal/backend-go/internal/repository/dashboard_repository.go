package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// LastStudySession represents the most recent study session
type LastStudySession struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	GroupID    int64     `json:"group_id"`
	GroupName  string    `json:"group_name"`
}

// StudyProgress represents the overall study progress
type StudyProgress struct {
	TotalWordsStudied     int `json:"total_words_studied"`
	TotalAvailableWords   int `json:"total_available_words"`
}

// QuickStats represents quick dashboard statistics
type QuickStats struct {
	SuccessRate         int `json:"success_rate"`
	TotalStudySessions  int `json:"total_study_sessions"`
	TotalActiveGroups   int `json:"total_active_groups"`
	CurrentStreak       int `json:"current_streak"`
}

// DashboardRepository defines methods for retrieving dashboard-related data
type DashboardRepository interface {
	// GetLastStudySession retrieves the most recent study session
	GetLastStudySession(ctx context.Context) (*LastStudySession, error)

	// GetStudyProgress calculates the overall study progress
	GetStudyProgress(ctx context.Context) (*StudyProgress, error)

	// GetQuickStats retrieves quick dashboard statistics
	GetQuickStats(ctx context.Context) (*QuickStats, error)
}

// SQLDashboardRepository implements DashboardRepository using SQLite
type SQLDashboardRepository struct {
	db *sql.DB
}

// NewDashboardRepository creates a new instance of SQLDashboardRepository
func NewDashboardRepository(db *sql.DB) *SQLDashboardRepository {
	return &SQLDashboardRepository{db: db}
}

// GetLastStudySession retrieves the most recent study session
func (r *SQLDashboardRepository) GetLastStudySession(ctx context.Context) (*LastStudySession, error) {
	query := `
		SELECT 
			ss.id, 
			sa.name as activity_name, 
			ss.created_at,
			g.id as group_id,
			g.name as group_name
		FROM study_sessions ss
		JOIN study_activities sa ON ss.study_activity_id = sa.id
		JOIN groups g ON ss.group_id = g.id
		ORDER BY ss.created_at DESC
		LIMIT 1
	`
	var session LastStudySession
	err := r.db.QueryRowContext(ctx, query).Scan(
		&session.ID,
		&session.Name,
		&session.CreatedAt,
		&session.GroupID,
		&session.GroupName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no study sessions found")
		}
		return nil, fmt.Errorf("failed to retrieve last study session: %w", err)
	}

	return &session, nil
}

// GetStudyProgress calculates the overall study progress
func (r *SQLDashboardRepository) GetStudyProgress(ctx context.Context) (*StudyProgress, error) {
	// Count total words studied (reviewed at least once)
	studiedWordsQuery := `
		SELECT COUNT(DISTINCT word_id) 
		FROM word_review_items
	`
	var totalWordsStudied int
	err := r.db.QueryRowContext(ctx, studiedWordsQuery).Scan(&totalWordsStudied)
	if err != nil {
		return nil, fmt.Errorf("failed to count studied words: %w", err)
	}

	// Count total available words
	availableWordsQuery := `
		SELECT COUNT(*) 
		FROM words
	`
	var totalAvailableWords int
	err = r.db.QueryRowContext(ctx, availableWordsQuery).Scan(&totalAvailableWords)
	if err != nil {
		return nil, fmt.Errorf("failed to count available words: %w", err)
	}

	return &StudyProgress{
		TotalWordsStudied:   totalWordsStudied,
		TotalAvailableWords: totalAvailableWords,
	}, nil
}

// GetQuickStats retrieves quick dashboard statistics
func (r *SQLDashboardRepository) GetQuickStats(ctx context.Context) (*QuickStats, error) {
	// Calculate success rate
	successRateQuery := `
		SELECT 
			CASE 
				WHEN COUNT(*) > 0 
				THEN ROUND(100.0 * SUM(CASE WHEN correct = 1 THEN 1 ELSE 0 END) / COUNT(*), 0)
				ELSE 0 
			END as success_rate
		FROM word_review_items
	`
	var successRate int
	err := r.db.QueryRowContext(ctx, successRateQuery).Scan(&successRate)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate success rate: %w", err)
	}

	// Count total study sessions
	sessionsQuery := `
		SELECT COUNT(*) 
		FROM study_sessions
	`
	var totalStudySessions int
	err = r.db.QueryRowContext(ctx, sessionsQuery).Scan(&totalStudySessions)
	if err != nil {
		return nil, fmt.Errorf("failed to count study sessions: %w", err)
	}

	// Count active groups
	activeGroupsQuery := `
		SELECT COUNT(*) 
		FROM groups 
		WHERE words_count > 0
	`
	var totalActiveGroups int
	err = r.db.QueryRowContext(ctx, activeGroupsQuery).Scan(&totalActiveGroups)
	if err != nil {
		return nil, fmt.Errorf("failed to count active groups: %w", err)
	}

	// Calculate current streak (consecutive days with study sessions)
	// Note: This is a simplified implementation and might need more complex logic
	currentStreakQuery := `
		WITH study_days AS (
			SELECT DISTINCT date(created_at) as study_date
			FROM study_sessions
			ORDER BY study_date DESC
		), consecutive_days AS (
			SELECT 
				study_date, 
				julianday('now') - julianday(study_date) as days_ago,
				ROW_NUMBER() OVER (ORDER BY study_date DESC) as rn
			FROM study_days
		)
		SELECT COUNT(*) 
		FROM consecutive_days 
		WHERE days_ago = rn - 1
	`
	var currentStreak int
	err = r.db.QueryRowContext(ctx, currentStreakQuery).Scan(&currentStreak)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate current streak: %w", err)
	}

	return &QuickStats{
		SuccessRate:        successRate,
		TotalStudySessions: totalStudySessions,
		TotalActiveGroups:  totalActiveGroups,
		CurrentStreak:      currentStreak,
	}, nil
}
