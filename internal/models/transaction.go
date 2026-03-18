package models

import (
	"time"

	"gorm.io/gorm"
)

// Transaction represents a financial transaction
type Transaction struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	AccountID       uint           `gorm:"not null;index" json:"account_id"`
	Type            string         `gorm:"not null" json:"type"` // deposit, withdrawal, transfer
	Amount          float64        `gorm:"not null" json:"amount"`
	Currency        string         `gorm:"not null;default:'USD'" json:"currency"`
	Description     string         `json:"description"`
	Status          string         `gorm:"not null;default:'pending'" json:"status"` // pending, completed, failed
	ReferenceNumber string         `gorm:"uniqueIndex" json:"reference_number"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	Account         Account        `gorm:"foreignKey:AccountID" json:"account,omitempty"`
}
