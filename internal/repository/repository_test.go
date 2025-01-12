package repository

import (
	"API/internal/models"
	"API/internal/repository/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllUsersDB(t *testing.T) {
	mockRepo := new(mocks.UserRepository)

	t.Run("should return all users successfully", func(t *testing.T) {
		expectedUsers := []models.Users{
			{ID: 1, Name: "john"},
			{ID: 2, Name: "Artorias"},
		}
		mockRepo.On("GetAllUsersDB").Return(expectedUsers, nil).Once()

		actualUsers, err := mockRepo.GetAllUsersDB()
		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, actualUsers)

		mockRepo.AssertExpectations(t)

	})

	t.Run("should return error when repository fails", func(t *testing.T) {

		mockRepo.On("GetAllUsersDB").Return(nil, errors.New("database error")).Once()

		actualUsers, err := mockRepo.GetAllUsersDB()

		assert.Nil(t, actualUsers)

		assert.EqualError(t, err, "database error")

		mockRepo.AssertExpectations(t)

	})

}

func TestCreateUserDB(t *testing.T) {
	mockRepo := new(mocks.UserRepository)

	t.Run("should return nil error when user is created successfully", func(t *testing.T) {
		createdUsers := []models.Users{
			{ID: 1, Name: "john"},
			{ID: 2, Name: "Artorias"},
		}

		for _, user := range createdUsers {
			mockRepo.On("CreateUserDB", &user).Return(nil).Once()

			err := mockRepo.CreateUserDB(&user)

			assert.NoError(t, err)
		}

		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		invalidUser := models.Users{
			ID:   -4,
			Name: "awd",
		}
		mockRepo.On("CreateUserDB", &invalidUser).Return(errors.New("invalid userID")).Once()

		err := mockRepo.CreateUserDB(&invalidUser)

		assert.EqualError(t, err, "invalid userID")

		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when user is nil", func(t *testing.T) {
		mockRepo.On("CreateUserDB", (*models.Users)(nil)).Return(errors.New("user cannot be nil")).Once()

		err := mockRepo.CreateUserDB(nil)

		assert.EqualError(t, err, "user cannot be nil")

		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteUserDB(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	t.Run("should delete user successfully", func(t *testing.T) {
		var userID int64 = 1
		mockRepo.On("DeleteUserDB", userID).Return(nil).Once()

		err := mockRepo.DeleteUserDB(userID)

		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		var userID int64 = 2
		mockRepo.On("DeleteUserDB", userID).Return(errors.New("delete error")).Once()

		err := mockRepo.DeleteUserDB(userID)

		assert.EqualError(t, err, "delete error")

		mockRepo.AssertExpectations(t)
	})
}
