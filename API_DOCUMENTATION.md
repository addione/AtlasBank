# AtlasBank API Documentation

## Base URL
```
http://localhost:8080
```

## Health Check

### GET /health
Check if the API is running.

**Response:**
```json
{
  "status": "healthy",
  "service": "atlasbank-api"
}
```

---

## User Endpoints

### POST /api/v1/users
Create a new user.

**Request Body:**
```json
{
  "email": "user@example.com",
  "first_name": "John",
  "last_name": "Doe",
  "password": "securepassword123"
}
```

**Validation Rules:**
- `email`: Required, must be a valid email format
- `first_name`: Required
- `last_name`: Required
- `password`: Required, minimum 6 characters

**Success Response (201 Created):**
```json
{
  "message": "User created successfully",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "created_at": "2026-03-18T14:30:00Z",
    "updated_at": "2026-03-18T14:30:00Z"
  }
}
```

**Error Responses:**

400 Bad Request - Invalid input:
```json
{
  "error": "Invalid request",
  "details": "Key: 'CreateUserRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag"
}
```

409 Conflict - Email already exists:
```json
{
  "error": "Email already exists"
}
```

500 Internal Server Error:
```json
{
  "error": "Failed to create user",
  "details": "error message"
}
```

**Example cURL:**
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "password": "password123"
  }'
```

---

### GET /api/v1/users/:id
Get a user by ID.

**URL Parameters:**
- `id` (integer): User ID

**Success Response (200 OK):**
```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "created_at": "2026-03-18T14:30:00Z",
    "updated_at": "2026-03-18T14:30:00Z"
  }
}
```

**Error Responses:**

404 Not Found:
```json
{
  "error": "User not found"
}
```

500 Internal Server Error:
```json
{
  "error": "Failed to get user",
  "details": "error message"
}
```

**Example cURL:**
```bash
curl http://localhost:8080/api/v1/users/1
```

---

### GET /api/v1/users
Get all users.

**Success Response (200 OK):**
```json
{
  "users": [
    {
      "id": 1,
      "email": "user1@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "created_at": "2026-03-18T14:30:00Z",
      "updated_at": "2026-03-18T14:30:00Z"
    },
    {
      "id": 2,
      "email": "user2@example.com",
      "first_name": "Jane",
      "last_name": "Smith",
      "created_at": "2026-03-18T14:35:00Z",
      "updated_at": "2026-03-18T14:35:00Z"
    }
  ],
  "count": 2
}
```

**Error Response:**

500 Internal Server Error:
```json
{
  "error": "Failed to get users",
  "details": "error message"
}
```

**Example cURL:**
```bash
curl http://localhost:8080/api/v1/users
```

---

## Account Endpoints (Placeholder)

### GET /api/v1/accounts
Get all accounts (not yet implemented).

### POST /api/v1/accounts
Create an account (not yet implemented).

### GET /api/v1/accounts/:id
Get account by ID (not yet implemented).

---

## Transaction Endpoints (Placeholder)

### GET /api/v1/transactions
Get all transactions (not yet implemented).

### POST /api/v1/transactions
Create a transaction (not yet implemented).

### GET /api/v1/transactions/:id
Get transaction by ID (not yet implemented).

---

## Test Endpoints

### GET /api/v1/test/redis
Test Redis connection.

**Success Response (200 OK):**
```json
{
  "message": "Redis test successful",
  "value": "test-value"
}
```

---

### POST /api/v1/test/kafka
Test Kafka producer.

**Success Response (200 OK):**
```json
{
  "message": "Kafka message sent successfully"
}
```

---

### POST /api/v1/test/log
Test Elasticsearch logging.

**Success Response (200 OK):**
```json
{
  "message": "Log sent to Elasticsearch"
}
```

---

## Project Structure

```
internal/
├── controllers/       # HTTP request handlers
│   └── user_controller.go
├── services/         # Business logic layer
│   └── user_service.go
├── routes/           # Route definitions
│   ├── user_routes.go
│   ├── account_routes.go
│   ├── transaction_routes.go
│   └── test_routes.go
├── database/         # Database models and migrations
│   ├── models.go
│   ├── postgres.go
│   └── migrations/
├── router/           # Main router setup
│   └── router.go
└── ...
```

## Security Notes

- Passwords are hashed using bcrypt before storage
- Password fields are never returned in API responses (marked with `json:"-"`)
- Email uniqueness is enforced at the database level
- Input validation is performed using Gin's binding tags

## Running the API

```bash
# Start all services (PostgreSQL, Redis, Kafka, Elasticsearch)
docker-compose up -d

# Run the API
make run
# or
go run cmd/api/main.go
```

## Testing the API

You can use the provided test endpoints or tools like:
- cURL (command line)
- Postman
- Thunder Client (VS Code extension)
- HTTPie

Example test flow:
```bash
# 1. Check health
curl http://localhost:8080/health

# 2. Create a user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","first_name":"Test","last_name":"User","password":"password123"}'

# 3. Get all users
curl http://localhost:8080/api/v1/users

# 4. Get specific user
curl http://localhost:8080/api/v1/users/1
```
