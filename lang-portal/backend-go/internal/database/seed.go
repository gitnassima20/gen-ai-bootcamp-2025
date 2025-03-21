package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"lang-portal/internal/models"
)

// Word represents the structure of a word for seeding
type Word struct {
	Kanji   string          `json:"kanji"`
	Romaji  string          `json:"romaji"`
	English string          `json:"english"`
	Parts   json.RawMessage `json:"parts"`
}

// SeedDatabase populates the database with initial data
func SeedDatabase(db *sql.DB) error {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Begin a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Seed words
	wordIDs, err := seedWords(tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to seed words: %w", err)
	}

	// Seed groups
	groupIDs, err := seedGroups(tx, wordIDs)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to seed groups: %w", err)
	}

	// Seed word-groups associations
	err = seedWordGroups(tx, wordIDs, groupIDs)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to seed word-groups associations: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func seedWords(tx *sql.Tx) ([]int64, error) {
	// Word files to seed
	wordFiles := []string{
		filepath.Join("c:\\Users\\nassima\\Desktop\\gen-ai-bootcamp-2025\\lang-portal\\backend-go", "seed", "words_adjectives.json"),
		filepath.Join("c:\\Users\\nassima\\Desktop\\gen-ai-bootcamp-2025\\lang-portal\\backend-go", "seed", "words_verbs.json"),
	}

	// Prepare word insert statement
	wordStmt, err := tx.Prepare(`
		INSERT INTO words (kanji, romaji, english, parts)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare word insert statement: %w", err)
	}
	defer wordStmt.Close()

	// Insert words and track their IDs
	var wordIDs []int64
	for _, seedPath := range wordFiles {
		// Read seed data
		data, err := os.ReadFile(seedPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read seed file %s: %w", seedPath, err)
		}

		var words []Word
		if err := json.Unmarshal(data, &words); err != nil {
			return nil, fmt.Errorf("failed to parse seed data from %s: %w", seedPath, err)
		}

		for _, word := range words {
			// Convert parts to JSON string
			partsJSON, err := json.Marshal(word.Parts)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal parts: %w", err)
			}

			// Insert word
			result, err := wordStmt.Exec(word.Kanji, word.Romaji, word.English, string(partsJSON))
			if err != nil {
				return nil, fmt.Errorf("failed to insert word: %w", err)
			}

			// Get word ID
			wordID, err := result.LastInsertId()
			if err != nil {
				return nil, fmt.Errorf("failed to get last insert ID: %w", err)
			}

			wordIDs = append(wordIDs, wordID)
		}
	}

	return wordIDs, nil
}

func seedGroups(tx *sql.Tx, wordIDs []int64) ([]int64, error) {
	// Read seed data
	seedPath := filepath.Join("c:\\Users\\nassima\\Desktop\\gen-ai-bootcamp-2025\\lang-portal\\backend-go", "seed", "groups.json")
	data, err := os.ReadFile(seedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read seed file: %w", err)
	}

	var groups []models.Group
	if err := json.Unmarshal(data, &groups); err != nil {
		return nil, fmt.Errorf("failed to parse seed data: %w", err)
	}

	// Prepare group insert statement
	groupStmt, err := tx.Prepare(`
		INSERT INTO groups (name, words_count)
		VALUES (?, ?)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare group insert statement: %w", err)
	}
	defer groupStmt.Close()

	// Prepare word_groups insert statement
	wordGroupStmt, err := tx.Prepare(`
		INSERT INTO word_groups (word_id, group_id)
		VALUES (?, ?)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare word_groups insert statement: %w", err)
	}
	defer wordGroupStmt.Close()

	// Insert groups and track their IDs
	var groupIDs []int64
	for _, group := range groups {
		// Insert group
		result, err := groupStmt.Exec(group.Name, group.WordsCount)
		if err != nil {
			return nil, fmt.Errorf("failed to insert group: %w", err)
		}

		// Get group ID
		groupID, err := result.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("failed to get last insert ID: %w", err)
		}

		// Add words to group (use a subset of words)
		maxWords := min(len(wordIDs), group.WordsCount)
		wordsToAdd := wordIDs[:maxWords]
		for _, wordID := range wordsToAdd {
			_, err = wordGroupStmt.Exec(wordID, groupID)
			if err != nil {
				return nil, fmt.Errorf("failed to add word to group: %w", err)
			}
		}

		groupIDs = append(groupIDs, groupID)
	}

	return groupIDs, nil
}

func seedWordGroups(tx *sql.Tx, wordIDs []int64, groupIDs []int64) error {
	// Prepare word_groups insert statement
	wordGroupStmt, err := tx.Prepare(`
		INSERT INTO word_groups (word_id, group_id)
		VALUES (?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare word_groups insert statement: %w", err)
	}
	defer wordGroupStmt.Close()

	// Instead of using fixed IDs from a JSON file, create associations dynamically
	// This is a simplified example - adjust according to your needs
	for i, wordID := range wordIDs {
		// Assign each word to at least one group
		groupIndex := i % len(groupIDs)
		groupID := groupIDs[groupIndex]

		_, err = wordGroupStmt.Exec(wordID, groupID)
		if err != nil {
			return fmt.Errorf("failed to insert word_group association: %w", err)
		}
	}

	return nil
}

// Helper function to get the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
