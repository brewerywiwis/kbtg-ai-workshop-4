package port

import "workshop4-backend/internal/domain"

type UserRepository interface {
	GetAll() ([]domain.User, error)
	GetByID(id int) (*domain.User, error)
	Create(user *domain.User) error
	Update(user *domain.User) error
	Delete(id int) error
}
