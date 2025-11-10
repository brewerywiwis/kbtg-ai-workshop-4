package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"workshop4-backend/internal/domain"
)

// Mock repositories for transfer service tests
type MockTransferRepository struct {
	mock.Mock
}

func (m *MockTransferRepository) Create(transfer *domain.Transfer) error {
	args := m.Called(transfer)
	return args.Error(0)
}

func (m *MockTransferRepository) GetByIdempotencyKey(key string) (*domain.Transfer, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Transfer), args.Error(1)
}

func (m *MockTransferRepository) GetByUserID(userID int, page, pageSize int) ([]domain.Transfer, int, error) {
	args := m.Called(userID, page, pageSize)
	return args.Get(0).([]domain.Transfer), args.Get(1).(int), args.Error(2)
}

func (m *MockTransferRepository) UpdateStatus(id int, status domain.TransferStatus, completedAt *string, failReason *string) error {
	args := m.Called(id, status, completedAt, failReason)
	return args.Error(0)
}

type MockPointLedgerRepository struct {
	mock.Mock
}

func (m *MockPointLedgerRepository) Create(entry *domain.PointLedger) error {
	args := m.Called(entry)
	return args.Error(0)
}

func (m *MockPointLedgerRepository) GetByUserID(userID int) ([]domain.PointLedger, error) {
	args := m.Called(userID)
	return args.Get(0).([]domain.PointLedger), args.Error(1)
}

func (m *MockPointLedgerRepository) GetUserBalance(userID int) (int, error) {
	args := m.Called(userID)
	return args.Get(0).(int), args.Error(1)
}

// Focus on validation logic tests only for unit tests
// Integration tests with real database transactions would be separate
func TestTransferService_CreateTransfer_ValidationTests(t *testing.T) {
	transferRepo := new(MockTransferRepository)
	ledgerRepo := new(MockPointLedgerRepository)
	userRepo := new(MockUserRepository)
	service := NewTransferService(transferRepo, ledgerRepo, userRepo, nil)

	t.Run("same user transfer", func(t *testing.T) {
		result, err := service.CreateTransfer(1, 1, 500, nil)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrSelfTransfer, err)
	})

	t.Run("invalid amount", func(t *testing.T) {
		result, err := service.CreateTransfer(1, 2, 0, nil)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "amount must be greater than 0")
	})

	t.Run("user not found", func(t *testing.T) {
		userRepo.On("GetByID", 999).Return(nil, errors.New("user not found"))
		result, err := service.CreateTransfer(999, 2, 500, nil)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrUserNotFound, err)
	})

	t.Run("insufficient balance", func(t *testing.T) {
		fromUser := &domain.User{ID: 1, Points: 100}
		toUser := &domain.User{ID: 2, Points: 500}
		userRepo.On("GetByID", 1).Return(fromUser, nil)
		userRepo.On("GetByID", 2).Return(toUser, nil)
		ledgerRepo.On("GetUserBalance", 1).Return(100, nil)

		result, err := service.CreateTransfer(1, 2, 500, nil)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrInsufficientBalance, err)
	})
}
