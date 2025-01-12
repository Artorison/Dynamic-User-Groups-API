package repository

import (
	"API/internal/models"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetUserSegmentsDВ(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	repo := NewUserSegmentRepository(mockDB, nil, nil, nil)

	t.Run("should return user segments successfully", func(t *testing.T) {
		expectedSegments := models.UserSegments{
			UserID: 1,
			Segments: []models.Slug{
				"segment1",
				"segment2",
			},
		}

		rows := sqlmock.NewRows([]string{"slug"}).
			AddRow("segment1").
			AddRow("segment2")

		query := `
			SELECT s.slug 
			FROM user_segments us
			JOIN segments s ON us.segment_id = s.id
			WHERE us.user_id = $1;
		`

		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(1).WillReturnRows(rows)

		actualSegments, err := repo.GetUserSegmentsDВ(1)

		assert.NoError(t, err)
		assert.Equal(t, expectedSegments, actualSegments)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
	t.Run("should return error when query fails", func(t *testing.T) {
		query := `
			SELECT s.slug 
			FROM user_segments us
			JOIN segments s ON us.segment_id = s.id
			WHERE us.user_id = $1;
		`
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(1).WillReturnError(fmt.Errorf("database error"))

		actualSegments, err := repo.GetUserSegmentsDВ(1)

		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
		assert.Empty(t, actualSegments)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

}
