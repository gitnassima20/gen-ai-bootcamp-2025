package models

import (
	"encoding/json"
)

// Group represents a collection of words
type Group struct {
	ID         int64  `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	WordsCount int    `json:"word_count" db:"words_count"`
}

// GroupWord represents a word within a group
type GroupWord struct {
	ID           int64           `json:"id"`
	Kanji        string          `json:"kanji"`
	Romaji       string          `json:"romaji"`
	English      string          `json:"english"`
	CorrectCount int             `json:"correct_count"`
	WrongCount   int             `json:"wrong_count"`
	Parts        json.RawMessage `json:"parts"`
}

// GroupWordsResponse represents the response for group words
type GroupWordsResponse struct {
	Words       []GroupWord `json:"words"`
	TotalPages  int         `json:"total_pages"`
	CurrentPage int         `json:"current_page"`
}

// GroupsResponse represents the response for groups list
type GroupsResponse struct {
	Groups      []Group `json:"groups"`
	TotalPages  int     `json:"total_pages"`
	CurrentPage int     `json:"current_page"`
}

// GroupQueryParams represents the query parameters for group-related requests
type GroupQueryParams struct {
	Page    int
	PerPage int
	SortBy  string
	Order   string
}

// DefaultGroupQueryParams returns the default query parameters
func DefaultGroupQueryParams() GroupQueryParams {
	return GroupQueryParams{
		Page:    1,
		PerPage: 10,
		SortBy:  "name",
		Order:   "asc",
	}
}

// RawGroupWordsResponse represents the raw words in a group
type RawGroupWordsResponse struct {
	GroupID   int64     `json:"group_id"`
	GroupName string    `json:"group_name"`
	Words     []RawWord `json:"words"`
}

// RawWord represents a word with full details
type RawWord struct {
	ID      int64           `json:"id"`
	Kanji   string          `json:"kanji"`
	Romaji  string          `json:"romaji"`
	English string          `json:"english"`
	Parts   json.RawMessage `json:"parts"`
}
