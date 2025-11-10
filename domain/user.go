package domain

import "time"

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
