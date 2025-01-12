package models

// Users represents a user entity in the system.
type Users struct {
	ID   int64  `json:"user_id"` // User's unique ID
	Name string `json:"name"`    // User's name
}
