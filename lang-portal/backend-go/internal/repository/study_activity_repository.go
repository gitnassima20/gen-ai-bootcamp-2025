package repository

import (
	"context"
	"database/sql"
	"fmt"
	"lang-portal/internal/models"
)

// StudyActivityRepository defines the interface for study activity-related database operations
type StudyActivityRepository interface {
	// List retrieves all study activities with optional pagination
	List(ctx context.Context, page, pageSize int) ([]models.StudyActivity, int, error)

	// GetByID retrieves a specific study activity by its ID
	GetByID(ctx context.Context, id int64) (*models.StudyActivity, error)

	// GetActivityDetails retrieves additional details for a study activity
	GetActivityDetails(ctx context.Context, id int64) (*StudyActivityDetails, error)
}

// StudyActivityDetails contains additional information about a study activity
type StudyActivityDetails struct {
	TotalSessions int `json:"total_sessions"`
}

// SQLStudyActivityRepository implements StudyActivityRepository using SQLite
type SQLStudyActivityRepository struct {
	db *sql.DB
}

// NewStudyActivityRepository creates a new instance of SQLStudyActivityRepository
func NewStudyActivityRepository(db *sql.DB) *SQLStudyActivityRepository {
	return &SQLStudyActivityRepository{db: db}
}

// List retrieves study activities with pagination
func (r *SQLStudyActivityRepository) List(ctx context.Context, page, pageSize int) ([]models.StudyActivity, int, error) {
	// Count total activities
	countQuery := `SELECT COUNT(*) FROM study_activities`
	var totalCount int
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count study activities: %w", err)
	}

	// Calculate pagination
	offset := (page - 1) * pageSize

	// Fetch activities
	query := `
		SELECT 
			sa.id, 
			sa.name, 
			sa.url
		FROM study_activities sa
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list study activities: %w", err)
	}
	defer rows.Close()

	var activities []models.StudyActivity
	for rows.Next() {
		var activity models.StudyActivity
		if err := rows.Scan(
			&activity.ID,
			&activity.Name,
			&activity.URL,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan study activity: %w", err)
		}

		activities = append(activities, activity)
	}

	return activities, totalCount, nil
}

// GetByID retrieves a specific study activity by its ID
func (r *SQLStudyActivityRepository) GetByID(ctx context.Context, id int64) (*models.StudyActivity, error) {
	query := `
		SELECT 
			sa.id, 
			sa.name, 
			sa.url
		FROM study_activities sa
		WHERE sa.id = ?
	`
	var activity models.StudyActivity
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&activity.ID,
		&activity.Name,
		&activity.URL,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("study activity not found")
		}
		return nil, fmt.Errorf("failed to get study activity: %w", err)
	}

	return &activity, nil
}

// GetActivityDetails retrieves additional details for a study activity
func (r *SQLStudyActivityRepository) GetActivityDetails(ctx context.Context, id int64) (*StudyActivityDetails, error) {
	query := `
		SELECT 
			COUNT(ss.id) as total_sessions
		FROM study_activities sa
		LEFT JOIN study_sessions ss ON ss.study_activity_id = sa.id
		WHERE sa.id = ?
		GROUP BY sa.id
	`
	var details StudyActivityDetails
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&details.TotalSessions,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no study activity found with ID: %d", id)
		}
		return nil, fmt.Errorf("failed to get study activity details: %w", err)
	}

	return &details, nil
}
