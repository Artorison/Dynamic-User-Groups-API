package models

// UserSegment represents a user and their associated segments.
// @description Model representing a user and their associated segments.
type UserSegments struct {
	UserID   int64  `json:"user_id"`  // User's unique ID
	Segments []Slug `json:"segments"` // List of associated segments
}

// UserSegment represents a single user-segment relationship.
// @description Model representing a relationship between a user and a segment.
type UserSegment struct {
	UserID   int64 `json:"user_id"`  // User's unique ID
	Segments Slug  `json:"segments"` // Associated segment
}

// UpdateSegmentsRequest is used for updating a user's segments.
// @description Request payload for updating a user's associated segments.
type UpdateSegmentsRequest struct {
	AddSegments    []Slug  `json:"add_segments" example:"[\"VOICE_MESSAGES\"]"`  // Segments to add
	DeleteSegments []Slug  `json:"delete_segments" example:"[\"CHAT_SUPPORT\"]"` // Segments to delete
	UserID         int64   `json:"user_id" example:"123"`                        // User's unique ID
	TTL            *string `json:"ttl"`                                          // TTL default NULL
}
