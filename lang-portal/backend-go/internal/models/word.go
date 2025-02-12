package models

import (
	"encoding/json"
	"time"
)

// Word represents a Japanese vocabulary word
type Word struct {
	ID      int64           `json:"id"`
	Kanji   string         `json:"kanji"`
	Romaji  string         `json:"romaji"`
	English string         `json:"english"`
	Parts   json.RawMessage `json:"parts"`
}

// WordGroup represents the many-to-many relationship between words and groups
type WordGroup struct {
	WordID  int64 `json:"word_id"`
	GroupID int64 `json:"group_id"`
}

// Group represents a collection of words
type Group struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	WordsCount int    `json:"words_count"`
}

// StudyActivity represents different types of study activities
type StudyActivity struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
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
