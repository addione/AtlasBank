package models

import (
	"time"
)

// OTP represents a one-time password for user verification
type OTP struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	OTPCode   string    `gorm:"size:4;not null" json:"-"`
	Action    string    `gorm:"size:50;not null" json:"action"` // account_verification, password_reset, etc.
	IsUsed    bool      `gorm:"not null;default:false" json:"is_used"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
