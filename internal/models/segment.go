package models

type Slug string

type Segments struct {
	ID   int64 `json:"id"`
	Slug Slug  `json:"slug"`
}

// SegmentRequest used to create segment
type SegmentRequest struct {
	Slug Slug `json:"slug" example:"DISCOUNT_30"` // segment name
}
