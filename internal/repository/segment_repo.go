package repository

import (
	"API/internal/models"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/lib/pq"
)

type SegmentRepository interface {
	CreateSegmentDB(slug models.Slug) error
	DeleteSegmentDB(slug models.Slug) error
	SelectAllSegmentsDB() ([]models.Segments, error)
	GetSegmentID(slugs []models.Slug) ([]int64, error)
	GetOneSegmentID(slug models.Slug) (int64, error)
}

type SegmentRepositoryDB struct {
	DB *sql.DB
}

func NewSegmentRepository(db *sql.DB) *SegmentRepositoryDB {
	return &SegmentRepositoryDB{DB: db}
}

func (r *SegmentRepositoryDB) CreateSegmentDB(slug models.Slug) error {
	const op = "internal/repository/CreateSegment"

	queryCheck := `SELECT id FROM segments WHERE slug = $1`
	var id int64
	err := r.DB.QueryRow(queryCheck, slug).Scan(&id)
	if err == nil {
		return fmt.Errorf("segment already exists with id %d", id)
	} else if err != sql.ErrNoRows {
		return err
	}

	query := `INSERT INTO segments (slug) VALUES ($1);`
	if _, err := r.DB.Exec(query, slug); err != nil {
		slog.String("op", op)
		return err
	}
	return nil
}

func (r *SegmentRepositoryDB) DeleteSegmentDB(slug models.Slug) error {
	const op = "internal/repository/DeleteSegment"
	query := `DELETE FROM segments WHERE slug = $1`

	if _, err := r.DB.Exec(query, slug); err != nil {
		slog.String("op", op)
		return err
	}
	return nil
}

func (r *SegmentRepositoryDB) SelectAllSegmentsDB() ([]models.Segments, error) {
	const op = "internal/repository/SelectAllSegments"
	query := `SELECT id, slug FROM segments`

	rows, err := r.DB.Query(query)
	if err != nil {
		slog.String("op", op)
		return nil, err
	}
	defer rows.Close()
	var segments []models.Segments
	for rows.Next() {
		var segment models.Segments
		if err := rows.Scan(&segment.ID, &segment.Slug); err != nil {
			slog.String("op", op)
			return nil, err
		}
		segments = append(segments, segment)
	}
	if err := rows.Err(); err != nil {
		slog.String("op", op)
		return nil, err
	}
	return segments, nil
}

func (r *SegmentRepositoryDB) GetSegmentID(slugs []models.Slug) ([]int64, error) {
	query := `SELECT id FROM segments WHERE slug = ANY($1);`

	rows, err := r.DB.Query(query, pq.Array(slugs))
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()
	var segmentIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan segment ID: %w", err)
		}
		segmentIDs = append(segmentIDs, id)
	}

	if len(segmentIDs) != len(slugs) {
		return nil, fmt.Errorf("some slugs do not exist in segments table")
	}

	return segmentIDs, nil
}

func (r *SegmentRepositoryDB) GetOneSegmentID(slug models.Slug) (int64, error) {
	const query = `SELECT id FROM segments WHERE slug = $1;`

	var slugID int64

	err := r.DB.QueryRow(query, slug).Scan(&slugID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("segment with slug '%s' not found", slug)
		}
		return 0, fmt.Errorf("failed to get segment ID: %w", err)
	}

	return slugID, nil
}
