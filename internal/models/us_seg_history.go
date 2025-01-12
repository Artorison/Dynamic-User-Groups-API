package models

import "time"

type OperationType string

const (
	ADD    OperationType = "ADD"
	DELETE OperationType = "DELETE"
)

type UserHistory struct {
	UserID int64         `json:"user_id"`
	Date   OperationType `json:"operation_type"`
}

type UserSegmentsHistory struct {
	ID            int64         `json:"id"`
	UserID        int64         `json:"user_id"`
	SegmentSlug   Slug          `json:"segment_slug"`
	OperationType OperationType `json:"operation_type"`
	OperationDate time.Time     `json:"operation_date"`
}
