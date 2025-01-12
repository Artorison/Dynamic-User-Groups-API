package services

import (
	"API/internal/models"
	"API/internal/repository"
	"fmt"
)

//go:generate mockery --name=ISegmentService --output=mocks --outpkg=mocks
type ISegmentService interface {
	GetAllSegments() ([]models.Segments, error)
	CreateSegment(slug models.Slug) error
	DeleteSegment(slug models.Slug) error
}

type SegmentService struct {
	Repo repository.SegmentRepository
}

func NewSegmentService(repo repository.SegmentRepository) *SegmentService {
	return &SegmentService{Repo: repo}
}

func (s *SegmentService) GetAllSegments() ([]models.Segments, error) {
	segments, err := s.Repo.SelectAllSegmentsDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get all segments: %w", err)
	}
	return segments, nil
}

func (s *SegmentService) CreateSegment(slug models.Slug) error {
	if err := s.Repo.CreateSegmentDB(slug); err != nil {
		return fmt.Errorf("failed to create segment: %w", err)
	}
	return nil
}

func (s *SegmentService) DeleteSegment(slug models.Slug) error {
	if err := s.Repo.DeleteSegmentDB(slug); err != nil {
		return fmt.Errorf("failed to delete segment: %w", err)
	}
	return nil
}
