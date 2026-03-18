package handlers

import (
	"net/http"
	"strconv"

	"github.com/atlasbank/api/internal/middleware"
	"github.com/atlasbank/api/internal/services"
	"github.com/gin-gonic/gin"
)

// AccountHandler handles account-related HTTP requests
type AccountHandler struct {
	accountService *services.AccountService
}

// NewAccountHandler creates a new AccountHandler
func NewAccountHandler(accountService *services.AccountService) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
	}
}

// DepositRequest represents a deposit request
type DepositRequest struct {
	AccountID   uint    `json:"account_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description"`
}

// WithdrawalRequest represents a withdrawal request
type WithdrawalRequest struct {
	AccountID   uint    `json:"account_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description"`
}

// Deposit handles deposit requests
// @Summary Deposit money into account
// @Description Deposit money into a user's account (requires authentication)
// @Tags accounts
// @Accept json
// @Produce json
// @Param request body DepositRequest true "Deposit request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /accounts/deposit [post]
// @Security BearerAuth
func (h *AccountHandler) Deposit(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var req DepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	// Validate account ownership
	if err := h.accountService.ValidateAccountOwnership(c.Request.Context(), userID, req.AccountID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Perform deposit
	transaction, err := h.accountService.Deposit(c.Request.Context(), req.AccountID, req.Amount, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process deposit: " + err.Error(),
		})
		return
	}

	// Get updated balance
	balance, err := h.accountService.GetAccountBalance(c.Request.Context(), req.AccountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Deposit successful but failed to retrieve balance: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Deposit successful",
		"transaction": transaction,
		"new_balance": balance,
	})
}

// Withdraw handles withdrawal requests
// @Summary Withdraw money from account
// @Description Withdraw money from a user's account (requires authentication and sufficient balance)
// @Tags accounts
// @Accept json
// @Produce json
// @Param request body WithdrawalRequest true "Withdrawal request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /accounts/withdraw [post]
// @Security BearerAuth
func (h *AccountHandler) Withdraw(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var req WithdrawalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	// Validate account ownership
	if err := h.accountService.ValidateAccountOwnership(c.Request.Context(), userID, req.AccountID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Perform withdrawal
	transaction, err := h.accountService.Withdraw(c.Request.Context(), req.AccountID, req.Amount, req.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Get updated balance
	balance, err := h.accountService.GetAccountBalance(c.Request.Context(), req.AccountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Withdrawal successful but failed to retrieve balance: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Withdrawal successful",
		"transaction": transaction,
		"new_balance": balance,
	})
}

// GetBalance retrieves the account balance
// @Summary Get account balance
// @Description Get the current balance of a user's account
// @Tags accounts
// @Produce json
// @Param account_id path int true "Account ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /accounts/{account_id}/balance [get]
// @Security BearerAuth
func (h *AccountHandler) GetBalance(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Get account ID from URL parameter
	accountIDStr := c.Param("account_id")
	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid account ID",
		})
		return
	}

	// Validate account ownership
	if err := h.accountService.ValidateAccountOwnership(c.Request.Context(), userID, uint(accountID)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Get balance
	balance, err := h.accountService.GetAccountBalance(c.Request.Context(), uint(accountID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"account_id": accountID,
		"balance":    balance,
	})
}

// GetUserAccounts retrieves all accounts for the authenticated user
// @Summary Get user accounts
// @Description Get all accounts belonging to the authenticated user
// @Tags accounts
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /accounts [get]
// @Security BearerAuth
func (h *AccountHandler) GetUserAccounts(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Get user accounts
	accounts, err := h.accountService.GetUserAccounts(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve accounts: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accounts": accounts,
		"count":    len(accounts),
	})
}
