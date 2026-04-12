-- Migration 006: Create comments table
-- Represents user commentary on travel plans
-- Supports soft-delete via is_deleted_by_admin flag
-- Maintains denormalization: TravelPlan.comment_count updated on create/delete

CREATE TABLE IF NOT EXISTS comments (
    id TEXT PRIMARY KEY,
    plan_id TEXT NOT NULL,
    author_id TEXT NOT NULL,
    text TEXT NOT NULL CHECK (LENGTH(text) >= 1 AND LENGTH(text) <= 1000),
    is_deleted_by_admin BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    
    -- Foreign key constraints
    FOREIGN KEY (plan_id) REFERENCES travel_plans(id) ON DELETE CASCADE,
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Indexes for common queries
CREATE INDEX IF NOT EXISTS idx_comments_plan_id ON comments(plan_id);
CREATE INDEX IF NOT EXISTS idx_comments_author_id ON comments(author_id);
CREATE INDEX IF NOT EXISTS idx_comments_created_at ON comments(created_at);
CREATE INDEX IF NOT EXISTS idx_comments_not_deleted ON comments(is_deleted_by_admin, plan_id);
