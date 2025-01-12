package services

import (
	"API/internal/kafka"
	"API/internal/models"
	"API/internal/repository"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/IBM/sarama"
)

//go:generate mockery --name=IUserSegmentService --output=mocks --outpkg=mocks
type IUserSegmentService interface {
	GetUserSegments(userID int64) (models.UserSegments, error)
	GetAllUserSegments() ([]models.UserSegment, error)
	UpdateUserSegments(userID int64, slugsToAdd, slugsToDelete []models.Slug, ttl *time.Time) error
	DeleteUserSegment(userID int64, slug models.Slug) error
}

type UserSegmentService struct {
	Repo        repository.UserSegmentRepository
	HistoryRepo repository.UserSegmentHistoryRepository
	Producer    *kafka.Producer
}

func NewUserSegmentService(repo repository.UserSegmentRepository, historyRepo repository.UserSegmentHistoryRepository, producer *kafka.Producer) *UserSegmentService {
	return &UserSegmentService{
		Repo:        repo,
		HistoryRepo: historyRepo,
		Producer:    producer,
	}
}

func (s *UserSegmentService) GetUserSegments(userID int64) (models.UserSegments, error) {
	return s.Repo.GetUserSegmentsDВ(userID)
}

func (s *UserSegmentService) GetAllUserSegments() ([]models.UserSegment, error) {
	return s.Repo.GetAllUserSegmentsDB()
}

func (s *UserSegmentService) UpdateUserSegments(userID int64, slugsToAdd, slugsToDelete []models.Slug, ttl *time.Time) error {
	err := s.Repo.UpdateUserSegments(slugsToAdd, slugsToDelete, userID, ttl)
	if err != nil {
		return err
	}

	for _, slug := range slugsToAdd {
		event := map[string]interface{}{
			"user_id": userID,
			"segment": slug,
			"action":  "add",
		}
		if ttl != nil {
			event["ttl"] = ttl.Format(time.RFC3339)
			err := s.Producer.SendMessage("segment_expiry", strconv.FormatInt(userID, 10), event)
			if err != nil {
				log.Printf("Failed to send TTL Kafka message: %v", err)
				return err
			}
			log.Printf("Kafka message sent for TTL expiry: %v", event)
		}
		key := strconv.FormatInt(userID, 10)
		err := s.Producer.SendMessage("user-segments", key, event)
		if err != nil {
			return err
		}
	}

	for _, slug := range slugsToDelete {
		event := map[string]interface{}{
			"user_id": userID,
			"segment": slug,
			"action":  "delete",
		}
		key := strconv.FormatInt(userID, 10)
		err := s.Producer.SendMessage("user-segments", key, event)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *UserSegmentService) DeleteUserSegment(userID int64, slug models.Slug) error {
	if err := s.Repo.DeleteUserSegment(userID, slug); err != nil {
		return err
	}
	return nil
}

func (s *UserSegmentService) ProcessKafkaMessage(message *sarama.ConsumerMessage) {
	var event struct {
		Action  string `json:"action"`
		UserID  int64  `json:"user_id"`
		Segment string `json:"segment"`
		TTL     string `json:"ttl"`
	}
	if err := json.Unmarshal(message.Value, &event); err != nil {
		log.Printf("Failed to parse Kafka message: %v", err)
		return
	}

	log.Printf("Processing Kafka message: action=%s, userID=%d, segment=%s, ttl=%s", event.Action, event.UserID, event.Segment, event.TTL)

	segmentSlug := models.Slug(event.Segment)

	var ttl *time.Time
	if event.TTL != "" {
		parsedTTL, err := time.Parse(time.RFC3339, event.TTL)
		if err != nil {
			log.Printf("Invalid TTL format in Kafka message: %v", err)
			return
		}
		ttl = &parsedTTL
	}

	switch event.Action {
	case "add":
		err := s.Repo.UpdateUserSegments([]models.Slug{segmentSlug}, nil, event.UserID, ttl)
		if err != nil {
			log.Printf("Failed to add user to segment: %v", err)
			return
		}

		err = s.HistoryRepo.SaveHistoryEntry(models.UserSegmentsHistory{
			UserID:        event.UserID,
			SegmentSlug:   segmentSlug,
			OperationType: models.ADD,
			OperationDate: time.Now(),
		})
		if err != nil {
			log.Printf("Failed to save add action to history: %v", err)
			return
		}

		log.Printf("User %d successfully added to segment %s", event.UserID, event.Segment)

	case "remove":
		err := s.DeleteUserSegment(event.UserID, segmentSlug)
		if err != nil {
			log.Printf("Failed to delete user segment: %v", err)
			return
		}

		err = s.HistoryRepo.SaveHistoryEntry(models.UserSegmentsHistory{
			UserID:        event.UserID,
			SegmentSlug:   segmentSlug,
			OperationType: models.DELETE,
			OperationDate: time.Now(),
		})
		if err != nil {
			log.Printf("Failed to save delete action to history: %v", err)
			return
		}

		log.Printf("User %d successfully removed from segment %s", event.UserID, event.Segment)

	default:
		log.Printf("Unknown action: %s", event.Action)
	}
}

func (s *UserSegmentService) ProcessTTLExpiryMessage(msg *sarama.ConsumerMessage) error {
	var event struct {
		UserID   int64  `json:"user_id"`
		Segment  string `json:"segment"`
		ExpireAt string `json:"ttl"`
	}

	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return fmt.Errorf("failed to parse Kafka message: %v", err)
	}

	expireTime, err := time.Parse(time.RFC3339, event.ExpireAt)
	if err != nil {
		return fmt.Errorf("invalid expire_at format: %v", err)
	}
	if time.Now().Before(expireTime) {
		log.Printf("TTL not expired yet for segment %s (user_id: %d)", event.Segment, event.UserID)
		return nil
	}

	if err := s.DeleteUserSegment(event.UserID, models.Slug(event.Segment)); err != nil {
		return fmt.Errorf("failed to delete expired segment: %v", err)
	}

	log.Printf("Successfully deleted expired segment %s for user %d", event.Segment, event.UserID)
	return nil
}
