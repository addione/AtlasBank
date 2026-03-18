package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/atlasbank/api/internal/config"
	"github.com/atlasbank/api/internal/database"
	"github.com/atlasbank/api/internal/elasticsearch"
	"github.com/atlasbank/api/internal/kafka"
	"github.com/atlasbank/api/internal/redis"
	"github.com/atlasbank/api/internal/router"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger := elasticsearch.NewLogger(cfg.ElasticsearchURL)
	defer logger.Close()

	// Log application startup
	logger.Info("Starting AtlasBank application", map[string]interface{}{
		"port": cfg.AppPort,
		"mode": cfg.GinMode,
	})

	// Initialize database
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		logger.Error("Failed to connect to database", map[string]interface{}{
			"error": err.Error(),
		})
		log.Fatalf("Failed to connect to database: %v", err)
	}
	logger.Info("Database connection established", nil)

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		logger.Error("Failed to run migrations", map[string]interface{}{
			"error": err.Error(),
		})
		log.Fatalf("Failed to run migrations: %v", err)
	}
	logger.Info("Database migrations completed", nil)

	// Initialize Redis
	redisClient := redis.NewRedisClient(cfg)
	defer redisClient.Close()

	// Test Redis connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Error("Failed to connect to Redis", map[string]interface{}{
			"error": err.Error(),
		})
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	logger.Info("Redis connection established", nil)

	// Initialize Kafka producer
	kafkaProducer := kafka.NewProducer(cfg.KafkaBroker)
	defer kafkaProducer.Close()
	logger.Info("Kafka producer initialized", nil)

	// Initialize Kafka consumer (in a separate goroutine)
	kafkaConsumer := kafka.NewConsumer(cfg.KafkaBroker, "atlasbank-group")
	go func() {
		if err := kafkaConsumer.Consume(ctx, logger); err != nil {
			logger.Error("Kafka consumer error", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()
	defer kafkaConsumer.Close()
	logger.Info("Kafka consumer started", nil)

	// Set Gin mode
	gin.SetMode(cfg.GinMode)

	// Initialize router
	r := router.SetupRouter(db, redisClient, kafkaProducer, logger)

	// Start server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.AppPort),
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to start", map[string]interface{}{
				"error": err.Error(),
			})
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	logger.Info("Server started successfully", map[string]interface{}{
		"port": cfg.AppPort,
	})

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...", nil)

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", map[string]interface{}{
			"error": err.Error(),
		})
		log.Fatal("Server forced to shutdown:", err)
	}

	logger.Info("Server exited", nil)
}
