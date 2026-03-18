package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupTransactionRoutes sets up transaction-related routes
func SetupTransactionRoutes(router *gin.RouterGroup) {
	transactions := router.Group("/transactions")
	{
		transactions.GET("", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Get all transactions",
			})
		})
		transactions.POST("", func(c *gin.Context) {
			c.JSON(http.StatusCreated, gin.H{
				"message": "Create transaction",
			})
		})
		transactions.GET("/:id", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Get transaction by ID",
				"id":      c.Param("id"),
			})
		})
	}
}
