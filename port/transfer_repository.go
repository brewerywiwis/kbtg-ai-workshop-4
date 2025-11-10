package port

import "workshop4-backend/domain"

type TransferRepository interface {
	Create(transfer *domain.Transfer) error
	GetByIdempotencyKey(key string) (*domain.Transfer, error)
	GetByUserID(userID int, page, pageSize int) ([]domain.Transfer, int, error)
	UpdateStatus(id int, status domain.TransferStatus, completedAt *string, failReason *string) error
}

type PointLedgerRepository interface {
	Create(entry *domain.PointLedger) error
	GetByUserID(userID int) ([]domain.PointLedger, error)
	GetUserBalance(userID int) (int, error)
}

// Additional methods for UserRepository to support transfers
type UserRepositoryWithBalance interface {
	UserRepository
	UpdatePoints(userID int, newBalance int) error
	GetUserBalance(userID int) (int, error)
}
