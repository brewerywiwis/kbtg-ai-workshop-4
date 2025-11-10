package service

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"workshop4-backend/internal/domain"
	"workshop4-backend/internal/port"

	"github.com/google/uuid"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrSelfTransfer        = errors.New("cannot transfer to yourself")
	ErrUserNotFound        = errors.New("user not found")
	ErrTransferNotFound    = errors.New("transfer not found")
)

type TransferService struct {
	transferRepo port.TransferRepository
	ledgerRepo   port.PointLedgerRepository
	userRepo     port.UserRepository
	db           *sql.DB // For transaction support
}

func NewTransferService(
	transferRepo port.TransferRepository,
	ledgerRepo port.PointLedgerRepository,
	userRepo port.UserRepository,
	db *sql.DB,
) *TransferService {
	return &TransferService{
		transferRepo: transferRepo,
		ledgerRepo:   ledgerRepo,
		userRepo:     userRepo,
		db:           db,
	}
}

func (s *TransferService) CreateTransfer(fromUserID, toUserID, amount int, note *string) (*domain.Transfer, error) {
	// Validate input
	if fromUserID == toUserID {
		return nil, ErrSelfTransfer
	}

	if amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	// Check if users exist
	fromUser, err := s.userRepo.GetByID(fromUserID)
	if err != nil || fromUser == nil {
		return nil, ErrUserNotFound
	}

	toUser, err := s.userRepo.GetByID(toUserID)
	if err != nil || toUser == nil {
		return nil, ErrUserNotFound
	}

	// Get current balance
	currentBalance, err := s.ledgerRepo.GetUserBalance(fromUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user balance: %w", err)
	}

	if currentBalance < amount {
		return nil, ErrInsufficientBalance
	}

	// Generate idempotency key
	idemKey := uuid.New().String()

	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Create transfer record
	now := time.Now()
	transfer := &domain.Transfer{
		FromUserID:     fromUserID,
		ToUserID:       toUserID,
		Amount:         amount,
		Status:         domain.TransferStatusCompleted, // For now, assume immediate completion
		Note:           note,
		IdempotencyKey: idemKey,
		CreatedAt:      now,
		UpdatedAt:      now,
		CompletedAt:    &now,
	}

	if err := s.transferRepo.Create(transfer); err != nil {
		return nil, fmt.Errorf("failed to create transfer: %w", err)
	}

	// Create ledger entries
	// Debit from sender
	debitEntry := &domain.PointLedger{
		UserID:       fromUserID,
		Change:       -amount,
		BalanceAfter: currentBalance - amount,
		EventType:    domain.EventTypeTransferOut,
		TransferID:   &transfer.ID,
		Reference:    nil,
		Metadata:     nil,
		CreatedAt:    now,
	}

	if err := s.ledgerRepo.Create(debitEntry); err != nil {
		return nil, fmt.Errorf("failed to create debit ledger entry: %w", err)
	}

	// Get recipient balance
	recipientBalance, err := s.ledgerRepo.GetUserBalance(toUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipient balance: %w", err)
	}

	// Credit to recipient
	creditEntry := &domain.PointLedger{
		UserID:       toUserID,
		Change:       amount,
		BalanceAfter: recipientBalance + amount,
		EventType:    domain.EventTypeTransferIn,
		TransferID:   &transfer.ID,
		Reference:    nil,
		Metadata:     nil,
		CreatedAt:    now,
	}

	if err := s.ledgerRepo.Create(creditEntry); err != nil {
		return nil, fmt.Errorf("failed to create credit ledger entry: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return transfer, nil
}

func (s *TransferService) GetTransferByIdempotencyKey(key string) (*domain.Transfer, error) {
	transfer, err := s.transferRepo.GetByIdempotencyKey(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get transfer: %w", err)
	}
	if transfer == nil {
		return nil, ErrTransferNotFound
	}
	return transfer, nil
}

func (s *TransferService) GetTransfersByUserID(userID int, page, pageSize int) ([]domain.Transfer, int, error) {
	// Validate pagination
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 200 {
		pageSize = 20
	}

	transfers, total, err := s.transferRepo.GetByUserID(userID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get transfers: %w", err)
	}

	return transfers, total, nil
}
