# Models Refactoring Documentation

## Overview
This document describes the refactoring of database models from a single `internal/database/models.go` file into separate, organized model files in the `internal/models` package.

## Changes Made

### 1. Model Files Created
All models have been separated into individual files in `internal/models/`:

- **`user.go`** - User model with authentication and verification fields
- **`otp.go`** - OTP (One-Time Password) model for user verification
- **`account.go`** - Bank account model with balance and status tracking
- **`transaction.go`** - Financial transaction model with reference tracking
- **`notification.go`** - Notification logging model for email/SMS/push notifications

### 2. Import Path Changes
**Old import:**
```go
import "github.com/atlasbank/api/internal/database"

// Usage
var user database.User
```

**New import:**
```go
import "github.com/atlasbank/api/internal/models"

// Usage
var user models.User
```

### 3. Updated Files
The following files have been updated to use the new `internal/models` package:

- `internal/database/postgres.go` - Updated AutoMigrate to use `models.*`
- `internal/services/user_service.go` - Updated all User references
- `internal/services/otp_service.go` - Updated all OTP references
- `internal/services/notification_service.go` - Updated all Notification references

### 4. Migration Optimizations

#### Users Table (000001)
- Added `is_verified` column directly in initial migration
- Added index on `is_verified` for faster verification queries

#### Accounts Table (000002)
- Added index on `status` for filtering active/frozen/closed accounts

#### Transactions Table (000003)
- Added index on `status` for filtering pending/completed/failed transactions
- Added index on `type` for filtering by transaction type

#### OTP Table (000004)
- Removed redundant `is_verified` column addition (now in users table)
- Added composite index `idx_otps_user_action_used` for optimal OTP lookup queries

#### Notifications Table (000005)
- Added composite index `idx_notifications_user_status` for efficient user notification queries

## Benefits

### 1. Better Code Organization
- Each model is in its own file, making it easier to locate and maintain
- Clear separation of concerns
- Easier to navigate in IDEs

### 2. Improved Maintainability
- Changes to one model don't affect others
- Easier to review changes in version control
- Reduced merge conflicts

### 3. Enhanced Performance
- Optimized database indexes for common query patterns
- Composite indexes for frequently used query combinations
- Better query performance for filtering and lookups

### 4. Scalability
- Easy to add new models without cluttering a single file
- Clear structure for future model additions
- Better support for team collaboration

## Migration Guide

If you have existing code using the old `database.*` models:

1. Update imports:
   ```go
   // Old
   import "github.com/atlasbank/api/internal/database"
   
   // New
   import "github.com/atlasbank/api/internal/models"
   ```

2. Update type references:
   ```go
   // Old
   var user database.User
   var otp database.OTP
   
   // New
   var user models.User
   var otp models.OTP
   ```

3. The old `internal/database/models.go` file has been deprecated with migration notes.

## Database Migration Notes

### For New Installations
Run migrations in order:
```bash
make migrate-up
```

### For Existing Installations
The optimized migrations are backward compatible. However, to benefit from new indexes:

1. Backup your database
2. Run the migration update:
   ```bash
   make migrate-down
   make migrate-up
   ```

Or manually add the new indexes:
```sql
-- Users
CREATE INDEX IF NOT EXISTS idx_users_is_verified ON users(is_verified);

-- Accounts
CREATE INDEX IF NOT EXISTS idx_accounts_status ON accounts(status);

-- Transactions
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);
CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions(type);

-- OTPs
CREATE INDEX IF NOT EXISTS idx_otps_user_action_used ON otps(user_id, action, is_used);

-- Notifications
CREATE INDEX IF NOT EXISTS idx_notifications_user_status ON notifications(user_id, status);
```

## Testing

The refactoring has been validated:
- ✅ All code compiles successfully
- ✅ No breaking changes to existing functionality
- ✅ All imports updated correctly
- ✅ Database migrations optimized

## Future Improvements

Consider these enhancements for future iterations:

1. **Model Validation** - Add validation tags and methods
2. **Custom Types** - Create custom types for status fields (e.g., `AccountStatus`, `TransactionType`)
3. **Model Methods** - Add business logic methods to models
4. **Relationships** - Enhance model relationships with proper eager loading
5. **Soft Delete Scopes** - Add global scopes for soft delete handling

## Questions or Issues?

If you encounter any issues with the refactored models, please check:
1. Import paths are correct (`internal/models` not `internal/database`)
2. All references use `models.*` prefix
3. Database migrations have been run
4. Go modules are up to date (`go mod tidy`)
