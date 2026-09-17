package service

import (
	"user-transfer-go/internal/models"
	"user-transfer-go/internal/repository"
)

// UserService berisi business logic terkait user.
type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetAllUsers() ([]*models.User, error) {
	return s.userRepo.FindAll()
}

func (s *UserService) GetUserByID(id int64) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) CreateUser(name string, age int, balance float64) (*models.User, error) {
	return s.userRepo.Create(name, age, balance)
}

func (s *UserService) DeleteUser(id int64) error {
	deleted, err := s.userRepo.Delete(id)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrUserNotFound
	}
	return nil
}
