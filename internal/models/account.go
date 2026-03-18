package models

import (
	"time"

	"gorm.io/gorm"
)

// Account represents a bank account
type Account struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"not null;index" json:"user_id"`
	AccountNumber string         `gorm:"uniqueIndex;not null" json:"account_number"`
	AccountType   string         `gorm:"not null" json:"account_type"` // savings, checking, etc.
	Balance       float64        `gorm:"not null;default:0" json:"balance"`
	Currency      string         `gorm:"not null;default:'USD'" json:"currency"`
	Status        string         `gorm:"not null;default:'active'" json:"status"` // active, frozen, closed
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	User          User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Transactions  []Transaction  `gorm:"foreignKey:AccountID" json:"transactions,omitempty"`
}
