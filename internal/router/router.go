package router

import (
	"net/http"

	"github.com/atlasbank/api/internal/config"
	"github.com/atlasbank/api/internal/controllers"
	"github.com/atlasbank/api/internal/elasticsearch"
	"github.com/atlasbank/api/internal/handlers"
	"github.com/atlasbank/api/internal/kafka"
	"github.com/atlasbank/api/internal/routes"
	"github.com/atlasbank/api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// SetupRouter sets up the Gin router with all routes
func SetupRouter(db *gorm.DB, redisClient *redis.Client, kafkaProducer *kafka.Producer, logger *elasticsearch.Logger, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// Initialize services
	userService := services.NewUserService(db)
	accountService := services.NewAccountService(db)
	jwtService := services.NewJWTService(cfg.JWTSecret)

	// Initialize controllers
	userController := controllers.NewUserController(userService)

	// Initialize handlers
	accountHandler := handlers.NewAccountHandler(accountService)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "atlasbank-api",
		})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Setup modular routes
		routes.SetupUserRoutes(v1, userController)
		routes.SetupAccountRoutes(v1, accountHandler, jwtService)
		routes.SetupTransactionRoutes(v1)
		routes.SetupTestRoutes(v1, redisClient, kafkaProducer, logger)
	}

	return r
}
