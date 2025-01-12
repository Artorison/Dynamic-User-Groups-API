package services_test

import (
	"API/internal/models"
	"API/internal/repository/mocks"
	"API/internal/services"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserService_GetAllUsers(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	userService := services.NewUserService(mockRepo)

	expectedUsers := []models.Users{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}

	mockRepo.On("GetAllUsersDB").Return(expectedUsers, nil)

	users, err := userService.GetAllUsers()

	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)

	mockRepo.AssertExpectations(t)
}

func TestUserService_CreateUser(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	userService := services.NewUserService(mockRepo)

	newUser := &models.Users{ID: 3, Name: "Charlie"}

	mockRepo.On("CreateUserDB", newUser).Return(nil)

	err := userService.CreateUser(newUser)

	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "CreateUserDB", newUser)
}

func TestUserService_CreateUser_Error(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	userService := services.NewUserService(mockRepo)

	newUser := &models.Users{ID: 3, Name: "Charlie"}

	mockRepo.On("CreateUserDB", newUser).Return(errors.New("failed to create user"))

	err := userService.CreateUser(newUser)

	assert.Error(t, err)
	mockRepo.AssertCalled(t, "CreateUserDB", newUser)
}

func TestUserService_DeleteUser(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	userService := services.NewUserService(mockRepo)

	userID := int64(1)

	mockRepo.On("DeleteUserDB", userID).Return(nil)

	err := userService.DeleteUser(userID)

	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "DeleteUserDB", userID)
}

func TestUserService_DeleteUser_Error(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	userService := services.NewUserService(mockRepo)

	userID := int64(1)

	mockRepo.On("DeleteUserDB", userID).Return(errors.New("failed to delete user"))

	err := userService.DeleteUser(userID)

	assert.Error(t, err)
	mockRepo.AssertCalled(t, "DeleteUserDB", userID)
}
