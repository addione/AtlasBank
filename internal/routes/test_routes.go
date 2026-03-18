package routes

import (
	"net/http"

	"github.com/atlasbank/api/internal/elasticsearch"
	"github.com/atlasbank/api/internal/kafka"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// SetupTestRoutes sets up test endpoints for services
func SetupTestRoutes(router *gin.RouterGroup, redisClient *redis.Client, kafkaProducer *kafka.Producer, logger *elasticsearch.Logger) {
	test := router.Group("/test")
	{
		// Test Redis
		test.GET("/redis", func(c *gin.Context) {
			ctx := c.Request.Context()
			err := redisClient.Set(ctx, "test-key", "test-value", 0).Err()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			val, err := redisClient.Get(ctx, "test-key").Result()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Redis test successful",
				"value":   val,
			})
		})

		// Test Kafka
		test.POST("/kafka", func(c *gin.Context) {
			ctx := c.Request.Context()
			message := map[string]interface{}{
				"event": "test",
				"data":  "Hello from AtlasBank!",
			}

			err := kafkaProducer.SendMessage(ctx, "test", message)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Kafka message sent successfully",
			})
		})

		// Test Elasticsearch logging
		test.POST("/log", func(c *gin.Context) {
			logger.Info("Test log from API", map[string]interface{}{
				"endpoint": "/api/v1/test/log",
				"method":   "POST",
			})

			c.JSON(http.StatusOK, gin.H{
				"message": "Log sent to Elasticsearch",
			})
		})
	}
}
