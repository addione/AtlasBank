package services

import (
	"context"
	"fmt"
	"time"

	"github.com/atlasbank/api/internal/models"
	"gorm.io/gorm"
)

// NotificationService handles notification-related business logic
type NotificationService struct {
	db *gorm.DB
}

// NewNotificationService creates a new NotificationService
func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{
		db: db,
	}
}

// Notification types
const (
	NotificationTypeEmail = "email"
	NotificationTypeSMS   = "sms"
	NotificationTypePush  = "push"
)

// Notification status
const (
	NotificationStatusPending = "pending"
	NotificationStatusSent    = "sent"
	NotificationStatusFailed  = "failed"
)

// LogNotification creates a notification record in the database
func (s *NotificationService) LogNotification(ctx context.Context, userID uint, notificationType, recipient, subject, message string) (*models.Notification, error) {
	notification := &models.Notification{
		UserID:           userID,
		NotificationType: notificationType,
		Recipient:        recipient,
		Subject:          subject,
		Message:          message,
		Status:           NotificationStatusPending,
	}

	if err := s.db.WithContext(ctx).Create(notification).Error; err != nil {
		return nil, fmt.Errorf("failed to log notification: %w", err)
	}

	return notification, nil
}

// MarkNotificationSent marks a notification as sent
func (s *NotificationService) MarkNotificationSent(ctx context.Context, notificationID uint) error {
	now := time.Now()
	return s.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("id = ?", notificationID).
		Updates(map[string]interface{}{
			"status":  NotificationStatusSent,
			"sent_at": now,
		}).Error
}

// MarkNotificationFailed marks a notification as failed with error message
func (s *NotificationService) MarkNotificationFailed(ctx context.Context, notificationID uint, errorMessage string) error {
	return s.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("id = ?", notificationID).
		Updates(map[string]interface{}{
			"status":        NotificationStatusFailed,
			"error_message": errorMessage,
		}).Error
}

// GetUserNotifications retrieves all notifications for a user
func (s *NotificationService) GetUserNotifications(ctx context.Context, userID uint, limit int) ([]models.Notification, error) {
	var notifications []models.Notification

	query := s.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&notifications).Error; err != nil {
		return nil, fmt.Errorf("failed to get user notifications: %w", err)
	}

	return notifications, nil
}

// GetNotificationStats returns notification statistics for a user
func (s *NotificationService) GetNotificationStats(ctx context.Context, userID uint) (map[string]int64, error) {
	stats := make(map[string]int64)

	// Count by status
	var statusCounts []struct {
		Status string
		Count  int64
	}

	if err := s.db.WithContext(ctx).
		Model(&models.Notification{}).
		Select("status, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("status").
		Scan(&statusCounts).Error; err != nil {
		return nil, fmt.Errorf("failed to get notification stats: %w", err)
	}

	for _, sc := range statusCounts {
		stats[sc.Status] = sc.Count
	}

	// Count by type
	var typeCounts []struct {
		NotificationType string
		Count            int64
	}

	if err := s.db.WithContext(ctx).
		Model(&models.Notification{}).
		Select("notification_type, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("notification_type").
		Scan(&typeCounts).Error; err != nil {
		return nil, fmt.Errorf("failed to get notification type stats: %w", err)
	}

	for _, tc := range typeCounts {
		stats[tc.NotificationType] = tc.Count
	}

	return stats, nil
}
