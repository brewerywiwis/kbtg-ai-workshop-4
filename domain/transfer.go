package domain

import "time"

type TransferStatus string

const (
	TransferStatusPending    TransferStatus = "pending"
	TransferStatusProcessing TransferStatus = "processing"
	TransferStatusCompleted  TransferStatus = "completed"
	TransferStatusFailed     TransferStatus = "failed"
	TransferStatusCancelled  TransferStatus = "cancelled"
	TransferStatusReversed   TransferStatus = "reversed"
)

type Transfer struct {
	ID             int            `json:"transferId,omitempty" db:"id"`
	FromUserID     int            `json:"fromUserId" db:"from_user_id"`
	ToUserID       int            `json:"toUserId" db:"to_user_id"`
	Amount         int            `json:"amount" db:"amount"`
	Status         TransferStatus `json:"status" db:"status"`
	Note           *string        `json:"note,omitempty" db:"note"`
	IdempotencyKey string         `json:"idemKey" db:"idempotency_key"`
	CreatedAt      time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time      `json:"updatedAt" db:"updated_at"`
	CompletedAt    *time.Time     `json:"completedAt,omitempty" db:"completed_at"`
	FailReason     *string        `json:"failReason,omitempty" db:"fail_reason"`
}

type EventType string

const (
	EventTypeTransferOut EventType = "transfer_out"
	EventTypeTransferIn  EventType = "transfer_in"
	EventTypeAdjust      EventType = "adjust"
	EventTypeEarn        EventType = "earn"
	EventTypeRedeem      EventType = "redeem"
)

type PointLedger struct {
	ID           int       `json:"id" db:"id"`
	UserID       int       `json:"userId" db:"user_id"`
	Change       int       `json:"change" db:"change"`
	BalanceAfter int       `json:"balanceAfter" db:"balance_after"`
	EventType    EventType `json:"eventType" db:"event_type"`
	TransferID   *int      `json:"transferId,omitempty" db:"transfer_id"`
	Reference    *string   `json:"reference,omitempty" db:"reference"`
	Metadata     *string   `json:"metadata,omitempty" db:"metadata"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
}
