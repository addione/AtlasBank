# AtlasBank - Banking Application

A microservices-based banking application built with Go (Golang) using the Gin framework, featuring PostgreSQL, Redis, Kafka, and Elasticsearch.

## Architecture

This application uses a Docker-based architecture with the following components:

- **Main Application**: Go application using Gin framework (Port: 8081)
- **PostgreSQL**: Database for persistent storage (Port: 5433)
- **Redis**: Caching layer (Port: 6380)
- **Kafka**: Message queue for event-driven architecture (Port: 9093)
- **Zookeeper**: Required for Kafka coordination (Port: 2182)
- **Elasticsearch**: Centralized logging (Port: 9201)
- **Kibana**: Log visualization dashboard (Port: 5602)

## Prerequisites

- Docker and Docker Compose installed
- Go 1.21 or higher (for local development)
- At least 4GB of available RAM for Docker

## Project Structure

```
AtlasBank/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go              # Configuration management
│   ├── database/
│   │   ├── models.go              # Database models
│   │   └── postgres.go            # PostgreSQL connection
│   ├── elasticsearch/
│   │   └── logger.go              # Elasticsearch logger
│   ├── kafka/
│   │   ├── consumer.go            # Kafka consumer
│   │   └── producer.go            # Kafka producer
│   ├── redis/
│   │   └── redis.go               # Redis client
│   └── router/
│       └── router.go              # API routes
├── docker-compose.yml             # Docker services configuration
├── Dockerfile                     # Application container
├── go.mod                         # Go module dependencies
├── go.sum                         # Go module checksums
└── README.md                      # This file
```

## Getting Started

### 1. Clone the repository

```bash
cd /home/akash/personal/AtlasBank
```

### 2. Start all services

```bash
docker-compose up -d
```

This will start all services in the background. The first run will take a few minutes to download images and build the application.

### 3. Check service health

```bash
docker-compose ps
```

All services should show as "healthy" or "running".

### 4. View logs

```bash
# View all logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f app
docker-compose logs -f postgres
docker-compose logs -f kafka
```

## API Endpoints

### Health Check
- `GET /health` - Check if the service is running

### Users
- `GET /api/v1/users` - Get all users
- `POST /api/v1/users` - Create a new user
- `GET /api/v1/users/:id` - Get user by ID

### Accounts
- `GET /api/v1/accounts` - Get all accounts
- `POST /api/v1/accounts` - Create a new account
- `GET /api/v1/accounts/:id` - Get account by ID

### Transactions
- `GET /api/v1/transactions` - Get all transactions
- `POST /api/v1/transactions` - Create a new transaction
- `GET /api/v1/transactions/:id` - Get transaction by ID

### Test Endpoints
- `GET /api/v1/test/redis` - Test Redis connection
- `POST /api/v1/test/kafka` - Test Kafka message publishing
- `POST /api/v1/test/log` - Test Elasticsearch logging

## Testing the Application

### Test the health endpoint
```bash
curl http://localhost:8081/health
```

### Test Redis
```bash
curl http://localhost:8081/api/v1/test/redis
```

### Test Kafka
```bash
curl -X POST http://localhost:8081/api/v1/test/kafka
```

### Test Elasticsearch logging
```bash
curl -X POST http://localhost:8081/api/v1/test/log
```

## Accessing Services

- **Application API**: http://localhost:8081
- **PostgreSQL**: localhost:5433 (user: atlasbank, password: atlasbank_password, db: atlasbank_db)
- **Redis**: localhost:6380
- **Kafka**: localhost:9093
- **Elasticsearch**: http://localhost:9201
- **Kibana**: http://localhost:5602

## Database Models

### User
- ID, Email, FirstName, LastName, Password
- One-to-many relationship with Accounts

### Account
- ID, UserID, AccountNumber, AccountType, Balance, Currency, Status
- Belongs to User
- One-to-many relationship with Transactions

### Transaction
- ID, AccountID, Type, Amount, Currency, Description, Status, ReferenceNumber
- Belongs to Account

## Development

### Local Development (without Docker)

1. Install dependencies:
```bash
go mod download
```

2. Set environment variables:
```bash
export DB_HOST=localhost
export DB_PORT=5433
export REDIS_HOST=localhost
export REDIS_PORT=6380
export KAFKA_BROKER=localhost:9093
export ELASTICSEARCH_URL=http://localhost:9201
export APP_PORT=8080
export GIN_MODE=debug
```

3. Run the application:
```bash
go run cmd/api/main.go
```

### Building the application

```bash
go build -o atlasbank ./cmd/api
```

## Stopping the Application

```bash
# Stop all services
docker-compose down

# Stop and remove volumes (WARNING: This will delete all data)
docker-compose down -v
```

## Troubleshooting

### Services not starting
- Check if ports are already in use: `netstat -tulpn | grep -E '8081|5433|6380|9093|9201|5602'`
- Check Docker logs: `docker-compose logs`

### Database connection issues
- Ensure PostgreSQL is healthy: `docker-compose ps postgres`
- Check database logs: `docker-compose logs postgres`

### Kafka connection issues
- Ensure both Zookeeper and Kafka are running
- Check Kafka logs: `docker-compose logs kafka`

### Application not building
- Run `go mod tidy` to clean up dependencies
- Ensure Go 1.21+ is installed

## Port Configuration

All external ports are configured to avoid conflicts with other applications:

| Service | Internal Port | External Port |
|---------|--------------|---------------|
| App | 8080 | 8081 |
| PostgreSQL | 5432 | 5433 |
| Redis | 6379 | 6380 |
| Zookeeper | 2181 | 2182 |
| Kafka | 9092 | 9093 |
| Elasticsearch | 9200 | 9201 |
| Elasticsearch | 9300 | 9301 |
| Kibana | 5601 | 5602 |

## Environment Variables

The application uses the following environment variables (configured in docker-compose.yml):

- `DB_HOST`: PostgreSQL host
- `DB_PORT`: PostgreSQL port
- `DB_USER`: PostgreSQL user
- `DB_PASSWORD`: PostgreSQL password
- `DB_NAME`: PostgreSQL database name
- `REDIS_HOST`: Redis host
- `REDIS_PORT`: Redis port
- `KAFKA_BROKER`: Kafka broker address
- `ELASTICSEARCH_URL`: Elasticsearch URL
- `APP_PORT`: Application port
- `GIN_MODE`: Gin framework mode (debug/release)

## License

This is a sample project for educational purposes.
