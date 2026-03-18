package routes

import (
	"github.com/atlasbank/api/internal/handlers"
	"github.com/atlasbank/api/internal/middleware"
	"github.com/atlasbank/api/internal/services"
	"github.com/gin-gonic/gin"
)

// SetupAccountRoutes sets up account-related routes
func SetupAccountRoutes(router *gin.RouterGroup, accountHandler *handlers.AccountHandler, jwtService *services.JWTService) {
	// Protected account routes (require authentication)
	accounts := router.Group("/accounts")
	accounts.Use(middleware.AuthMiddleware(jwtService))
	{
		// Get all user accounts
		accounts.GET("", accountHandler.GetUserAccounts)

		// Get account balance
		accounts.GET("/:account_id/balance", accountHandler.GetBalance)

		// Deposit money
		accounts.POST("/deposit", accountHandler.Deposit)

		// Withdraw money
		accounts.POST("/withdraw", accountHandler.Withdraw)
	}
}
