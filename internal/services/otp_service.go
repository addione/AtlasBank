package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/atlasbank/api/internal/models"
	"gorm.io/gorm"
)

// OTPService handles OTP-related business logic
type OTPService struct {
	db *gorm.DB
}

// NewOTPService creates a new OTPService
func NewOTPService(db *gorm.DB) *OTPService {
	return &OTPService{
		db: db,
	}
}

// OTP Actions
const (
	OTPActionAccountVerification = "account_verification"
	OTPActionPasswordReset       = "password_reset"
	OTPActionLogin               = "login"
	OTPActionTransaction         = "transaction"
)

// GenerateOTP creates a new OTP for a user
// For now, hardcoded as "0000" as per requirements
func (s *OTPService) GenerateOTP(ctx context.Context, userID uint, action string) (*models.OTP, error) {
	// Invalidate any existing unused OTPs for this user and action
	if err := s.db.WithContext(ctx).
		Model(&models.OTP{}).
		Where("user_id = ? AND action = ? AND is_used = ?", userID, action, false).
		Update("is_used", true).Error; err != nil {
		return nil, fmt.Errorf("failed to invalidate old OTPs: %w", err)
	}

	// Create new OTP (hardcoded as "0000" for now)
	otp := &models.OTP{
		UserID:    userID,
		OTPCode:   "0000", // Hardcoded as per requirements
		Action:    action,
		IsUsed:    false,
		ExpiresAt: time.Now().Add(15 * time.Minute), // OTP valid for 15 minutes
	}

	if err := s.db.WithContext(ctx).Create(otp).Error; err != nil {
		return nil, fmt.Errorf("failed to create OTP: %w", err)
	}

	return otp, nil
}

// VerifyOTP verifies an OTP for a user
func (s *OTPService) VerifyOTP(ctx context.Context, userID uint, otpCode, action string) error {
	var otp models.OTP

	// Find the most recent unused OTP for this user and action
	if err := s.db.WithContext(ctx).
		Where("user_id = ? AND action = ? AND is_used = ?", userID, action, false).
		Order("created_at DESC").
		First(&otp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("invalid or expired OTP")
		}
		return fmt.Errorf("failed to find OTP: %w", err)
	}

	// Check if OTP has expired
	if time.Now().After(otp.ExpiresAt) {
		return errors.New("OTP has expired")
	}

	// Verify OTP code
	if otp.OTPCode != otpCode {
		return errors.New("invalid OTP code")
	}

	// Mark OTP as used
	otp.IsUsed = true
	if err := s.db.WithContext(ctx).Save(&otp).Error; err != nil {
		return fmt.Errorf("failed to mark OTP as used: %w", err)
	}

	return nil
}

// GetActiveOTP retrieves the active OTP for a user and action (for testing/debugging)
func (s *OTPService) GetActiveOTP(ctx context.Context, userID uint, action string) (*models.OTP, error) {
	var otp models.OTP

	if err := s.db.WithContext(ctx).
		Where("user_id = ? AND action = ? AND is_used = ? AND expires_at > ?", userID, action, false, time.Now()).
		Order("created_at DESC").
		First(&otp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no active OTP found")
		}
		return nil, fmt.Errorf("failed to get active OTP: %w", err)
	}

	return &otp, nil
}

// CleanupExpiredOTPs removes expired OTPs from the database
func (s *OTPService) CleanupExpiredOTPs(ctx context.Context) error {
	result := s.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&models.OTP{})

	if result.Error != nil {
		return fmt.Errorf("failed to cleanup expired OTPs: %w", result.Error)
	}

	if result.RowsAffected > 0 {
		fmt.Printf("Cleaned up %d expired OTPs\n", result.RowsAffected)
	}

	return nil
}
