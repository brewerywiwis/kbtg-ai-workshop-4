package service

import (
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
	return s.repo.Create(user)
}

func (s *UserService) UpdateUser(user *domain.User) error {
	return s.repo.Update(user)
}

func (s *UserService) DeleteUser(id int) error {
	return s.repo.Delete(id)
}
