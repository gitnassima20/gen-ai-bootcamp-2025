package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// Word represents a Japanese vocabulary word
type Word struct {
	ID      int64           `json:"id"`
	Kanji   string          `json:"kanji"`
	Romaji  string          `json:"romaji"`
	English string          `json:"english"`
	Parts   json.RawMessage `json:"parts"`
}

// UnmarshalParts attempts to parse the parts column safely
func (w *Word) UnmarshalParts(data interface{}) error {
	switch v := data.(type) {
	case string:
		// If it's a string, try to convert to RawMessage
		return json.Unmarshal([]byte(v), &w.Parts)
	case []byte:
		// If it's already bytes, use directly
		w.Parts = json.RawMessage(v)
		return nil
	case nil:
		// If nil, set to empty JSON object
		w.Parts = json.RawMessage("{}")
		return nil
	default:
		return fmt.Errorf("unsupported type for parts: %T", data)
	}
}

// StudyActivity represents different types of study activities
type StudyActivity struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	URL    string `json:"thumbnail_url"`
}

// StudySession represents an individual study session
type StudySession struct {
	ID              int64     `json:"id"`
	GroupID         int64     `json:"group_id"`
	StudyActivityID int64     `json:"study_activity_id"`
	CreatedAt       time.Time `json:"created_at"`
}

// WordReviewItem represents a review of a word during a study session
type WordReviewItem struct {
	ID             int64     `json:"id"`
	WordID         int64     `json:"word_id"`
	StudySessionID int64     `json:"study_session_id"`
	Correct        bool      `json:"correct"`
	CreatedAt      time.Time `json:"created_at"`
}
