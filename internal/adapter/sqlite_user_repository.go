package adapter

import (
	"database/sql"

	"workshop4-backend/internal/domain"
	"workshop4-backend/internal/port"
)

type SqliteUserRepository struct {
	db *sql.DB
}

func NewSqliteUserRepository(db *sql.DB) port.UserRepository {
	return &SqliteUserRepository{db: db}
}

func (r *SqliteUserRepository) GetAll() ([]domain.User, error) {
	rows, err := r.db.Query(`SELECT id, name, phone, email, member_since, membership_level, member_id, points, created_at, updated_at FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Name, &user.Phone, &user.Email, &user.MemberSince, &user.MembershipLevel, &user.MemberID, &user.Points, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *SqliteUserRepository) GetByID(id int) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRow(`SELECT id, name, phone, email, member_since, membership_level, member_id, points, created_at, updated_at FROM users WHERE id = ?`, id).
		Scan(&user.ID, &user.Name, &user.Phone, &user.Email, &user.MemberSince, &user.MembershipLevel, &user.MemberID, &user.Points, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *SqliteUserRepository) Create(user *domain.User) error {
	result, err := r.db.Exec(`INSERT INTO users (name, phone, email, member_since, membership_level, member_id, points) VALUES (?, ?, ?, ?, ?, ?, ?)`, user.Name, user.Phone, user.Email, user.MemberSince, user.MembershipLevel, user.MemberID, user.Points)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = int(id)
	return nil
}

func (r *SqliteUserRepository) Update(user *domain.User) error {
	_, err := r.db.Exec(`UPDATE users SET name = ?, phone = ?, email = ?, member_since = ?, membership_level = ?, member_id = ?, points = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, user.Name, user.Phone, user.Email, user.MemberSince, user.MembershipLevel, user.MemberID, user.Points, user.ID)
	return err
}

func (r *SqliteUserRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE id = ?`, id)
	return err
}

func (r *SqliteUserRepository) UpdatePoints(userID int, newBalance int) error {
	_, err := r.db.Exec(`UPDATE users SET points = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, newBalance, userID)
	return err
}

func (r *SqliteUserRepository) GetUserBalance(userID int) (int, error) {
	var balance int
	err := r.db.QueryRow(`SELECT points FROM users WHERE id = ?`, userID).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}
