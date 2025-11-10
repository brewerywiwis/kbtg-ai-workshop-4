package adapter

import (
	"database/sql"
	"time"
	"workshop4-backend/domain"
	"workshop4-backend/port"
)

// Helper functions for time parsing
func parseTimeString(timeStr string, target *time.Time) error {
	t, err := time.Parse("2006-01-02T15:04:05Z07:00", timeStr)
	if err != nil {
		return err
	}
	*target = t
	return nil
}

func parseTimeStringPtr(timeStr string, target **time.Time) error {
	t, err := time.Parse("2006-01-02T15:04:05Z07:00", timeStr)
	if err != nil {
		return err
	}
	*target = &t
	return nil
}

func formatTimePtr(t *time.Time) interface{} {
	if t == nil {
		return nil
	}
	return t.Format("2006-01-02T15:04:05Z07:00")
}

type SqliteTransferRepository struct {
	db *sql.DB
}

func NewSqliteTransferRepository(db *sql.DB) port.TransferRepository {
	return &SqliteTransferRepository{db: db}
}

func (r *SqliteTransferRepository) Create(transfer *domain.Transfer) error {
	query := `
		INSERT INTO transfers (from_user_id, to_user_id, amount, status, note, idempotency_key, created_at, updated_at, completed_at, fail_reason)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query,
		transfer.FromUserID,
		transfer.ToUserID,
		transfer.Amount,
		transfer.Status,
		transfer.Note,
		transfer.IdempotencyKey,
		transfer.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		transfer.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		formatTimePtr(transfer.CompletedAt),
		transfer.FailReason)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	transfer.ID = int(id)
	return nil
}

func (r *SqliteTransferRepository) GetByIdempotencyKey(key string) (*domain.Transfer, error) {
	query := `
		SELECT id, from_user_id, to_user_id, amount, status, note, idempotency_key, created_at, updated_at, completed_at, fail_reason
		FROM transfers WHERE idempotency_key = ?
	`
	
	var transfer domain.Transfer
	var createdAtStr, updatedAtStr string
	var completedAtStr, note, failReason sql.NullString
	
	err := r.db.QueryRow(query, key).Scan(
		&transfer.ID,
		&transfer.FromUserID,
		&transfer.ToUserID,
		&transfer.Amount,
		&transfer.Status,
		&note,
		&transfer.IdempotencyKey,
		&createdAtStr,
		&updatedAtStr,
		&completedAtStr,
		&failReason)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Parse time strings
	if err := parseTimeString(createdAtStr, &transfer.CreatedAt); err != nil {
		return nil, err
	}
	if err := parseTimeString(updatedAtStr, &transfer.UpdatedAt); err != nil {
		return nil, err
	}

	// Handle nullable fields
	if note.Valid {
		transfer.Note = &note.String
	}
	if failReason.Valid {
		transfer.FailReason = &failReason.String
	}
	if completedAtStr.Valid {
		if err := parseTimeStringPtr(completedAtStr.String, &transfer.CompletedAt); err != nil {
			return nil, err
		}
	}

	return &transfer, nil
}

func (r *SqliteTransferRepository) GetByUserID(userID int, page, pageSize int) ([]domain.Transfer, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	countQuery := `
		SELECT COUNT(*) FROM transfers 
		WHERE from_user_id = ? OR to_user_id = ?
	`
	var total int
	err := r.db.QueryRow(countQuery, userID, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	query := `
		SELECT id, from_user_id, to_user_id, amount, status, note, idempotency_key, created_at, updated_at, completed_at, fail_reason
		FROM transfers 
		WHERE from_user_id = ? OR to_user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, userID, userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var transfers []domain.Transfer
	for rows.Next() {
		var transfer domain.Transfer
		var createdAtStr, updatedAtStr string
		var completedAtStr, note, failReason sql.NullString

		err := rows.Scan(
			&transfer.ID,
			&transfer.FromUserID,
			&transfer.ToUserID,
			&transfer.Amount,
			&transfer.Status,
			&note,
			&transfer.IdempotencyKey,
			&createdAtStr,
			&updatedAtStr,
			&completedAtStr,
			&failReason)

		if err != nil {
			return nil, 0, err
		}

		// Parse time strings
		if err := parseTimeString(createdAtStr, &transfer.CreatedAt); err != nil {
			return nil, 0, err
		}
		if err := parseTimeString(updatedAtStr, &transfer.UpdatedAt); err != nil {
			return nil, 0, err
		}

		// Handle nullable fields
		if note.Valid {
			transfer.Note = &note.String
		}
		if failReason.Valid {
			transfer.FailReason = &failReason.String
		}
		if completedAtStr.Valid {
			if err := parseTimeStringPtr(completedAtStr.String, &transfer.CompletedAt); err != nil {
				return nil, 0, err
			}
		}

		transfers = append(transfers, transfer)
	}

	return transfers, total, nil
}

func (r *SqliteTransferRepository) UpdateStatus(id int, status domain.TransferStatus, completedAt *string, failReason *string) error {
	query := `
		UPDATE transfers 
		SET status = ?, updated_at = datetime('now'), completed_at = ?, fail_reason = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query, status, completedAt, failReason, id)
	return err
}