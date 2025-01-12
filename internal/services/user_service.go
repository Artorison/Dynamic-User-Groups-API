package services

import (
	"API/internal/models"
	"API/internal/repository"
)

//go:generate mockery --name=IUserService --output=mocks --outpkg=mocks
type IUserService interface {
	GetAllUsers() ([]models.Users, error)
	CreateUser(user *models.Users) error
	DeleteUser(userID int64) error
}

type UserService struct {
	Repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) GetAllUsers() ([]models.Users, error) {
	return s.Repo.GetAllUsersDB()
}

func (s *UserService) CreateUser(user *models.Users) error {
	return s.Repo.CreateUserDB(user)
}

func (s *UserService) DeleteUser(userID int64) error {
	return s.Repo.DeleteUserDB(userID)
}
