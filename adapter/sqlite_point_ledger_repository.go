package adapter

import (
	"database/sql"
	"time"

	"workshop4-backend/domain"
	"workshop4-backend/port"
)

type SqlitePointLedgerRepository struct {
	db *sql.DB
}

func NewSqlitePointLedgerRepository(db *sql.DB) port.PointLedgerRepository {
	return &SqlitePointLedgerRepository{db: db}
}

func (r *SqlitePointLedgerRepository) Create(entry *domain.PointLedger) error {
	query := `
		INSERT INTO point_ledger (user_id, change, balance_after, event_type, transfer_id, reference, metadata, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query,
		entry.UserID,
		entry.Change,
		entry.BalanceAfter,
		entry.EventType,
		entry.TransferID,
		entry.Reference,
		entry.Metadata,
		entry.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	entry.ID = int(id)
	return nil
}

func (r *SqlitePointLedgerRepository) GetByUserID(userID int) ([]domain.PointLedger, error) {
	query := `
		SELECT id, user_id, change, balance_after, event_type, transfer_id, reference, metadata, created_at
		FROM point_ledger WHERE user_id = ? ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []domain.PointLedger
	for rows.Next() {
		var entry domain.PointLedger
		var createdAtStr string
		var transferID sql.NullInt64
		var reference, metadata sql.NullString

		err := rows.Scan(
			&entry.ID,
			&entry.UserID,
			&entry.Change,
			&entry.BalanceAfter,
			&entry.EventType,
			&transferID,
			&reference,
			&metadata,
			&createdAtStr)
		if err != nil {
			return nil, err
		}

		// Parse time string
		entry.CreatedAt, err = time.Parse("2006-01-02T15:04:05Z07:00", createdAtStr)
		if err != nil {
			return nil, err
		}

		// Handle nullable fields
		if transferID.Valid {
			tid := int(transferID.Int64)
			entry.TransferID = &tid
		}
		if reference.Valid {
			entry.Reference = &reference.String
		}
		if metadata.Valid {
			entry.Metadata = &metadata.String
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

func (r *SqlitePointLedgerRepository) GetUserBalance(userID int) (int, error) {
	query := `
		SELECT balance_after FROM point_ledger
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT 1
	`

	var balance int
	err := r.db.QueryRow(query, userID).Scan(&balance)
	if err == sql.ErrNoRows {
		// No ledger entries, check user table for initial points
		userQuery := `SELECT points FROM users WHERE id = ?`
		err = r.db.QueryRow(userQuery, userID).Scan(&balance)
		if err != nil {
			return 0, err
		}
		return balance, nil
	}
	if err != nil {
		return 0, err
	}

	return balance, nil
}
