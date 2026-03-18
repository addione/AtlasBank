package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/atlasbank/api/internal/models"
	"gorm.io/gorm"
)

// AccountService handles business logic for accounts
type AccountService struct {
	db *gorm.DB
}

// NewAccountService creates a new AccountService
func NewAccountService(db *gorm.DB) *AccountService {
	return &AccountService{
		db: db,
	}
}

// Account types
const (
	AccountTypeSalary    = "salary"
	AccountTypeRegular   = "regular"
	AccountTypeCorporate = "corporate"
)

// Transaction types
const (
	TransactionTypeDeposit    = "deposit"
	TransactionTypeWithdrawal = "withdrawal"
	TransactionTypeTransfer   = "transfer"
)

// Deposit adds money to an account
func (s *AccountService) Deposit(ctx context.Context, accountID uint, amount float64, description string) (*models.Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("deposit amount must be greater than zero")
	}

	// Start transaction
	tx := s.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get account with lock
	var account models.Account
	if err := tx.Clauses().Where("id = ? AND status = ?", accountID, "active").First(&account).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("account not found or inactive")
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	// Update balance
	account.Balance += amount
	if err := tx.Save(&account).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update account balance: %w", err)
	}

	// Create transaction record
	transaction := &models.Transaction{
		AccountID:   accountID,
		Type:        TransactionTypeDeposit,
		Amount:      amount,
		Currency:    account.Currency,
		Description: description,
		Status:      "completed",
	}

	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return transaction, nil
}

// Withdraw removes money from an account
func (s *AccountService) Withdraw(ctx context.Context, accountID uint, amount float64, description string) (*models.Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("withdrawal amount must be greater than zero")
	}

	// Start transaction
	tx := s.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get account with lock
	var account models.Account
	if err := tx.Clauses().Where("id = ? AND status = ?", accountID, "active").First(&account).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("account not found or inactive")
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	// Check if sufficient balance
	if account.Balance < amount {
		tx.Rollback()
		return nil, fmt.Errorf("insufficient balance: available %.2f, requested %.2f", account.Balance, amount)
	}

	// Update balance
	account.Balance -= amount
	if err := tx.Save(&account).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update account balance: %w", err)
	}

	// Create transaction record
	transaction := &models.Transaction{
		AccountID:   accountID,
		Type:        TransactionTypeWithdrawal,
		Amount:      amount,
		Currency:    account.Currency,
		Description: description,
		Status:      "completed",
	}

	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return transaction, nil
}

// GetAccountByID retrieves an account by ID
func (s *AccountService) GetAccountByID(ctx context.Context, accountID uint) (*models.Account, error) {
	var account models.Account
	if err := s.db.WithContext(ctx).First(&account, accountID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("account not found")
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return &account, nil
}

// GetUserAccounts retrieves all accounts for a user
func (s *AccountService) GetUserAccounts(ctx context.Context, userID uint) ([]models.Account, error) {
	var accounts []models.Account
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		return nil, fmt.Errorf("failed to get user accounts: %w", err)
	}

	return accounts, nil
}

// GetAccountBalance retrieves the current balance of an account
func (s *AccountService) GetAccountBalance(ctx context.Context, accountID uint) (float64, error) {
	var account models.Account
	if err := s.db.WithContext(ctx).Select("balance").First(&account, accountID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("account not found")
		}
		return 0, fmt.Errorf("failed to get account balance: %w", err)
	}

	return account.Balance, nil
}

// ValidateAccountOwnership checks if a user owns an account
func (s *AccountService) ValidateAccountOwnership(ctx context.Context, userID, accountID uint) error {
	var account models.Account
	if err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", accountID, userID).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("account not found or you don't have permission to access it")
		}
		return fmt.Errorf("failed to validate account ownership: %w", err)
	}

	return nil
}
