# Database Migrations

This directory contains all database schema migrations for the TravelLink application.

## Overview

- **Framework**: SQLite3 with file-based tracking via `schema_migrations` table
- **Format**: Plain SQL files named with version prefix (e.g., `001_create_users_table.sql`)
- **Execution**: Automatic on application startup; tracks applied migrations
- **Ordering**: Migrations execute in filename order (001, 002, 003, etc.)

## Migration File Format

**Naming Convention**: `<VERSION>_<DESCRIPTION>.sql`

Examples:
- `001_create_users_table.sql`
- `002_create_travel_plans_table.sql`
- `003_add_indexes.sql`

**File Structure**: Each file is pure SQL without any special markers

```sql
-- backend/internal/database/migrations/001_create_users_table.sql
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT DEFAULT 'simple',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
```

## How Migrations Work

1. **On Startup**: When the application starts, `migrations.InitMigrations()` is called
2. **Create Tracking Table**: If `schema_migrations` table doesn't exist, it's created
3. **Load Files**: All `.sql` files in this directory are loaded and sorted by version
4. **Check Status**: For each migration, check if it's listed in `schema_migrations`
5. **Execute Pending**: Run only migrations not yet applied to the database
6. **Record Success**: After execution, add entry to `schema_migrations` with version, name, and timestamp

## Schema Migrations Table

The framework automatically creates:

```sql
CREATE TABLE schema_migrations (
    version TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

Example after 3 migrations:

```
| version | name                      | applied_at          |
|---------|---------------------------|---------------------|
| 001     | create_users_table        | 2026-04-05 10:30:00 |
| 002     | create_travel_plans_table | 2026-04-05 10:30:01 |
| 003     | create_nodes_table        | 2026-04-05 10:30:02 |
```

## Adding a New Migration

1. **Create SQL file** in this directory: `<NEXT_VERSION>_<description>.sql`
2. **Write SQL** using SQLite3 syntax with proper constraints, indexes, foreign keys
3. **Restart application** - migration runs automatically on startup
4. **Verify** - Check logs for "✓ All migrations completed successfully"

Example: Adding a new table

```bash
# Create file: 004_create_comments_table.sql
```

```sql
CREATE TABLE comments (
    id TEXT PRIMARY KEY,
    plan_id TEXT NOT NULL,
    author_id TEXT NOT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (plan_id) REFERENCES travel_plans(id),
    FOREIGN KEY (author_id) REFERENCES users(id)
);

CREATE INDEX idx_comments_plan_id ON comments(plan_id);
CREATE INDEX idx_comments_author_id ON comments(author_id);
```

Then restart the app: `go run cmd/api/main.go`

## Checking Migration Status

Debug migrations programmatically:

```go
status, err := migrations.GetStatus(db, ".")
if err != nil {
    log.Fatal(err)
}

for version, applied := range status {
    if applied {
        log.Printf("✓ %s applied", version)
    } else {
        log.Printf("✗ %s pending", version)
    }
}
```

## Migration Best Practices

### ✅ DO

- ✅ Write idempotent migrations: `CREATE TABLE IF NOT EXISTS`
- ✅ Add constraints: `PRIMARY KEY`, `UNIQUE`, `NOT NULL`, `FOREIGN KEY`
- ✅ Add indexes for frequently queried columns
- ✅ Use meaningful column names and types
- ✅ Include timestamps: `created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP`
- ✅ Test migrations in development before deploying
- ✅ Keep migrations small and focused (one logical change per file)

### ❌ DON'T

- ❌ Don't modify existing migration files (they may have already been applied in production)
- ❌ Don't write app logic in migrations
- ❌ Don't use transactions manually (framework handles this)
- ❌ Don't skip version numbers
- ❌ Don't mix schema changes with data changes in migration files

## SQLite3 Column Types

Common SQLite3 types used in TravelLink:

| Type | Usage | Example |
|------|-------|---------|
| `TEXT` | Strings, UUIDs | `id TEXT PRIMARY KEY` |
| `INTEGER` | Whole numbers, counters | `sequence_position INTEGER` |
| `REAL` | Decimal numbers | `rating REAL` |
| `BOOLEAN` | True/False (stored as 0/1) | `is_active BOOLEAN DEFAULT TRUE` |
| `TIMESTAMP` | Date and time | `created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP` |
| `BLOB` | Binary data | User photos, etc. |

## Troubleshooting

### Migration File Not Running

**Symptoms**: New SQL file created but not executing

**Causes & Fixes**:
1. Filename format wrong - ensure `001_`, `002_`, etc. prefix
2. File not in `migrations/` directory - move to `backend/internal/database/migrations/`
3. `schema_migrations` table corrupted - delete `travellink.db` to reset
4. SQLite database locked - close other connections, ensure WAL mode working

### Migration Fails with SQL Error

**Symptoms**: Application starts but migration shows error

**Causes & Fixes**:
1. SQL syntax error - validate SQL with `sqlite3` CLI tool first
2. Table already exists - use `CREATE TABLE IF NOT EXISTS`
3. Foreign key references non-existent table - ensure reference table created earlier
4. Column type unsupported - verify SQLite3 supports the type

### Database Lock Issues

**Symptoms**: "database is locked" errors

**Causes & Fixes** (SQLite3 single-writer limitation):
1. Multiple processes accessing `travellink.db` - close other terminals/apps
2. WAL mode not enabled - check `PRAGMA journal_mode` returns WAL
3. Connection pool misconfigured - verify `SetMaxOpenConns(1)` in database_sqlite.go

## Related Files

- **Migration Executor**: `backend/internal/database/migrations/runner.go`
- **Public API**: `backend/internal/database/migrations/init.go`
- **SQLite Service**: `backend/internal/database/database_sqlite.go`
- **Configuration**: `.env` file (DB_PATH, DB_TYPE)

## Examples from TravelLink

See Phase 2 tasks (T008-T014) for example migrations:

- T008: `001_create_users_table.sql`
- T009: `002_create_travel_plans_table.sql`
- T010: `003_create_nodes_table.sql`
- T011: `004_create_node_details_tables.sql`
- T012: `005_create_plan_nodes_table.sql`
- T013: `006_create_comments_table.sql`
- T014: `007_create_ratings_table.sql`
