# Code Refactoring Proposal

This document outlines proposed improvements to variable naming and code comments across the AtlasBank codebase.

## Naming Conventions

### Current vs Proposed Variable Names

#### Controllers
| Current | Proposed | Reason |
|---------|----------|--------|
| `ctrl` | `userController` | More descriptive, clearer context |
| `req` | `createUserRequest` | Specific to the request type |
| `err` | `validationError`, `serviceError` | Distinguishes error types |
| `c` | `ginContext` or `ctx` | Clearer that it's Gin's context |

#### Services
| Current | Proposed | Reason |
|---------|----------|--------|
| `s` | `userService` or `service` | More readable in method bodies |
| `ctx` | `requestContext` | Clarifies it's the request context |
| `id` | `userID` | More specific |
| `db` | `database` | Clearer purpose |

#### Models
| Current | Proposed | Reason |
|---------|----------|--------|
| `cfg` | `config` | Full word is clearer |
| `dsn` | `databaseConnectionString` | More descriptive |
| `r` | `router` | Full word preferred |
| `v1` | `apiV1Group` | Clearer that it's a route group |

## Comment Standards

### File-level Comments
Every file should start with:
```go
// Package [name] provides [description]
// 
// This package handles [specific responsibility]
```

### Function Comments
Every exported function should have:
```go
// FunctionName performs [action]
//
// Parameters:
//   - param1: description
//   - param2: description
//
// Returns:
//   - type: description
//   - error: description of possible errors
//
// Example:
//   result, err := FunctionName(param1, param2)
```

### Inline Comments
- Complex logic should have explanatory comments
- Business rules should be documented
- Non-obvious decisions should be explained

## Sample Refactored File

### Before: `internal/controllers/user_controller.go`
```go
func (ctrl *UserController) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}
	// ... rest of code
}
```

### After: `internal/controllers/user_controller.go`
```go
// CreateUser handles the HTTP POST request to create a new user account
//
// This endpoint validates the incoming request, creates a new user with hashed password,
// generates an OTP for account verification, and returns the created user details.
//
// Request Body:
//   - email: User's email address (required, must be valid email format)
//   - first_name: User's first name (required)
//   - last_name: User's last name (required)
//   - password: User's password (required, minimum 6 characters)
//
// Response Codes:
//   - 201: User created successfully
//   - 400: Invalid request body or validation failed
//   - 409: Email already exists
//   - 500: Internal server error
func (userController *UserController) CreateUser(ginContext *gin.Context) {
	var createUserRequest CreateUserRequest

	// Validate and bind the JSON request body to the struct
	if validationError := ginContext.ShouldBindJSON(&createUserRequest); validationError != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": validationError.Error(),
		})
		return
	}

	// Create user through the service layer
	createdUser, serviceError := userController.userService.CreateUser(
		ginContext.Request.Context(),
		createUserRequest.Email,
		createUserRequest.FirstName,
		createUserRequest.LastName,
		createUserRequest.Password,
	)
	
	// Handle service errors
	if serviceError != nil {
		// Check if it's a duplicate email error (business logic error)
		if serviceError.Error() == "email already exists" {
			ginContext.JSON(http.StatusConflict, gin.H{
				"error": "Email already exists",
			})
			return
		}

		// Handle unexpected errors
		ginContext.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create user",
			"details": serviceError.Error(),
		})
		return
	}

	// Return success response with created user
	ginContext.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    createdUser,
	})
}
```

## Proposed Changes Summary

### 1. Controllers (`internal/controllers/`)
- Add comprehensive function documentation
- Rename `ctrl` → `userController`
- Rename `c` → `ginContext`
- Rename `req` → specific request names
- Add inline comments for business logic

### 2. Services (`internal/services/`)
- Add package-level documentation
- Rename `s` → `service` or specific service name
- Add detailed function comments
- Document error conditions
- Add examples for complex functions

### 3. Routes (`internal/routes/`)
- Add comments explaining route grouping
- Document middleware if any
- Explain route patterns

### 4. Models (`internal/database/`)
- Add field-level comments for non-obvious fields
- Document relationships
- Explain validation rules

### 5. Configuration (`internal/config/`)
- Document each configuration field
- Explain default values
- Add usage examples

## Review Instructions

Please review this proposal and:
1. ✅ Approve - Apply all changes
2. ✏️ Modify - Suggest specific changes
3. ❌ Reject - Keep current naming

**Your feedback:**
- Which naming conventions do you prefer?
- Should we use full names (`userController`) or abbreviations (`ctrl`)?
- Any specific files you want to see refactored first?
- Any additional naming standards you'd like to add?

---

**Note:** Once approved, I will:
1. Refactor all files systematically
2. Ensure consistency across the codebase
3. Add comprehensive comments
4. Update documentation to reflect changes
