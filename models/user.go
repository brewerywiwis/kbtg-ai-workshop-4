package models

import "time"

// User represents a user profile based on the LBK membership profile
type User struct {
	ID              int       `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`                         // ชื่อ (สมชาย ใจดี)
	Phone           string    `json:"phone" db:"phone"`                       // เบอร์โทรศัพท์ (081-234-5678)
	Email           string    `json:"email" db:"email"`                       // อีเมล (somchai@example.com)
	MemberSince     string    `json:"member_since" db:"member_since"`         // วันที่สมัครสมาชิก (15/6/2566)
	MembershipLevel string    `json:"membership_level" db:"membership_level"` // ระดับสมาชิก (Gold)
	MemberID        string    `json:"member_id" db:"member_id"`               // รหัสสมาชิก (LBK001234)
	Points          int       `json:"points" db:"points"`                     // แต้มคงเหลือ (15,420 แต้ม)
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}
