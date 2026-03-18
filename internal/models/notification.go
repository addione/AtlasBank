package models

import (
	"time"
)

// Notification represents a notification sent to a user
type Notification struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	UserID           uint       `gorm:"not null;index" json:"user_id"`
	NotificationType string     `gorm:"size:20;not null" json:"notification_type"` // email, sms, push
	Recipient        string     `gorm:"size:255;not null" json:"recipient"`        // email, phone, device token
	Subject          string     `gorm:"size:255" json:"subject"`
	Message          string     `gorm:"type:text;not null" json:"message"`
	Status           string     `gorm:"size:20;not null;default:'pending'" json:"status"` // pending, sent, failed
	ErrorMessage     string     `gorm:"type:text" json:"error_message,omitempty"`
	SentAt           *time.Time `json:"sent_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	User             User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
