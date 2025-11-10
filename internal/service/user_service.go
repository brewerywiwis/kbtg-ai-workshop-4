package service

import (
	"errors"
	"strings"
	"time"

	"workshop4-backend/internal/domain"
	"workshop4-backend/internal/port"
)

type UserService struct {
	repo port.UserRepository
}

func NewUserService(repo port.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetAllUsers() ([]domain.User, error) {
	return s.repo.GetAll()
}

func (s *UserService) GetUserByID(id int) (*domain.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) CreateUser(user *domain.User) error {
	// Validate user input
	if err := s.validateUser(user); err != nil {
		return err
	}

	// Set timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	return s.repo.Create(user)
}

func (s *UserService) validateUser(user *domain.User) error {
	if strings.TrimSpace(user.Name) == "" {
		return errors.New("validation error: name is required")
	}
	if strings.TrimSpace(user.Email) == "" {
		return errors.New("validation error: email is required")
	}
	if strings.TrimSpace(user.Phone) == "" {
		return errors.New("validation error: phone is required")
	}
	return nil
}

func (s *UserService) UpdateUser(user *domain.User) error {
	return s.repo.Update(user)
}

func (s *UserService) DeleteUser(id int) error {
	return s.repo.Delete(id)
}
