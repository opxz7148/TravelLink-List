-- 001_create_users_table.sql
-- Creates the users table as per data-model.md specification
-- User entity with 11 fields and role-based access control

CREATE TABLE IF NOT EXISTS users (
    -- Primary key: UUID v4 format stored as TEXT
    id TEXT PRIMARY KEY,

    -- Email: unique, required, RFC 5322 validated
    email TEXT UNIQUE NOT NULL,

    -- Username: unique, required, 3-50 alphanumeric + underscore
    username TEXT UNIQUE NOT NULL,

    -- Password hash: bcrypt hash only, never plaintext, required
    password_hash TEXT NOT NULL,

    -- Role: enum value (simple|traveller|admin), default: simple
    -- Defines user permissions:
    --   - simple: Browse, search, comment, rate (default)
    --   - traveller: Create/edit travel plans, create nodes (after promotion)
    --   - admin: Full moderation access
    role TEXT DEFAULT 'simple' NOT NULL,

    -- Display name: optional, max 100 chars, defaults to username in UI
    display_name TEXT,

    -- Bio: optional, max 500 chars
    bio TEXT,

    -- Profile picture URL: optional, stored as string (CDN URL or data URL)
    profile_picture_url TEXT,

    -- Is active: boolean soft-delete flag, default: true
    -- When false, user is hidden from public listings but content preserved
    is_active BOOLEAN DEFAULT 1 NOT NULL,

    -- Created at: UTC timestamp, immutable, auto-set on creation
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,

    -- Updated at: UTC timestamp, auto-updated on modification
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,

    -- Last login: nullable UTC timestamp, tracks successful authentication
    last_login_at TIMESTAMP,

    -- Indexes for query optimization
    UNIQUE INDEX idx_users_email ON email,
    UNIQUE INDEX idx_users_username ON username,
    INDEX idx_users_role ON role,
    INDEX idx_users_is_active ON is_active,
    INDEX idx_users_last_login_at ON last_login_at
);

-- Create trigger to auto-update updated_at timestamp on row modification
CREATE TRIGGER IF NOT EXISTS users_updated_at_trigger
AFTER UPDATE ON users
FOR EACH ROW
BEGIN
    UPDATE users SET updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.id;
END;
