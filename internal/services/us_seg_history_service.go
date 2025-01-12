package services

import (
	"API/internal/repository"
	"API/internal/utils"
	"encoding/csv"
	"fmt"
	"os"
)

//go:generate mockery --name=IUserSegmentHistoryService --output=mocks --outpkg=mocks
type IUserSegmentHistoryService interface {
	GenerateUserHistoryCSV(userID int64, date string) (string, error)
}

type UserSegmentHistoryService struct {
	Repository repository.UserSegmentHistoryRepository
}

func NewUserSegmentHistoryService(repo repository.UserSegmentHistoryRepository) *UserSegmentHistoryService {
	return &UserSegmentHistoryService{Repository: repo}
}

func (s *UserSegmentHistoryService) GenerateUserHistoryCSV(userID int64, date string) (string, error) {
	start, end, err := utils.ParseYearMonth(date)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %w", err)
	}

	histories, err := s.Repository.GetUserHistory(userID, start, end)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve user history: %w", err)
	}

	fileName := fmt.Sprintf("user_%d_history_%s.csv", userID, date)
	filePath := fmt.Sprintf("csv_reports/%s", fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"UserID", "SegmentSlug", "OperationType", "OperationDate"}); err != nil {
		return "", fmt.Errorf("failed to write CSV header: %w", err)
	}

	for _, history := range histories {
		record := []string{
			fmt.Sprintf("%d", history.UserID),
			string(history.SegmentSlug),
			string(history.OperationType),
			history.OperationDate.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(record); err != nil {
			return "", fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	return fileName, nil
}
