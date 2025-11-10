package domain

import (
	"errors"
	"strings"
	"time"
)

type User struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	Phone           string    `json:"phone"`
	Email           string    `json:"email"`
	MemberSince     string    `json:"member_since"`
	MembershipLevel string    `json:"membership_level"`
	MemberID        string    `json:"member_id"`
	Points          int       `json:"points"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Validate validates the user data
func (u *User) Validate() error {
	if strings.TrimSpace(u.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(u.Email) == "" {
		return errors.New("email is required")
	}
	if strings.TrimSpace(u.Phone) == "" {
		return errors.New("phone is required")
	}
	return nil
}
