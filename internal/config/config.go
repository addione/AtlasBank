package config

import (
	"os"
)

// Config holds all application configuration
type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Redis
	RedisHost string
	RedisPort string

	// Kafka
	KafkaBroker string

	// Elasticsearch
	ElasticsearchURL string

	// Application
	AppPort string
	GinMode string

	// JWT
	JWTSecret string
}

// Load reads configuration from environment variables
func Load() *Config {
	return &Config{
		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "atlasbank"),
		DBPassword: getEnv("DB_PASSWORD", "atlasbank_password"),
		DBName:     getEnv("DB_NAME", "atlasbank_db"),

		// Redis
		RedisHost: getEnv("REDIS_HOST", "localhost"),
		RedisPort: getEnv("REDIS_PORT", "6379"),

		// Kafka
		KafkaBroker: getEnv("KAFKA_BROKER", "localhost:9092"),

		// Elasticsearch
		ElasticsearchURL: getEnv("ELASTICSEARCH_URL", "http://localhost:9200"),

		// Application
		AppPort: getEnv("APP_PORT", "8080"),
		GinMode: getEnv("GIN_MODE", "debug"),

		// JWT
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key-change-this-in-production"),
	}
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
