package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"workshop4-backend/internal/domain"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetAll() ([]domain.User, error) {
	args := m.Called()
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(id int) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestUserService_Create_Success(t *testing.T) {
	repo := new(MockUserRepository)
	service := NewUserService(repo)
	user := &domain.User{
		Name:  "Test User",
		Email: "test@example.com",
		Phone: "081-234-5678",
	}
	repo.On("Create", user).Return(nil)
	err := service.CreateUser(user)
	assert.NoError(t, err)
	assert.NotZero(t, user.CreatedAt)
	assert.NotZero(t, user.UpdatedAt)
	repo.AssertExpectations(t)
}

func TestUserService_Create_ValidationError(t *testing.T) {
	repo := new(MockUserRepository)
	service := NewUserService(repo)
	user := &domain.User{Name: ""} // Invalid name
	// No repo.On expectation, since validation should fail before repo.Create is called
	err := service.CreateUser(user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation")
}

func TestUserService_Create_RepoError(t *testing.T) {
	repo := new(MockUserRepository)
	service := NewUserService(repo)
	user := &domain.User{
		Name:  "Test User",
		Email: "test@example.com",
		Phone: "081-234-5678",
	}
	repo.On("Create", user).Return(errors.New("db error"))
	err := service.CreateUser(user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
}

func TestUserService_GetAllUsers_Success(t *testing.T) {
	repo := new(MockUserRepository)
	service := NewUserService(repo)
	expectedUsers := []domain.User{
		{ID: 1, Name: "User 1", Email: "user1@example.com", Phone: "081-111-1111"},
		{ID: 2, Name: "User 2", Email: "user2@example.com", Phone: "081-222-2222"},
	}
	repo.On("GetAll").Return(expectedUsers, nil)
	users, err := service.GetAllUsers()
	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
	repo.AssertExpectations(t)
}

func TestUserService_GetUserByID_Success(t *testing.T) {
	repo := new(MockUserRepository)
	service := NewUserService(repo)
	expectedUser := &domain.User{ID: 1, Name: "Test User", Email: "test@example.com", Phone: "081-234-5678"}
	repo.On("GetByID", 1).Return(expectedUser, nil)
	user, err := service.GetUserByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	repo.AssertExpectations(t)
}

func TestUserService_GetUserByID_NotFound(t *testing.T) {
	repo := new(MockUserRepository)
	service := NewUserService(repo)
	repo.On("GetByID", 999).Return(nil, errors.New("user not found"))
	user, err := service.GetUserByID(999)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "user not found")
	repo.AssertExpectations(t)
}

func TestUserService_DeleteUser_Success(t *testing.T) {
	repo := new(MockUserRepository)
	service := NewUserService(repo)
	repo.On("Delete", 1).Return(nil)
	err := service.DeleteUser(1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}
