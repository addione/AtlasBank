package middleware

import (
	"net/http"
	"strings"

	"github.com/atlasbank/api/internal/services"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens and adds user info to context
func AuthMiddleware(jwtService *services.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		// Check if it's a Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format. Expected: Bearer <token>",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := jwtService.ValidateAccessToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Add user info to context
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("is_verified", claims.IsVerified)

		c.Next()
	}
}

// RequireVerified ensures the user is verified
func RequireVerified() gin.HandlerFunc {
	return func(c *gin.Context) {
		isVerified, exists := c.Get("is_verified")
		if !exists || !isVerified.(bool) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Account verification required. Please verify your account first.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserID retrieves the user ID from the context
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	return userID.(uint), true
}
