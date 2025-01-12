package repository

import (
	"API/internal/models"
	"database/sql"
	"fmt"
	"time"
)

type UserSegmentHistoryRepositoryDB struct {
	DB *sql.DB
}

func NewUserSegmentHistoryRepository(db *sql.DB) *UserSegmentHistoryRepositoryDB {
	return &UserSegmentHistoryRepositoryDB{DB: db}
}

type UserSegmentHistoryRepository interface {
	SaveHistoryEntry(record models.UserSegmentsHistory) error
	GetUserHistory(userID int64, start, end time.Time) ([]models.UserSegmentsHistory, error)
}

func (r *UserSegmentHistoryRepositoryDB) SaveHistoryEntry(record models.UserSegmentsHistory) error {
	query := `
		INSERT INTO user_segments_history (user_id, segment_slug, operation_type, operation_date)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.DB.Exec(query, record.UserID, record.SegmentSlug, record.OperationType, record.OperationDate)
	if err != nil {
		return fmt.Errorf("failed to save history entry: %w", err)
	}

	return nil
}

func (r *UserSegmentHistoryRepositoryDB) GetUserHistory(userID int64, start, end time.Time) ([]models.UserSegmentsHistory, error) {

	query := `
	SELECT id, user_id, segment_slug, operation_type, operation_date
	FROM user_segments_history
	WHERE user_id = $1
	AND operation_date BETWEEN $2 AND $3;
	`

	rows, err := r.DB.Query(query, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	histories := make([]models.UserSegmentsHistory, 0)

	for rows.Next() {
		var history models.UserSegmentsHistory
		if err := rows.Scan(
			&history.ID,
			&history.UserID,
			&history.SegmentSlug,
			&history.OperationType,
			&history.OperationDate,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		histories = append(histories, history)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return histories, nil
}
