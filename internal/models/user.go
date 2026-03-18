package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a bank user
type User struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Email      string         `gorm:"uniqueIndex;not null" json:"email"`
	FirstName  string         `gorm:"not null" json:"first_name"`
	LastName   string         `gorm:"not null" json:"last_name"`
	Password   string         `gorm:"not null" json:"-"`
	IsVerified bool           `gorm:"not null;default:false" json:"is_verified"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	Accounts   []Account      `gorm:"foreignKey:UserID" json:"accounts,omitempty"`
	OTPs       []OTP          `gorm:"foreignKey:UserID" json:"otps,omitempty"`
}
