# Database Migrations Guide

This project uses SQL-based migrations to track and manage database schema changes. The migration files are located in `internal/database/migrations/`.

## Migration Files

We have created the following migration files:

### 1. Users Table (000001)
- **Up**: `000001_create_users_table.up.sql`
- **Down**: `000001_create_users_table.down.sql`
- Creates the `users` table with fields: id, email, first_name, last_name, password, timestamps

### 2. Accounts Table (000002)
- **Up**: `000002_create_accounts_table.up.sql`
- **Down**: `000002_create_accounts_table.down.sql`
- Creates the `accounts` table with foreign key to users

### 3. Transactions Table (000003)
- **Up**: `000003_create_transactions_table.up.sql`
- **Down**: `000003_create_transactions_table.down.sql`
- Creates the `transactions` table with foreign key to accounts

## Running Migrations

### Using GORM Auto-Migrate (Current Implementation)

The application currently uses GORM's `AutoMigrate` feature which automatically creates/updates tables based on the model definitions in `internal/database/models.go`.

When you start the application, migrations run automatically:

```bash
make run
# or
go run cmd/api/main.go
```

### Using golang-migrate CLI (Recommended for Production)

For production environments, it's recommended to use the `golang-migrate` CLI tool for more control over migrations.

#### Install golang-migrate

```bash
# macOS
brew install golang-migrate

# Linux
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/

# Windows
scoop install migrate

# Or using Go
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

#### Migration Commands

We've provided a `Makefile.migrations` with helpful commands:

```bash
# Run all pending migrations
make -f Makefile.migrations migrate-up

# Rollback the last migration
make -f Makefile.migrations migrate-down

# Rollback all migrations
make -f Makefile.migrations migrate-down-all

# Create a new migration file
make -f Makefile.migrations migrate-create name=add_user_phone

# Check current migration version
make -f Makefile.migrations migrate-version

# Force set migration version (use with caution!)
make -f Makefile.migrations migrate-force version=1

# Drop everything (DANGEROUS!)
make -f Makefile.migrations migrate-drop
```

#### Manual Migration Commands

```bash
# Set your database URL
export DB_URL="postgresql://atlasbank:atlasbank123@localhost:5432/atlasbank?sslmode=disable"

# Run migrations up
migrate -path internal/database/migrations -database "$DB_URL" up

# Run migrations down
migrate -path internal/database/migrations -database "$DB_URL" down

# Create new migration
migrate create -ext sql -dir internal/database/migrations -seq add_new_field
```

## Migration File Naming Convention

Migration files follow this naming pattern:
```
{version}_{description}.{up|down}.sql
```

Examples:
- `000001_create_users_table.up.sql`
- `000001_create_users_table.down.sql`
- `000002_create_accounts_table.up.sql`
- `000002_create_accounts_table.down.sql`

## Best Practices

1. **Always create both UP and DOWN migrations** - This allows you to rollback changes if needed
2. **Test migrations locally first** - Run migrations on a local database before production
3. **Keep migrations small and focused** - One migration should do one thing
4. **Never modify existing migrations** - Once a migration is committed and deployed, create a new migration to make changes
5. **Use transactions** - Wrap DDL statements in transactions when possible
6. **Backup before migrating** - Always backup production databases before running migrations

## Database Connection

The database connection details are configured in `.env`:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=atlasbank
DB_PASSWORD=atlasbank123
DB_NAME=atlasbank
```

## Troubleshooting

### Migration version mismatch
If you encounter version conflicts:
```bash
# Check current version
make -f Makefile.migrations migrate-version

# Force to a specific version (use carefully!)
make -f Makefile.migrations migrate-force version=3
```

### Dirty database state
If a migration fails midway:
```bash
# Fix the issue in your database manually, then force the version
make -f Makefile.migrations migrate-force version=X
```

### Starting fresh
To completely reset the database:
```bash
# Drop all tables
make -f Makefile.migrations migrate-drop

# Run all migrations again
make -f Makefile.migrations migrate-up
```

## Migration History

The migration system tracks applied migrations in a `schema_migrations` table (golang-migrate) or `migration_histories` table (custom implementation).
