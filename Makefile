.PHONY: help build run stop clean logs test docker-up docker-down docker-rebuild

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the Go application
	@echo "Building application..."
	go build -o atlasbank ./cmd/api

run: ## Run the application locally
	@echo "Running application..."
	go run ./cmd/api/main.go

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

docker-up: ## Start all Docker services
	@echo "Starting Docker services..."
	docker-compose up -d

docker-down: ## Stop all Docker services
	@echo "Stopping Docker services..."
	docker-compose down

docker-rebuild: ## Rebuild and restart Docker services
	@echo "Rebuilding Docker services..."
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d

logs: ## View Docker logs
	docker-compose logs -f

logs-app: ## View application logs
	docker-compose logs -f app

logs-db: ## View database logs
	docker-compose logs -f postgres

logs-kafka: ## View Kafka logs
	docker-compose logs -f kafka

clean: ## Clean build artifacts and Docker volumes
	@echo "Cleaning..."
	rm -f atlasbank
	docker-compose down -v

tidy: ## Tidy Go modules
	@echo "Tidying Go modules..."
	GOPROXY=direct GOSUMDB=off go mod tidy

fmt: ## Format Go code
	@echo "Formatting code..."
	go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

lint: fmt vet ## Run formatters and linters

health: ## Check application health
	@echo "Checking application health..."
	curl -s http://localhost:8081/health | jq .

test-redis: ## Test Redis connection
	@echo "Testing Redis..."
	curl -s http://localhost:8081/api/v1/test/redis | jq .

test-kafka: ## Test Kafka connection
	@echo "Testing Kafka..."
	curl -s -X POST http://localhost:8081/api/v1/test/kafka | jq .

test-log: ## Test Elasticsearch logging
	@echo "Testing Elasticsearch logging..."
	curl -s -X POST http://localhost:8081/api/v1/test/log | jq .

ps: ## Show running containers
	docker-compose ps

restart: ## Restart all services
	docker-compose restart

restart-app: ## Restart only the application
	docker-compose restart app
