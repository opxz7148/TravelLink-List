-- Migration: Add deleted status to travel_plans table
-- Updates the CHECK constraint to include the new 'deleted' status
-- Allows plans to be marked as deleted instead of using is_deleted_by_admin flag

-- SQLite doesn't support direct ALTER TABLE constraint modification
-- Solution: Drop the constraint implicitly by recreating it via PRAGMA foreign_keys

-- For SQLite, we need to:
-- 1. Disable foreign key checks temporarily
-- 2. Rename old table to backup
-- 3. Create new table with updated constraint
-- 4. Copy data from backup
-- 5. Drop backup
-- 6. Re-enable foreign key checks

PRAGMA foreign_keys = OFF;

-- Create temporary table with updated schema
CREATE TABLE travel_plans_new (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL CHECK(LENGTH(title) >= 1 AND LENGTH(title) <= 150),
    description TEXT CHECK(LENGTH(description) IS NULL OR LENGTH(description) <= 1000),
    destination TEXT NOT NULL CHECK(LENGTH(destination) >= 1 AND LENGTH(destination) <= 200),
    author_id TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'draft' CHECK(status IN ('draft', 'published', 'suspended', 'deleted')),
    rating_count INTEGER NOT NULL DEFAULT 0 CHECK(rating_count >= 0),
    rating_sum INTEGER NOT NULL DEFAULT 0 CHECK(rating_sum >= 0),
    comment_count INTEGER NOT NULL DEFAULT 0 CHECK(comment_count >= 0),
    is_deleted_by_admin BOOLEAN NOT NULL DEFAULT false,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (author_id) REFERENCES users(id)
);

-- Copy all data from old table to new table
INSERT INTO travel_plans_new SELECT * FROM travel_plans;

-- Drop old table
DROP TABLE travel_plans;

-- Rename new table to original name
ALTER TABLE travel_plans_new RENAME TO travel_plans;

-- Re-enable foreign key checks
PRAGMA foreign_keys = ON;
