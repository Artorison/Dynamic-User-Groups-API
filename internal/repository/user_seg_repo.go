package repository

import (
	"API/internal/models"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type UserSegmentRepository interface {
	GetUserSegmentsDВ(id int64) (models.UserSegments, error)
	GetAllUserSegmentsDB() ([]models.UserSegment, error)
	UpdateUserSegments(slugsToAdd []models.Slug, slugsToDelete []models.Slug, userID int64, ttl *time.Time) error
	DeleteUserSegment(userID int64, slug models.Slug) error
}

type UserSegmentRepositoryDB struct {
	DB                *sql.DB
	UserRepository    UserRepository
	SegmentRepository SegmentRepository
	HistoryRepository UserSegmentHistoryRepository
}

func NewUserSegmentRepository(db *sql.DB, userRepo UserRepository, segmentRepo SegmentRepository, historyRepo UserSegmentHistoryRepository) *UserSegmentRepositoryDB {
	return &UserSegmentRepositoryDB{
		DB:                db,
		UserRepository:    userRepo,
		SegmentRepository: segmentRepo,
		HistoryRepository: historyRepo,
	}
}

func (r *UserSegmentRepositoryDB) GetUserSegmentsDВ(id int64) (models.UserSegments, error) {
	query := `
		SELECT s.slug 
		FROM user_segments us
		JOIN segments s ON us.segment_id = s.id
		WHERE us.user_id = $1;
	`
	var segments models.UserSegments
	rows, err := r.DB.Query(query, id)
	if err != nil {
		return segments, err
	}
	defer rows.Close()
	segments.UserID = id

	for rows.Next() {
		var slug models.Slug
		if err := rows.Scan(&slug); err != nil {
			return segments, err
		}
		segments.Segments = append(segments.Segments, slug)
	}
	if err := rows.Err(); err != nil {
		return segments, err
	}

	return segments, nil
}

func (r *UserSegmentRepositoryDB) GetAllUserSegmentsDB() ([]models.UserSegment, error) {
	query := `
	SELECT us.user_id, segments.slug
	FROM user_segments us
	JOIN segments ON us.segment_id = segments.id;`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uSegments []models.UserSegment
	for rows.Next() {
		var uSeg models.UserSegment
		if err := rows.Scan(&uSeg.UserID, &uSeg.Segments); err != nil {
			return nil, err
		}
		uSegments = append(uSegments, uSeg)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return uSegments, nil
}

func (r *UserSegmentRepositoryDB) addSegmentsToUser(tx *sql.Tx, userID int64, slugsID []int64, ttl *time.Time) error {
	if len(slugsID) == 0 {
		return nil
	}

	const query = `
	INSERT INTO user_segments (user_id, segment_id, ttl)
	SELECT $1, UNNEST($2::BIGINT[]), $3
	ON CONFLICT (user_id, segment_id) DO UPDATE
	SET ttl = EXCLUDED.ttl;`

	if _, err := tx.Exec(query, userID, pq.Array(slugsID), ttl); err != nil {
		return fmt.Errorf("failed to add segments to user %d: %w", userID, err)
	}

	return nil
}

func (r *UserSegmentRepositoryDB) removeSegmentsFromUser(tx *sql.Tx, userID int64, slugsID []int64) error {
	if len(slugsID) == 0 {
		return nil
	}
	const query = `
	DELETE FROM user_segments
	WHERE user_id = $1
	AND segment_id = ANY($2)
	`

	if _, err := tx.Exec(query, userID, pq.Array(slugsID)); err != nil {
		return fmt.Errorf("failed to remove segments from user %d: %w", userID, err)
	}

	return nil
}

func (r *UserSegmentRepositoryDB) UpdateUserSegments(slugsToAdd []models.Slug, slugsToDelete []models.Slug, userID int64, ttl *time.Time) error {

	isexists, err := r.UserRepository.CheckUserExists(userID)
	if err != nil {
		return err
	}
	if !isexists {
		return fmt.Errorf("user with id = %d is not exists", userID)
	}

	SAddID, _ := r.SegmentRepository.GetSegmentID(slugsToAdd)
	SDeleteID, _ := r.SegmentRepository.GetSegmentID(slugsToDelete)

	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := r.addSegmentsToUser(tx, userID, SAddID, ttl); err != nil {
		return fmt.Errorf("failed to get ID for segments to add: %w", err)
	}

	if err := r.removeSegmentsFromUser(tx, userID, SDeleteID); err != nil {
		return fmt.Errorf("failed to get ID for segments to delete: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	for _, slug := range slugsToAdd {
		var record models.UserSegmentsHistory
		record.UserID = userID
		record.SegmentSlug = slug
		record.OperationType = models.ADD
		record.OperationDate = time.Now()
		r.HistoryRepository.SaveHistoryEntry(record)
	}

	for _, slug := range slugsToDelete {
		var record models.UserSegmentsHistory
		record.UserID = userID
		record.SegmentSlug = slug
		record.OperationType = models.DELETE
		record.OperationDate = time.Now()
		r.HistoryRepository.SaveHistoryEntry(record)
	}

	return nil
}

func (r *UserSegmentRepositoryDB) DeleteUserSegment(userID int64, slug models.Slug) error {
	isexists, err := r.UserRepository.CheckUserExists(userID)
	if err != nil {
		return err
	}
	if !isexists {
		return fmt.Errorf("user with id = %d is not exists", userID)
	}

	slugID, err := r.SegmentRepository.GetOneSegmentID(slug)
	if err != nil {
		return err
	}

	const query = `
	DELETE FROM user_segments
	WHERE user_id = $1 AND segment_id = $2;
	`

	if _, err = r.DB.Exec(query, userID, slugID); err != nil {
		return fmt.Errorf("failed to delete user segment (user_id: %d, segment_id: %d): %w", userID, slugID, err)
	}

	return nil
}
