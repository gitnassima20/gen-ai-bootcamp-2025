package models

type Word struct {
	ID          int64  `json:"id"`
	Word        string `json:"word"`
	Translation string `json:"translation"`
	// TODO: Add other fields as needed
}
