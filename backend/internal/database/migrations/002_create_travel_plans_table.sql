-- Migration: Create travel_plans table
-- Represents a complete travel itinerary as a linked list of nodes
-- Status can be draft (private), published (public), suspended (admin action), or deleted (user/admin action)

CREATE TABLE IF NOT EXISTS travel_plans (
    -- Primary key: UUID string
    id TEXT PRIMARY KEY,
    
    -- Title: required, max 150 chars, searchable
    title TEXT NOT NULL CHECK(LENGTH(title) >= 1 AND LENGTH(title) <= 150),
    
    -- Description: optional, max 1000 chars
    description TEXT CHECK(LENGTH(description) IS NULL OR LENGTH(description) <= 1000),
    
    -- Destination: required, max 200 chars, searchable
    -- Used for search/filtering, not strict geocoding
    destination TEXT NOT NULL CHECK(LENGTH(destination) >= 1 AND LENGTH(destination) <= 200),
    
    -- Author ID: foreign key to users table, denormalized for query efficiency
    author_id TEXT NOT NULL,
    
    -- Status: enum value (draft|published|suspended|deleted), default: draft
    -- Controls visibility and lifecycle of the plan
    status TEXT NOT NULL DEFAULT 'draft' CHECK(status IN ('draft', 'published', 'suspended', 'deleted')),
    
    -- Rating count: denormalized count for fast retrieval
    -- Number of ratings received
    rating_count INTEGER NOT NULL DEFAULT 0 CHECK(rating_count >= 0),
    
    -- Rating sum: used to calculate average rating
    -- Sum of all star ratings (average = rating_sum / rating_count when count > 0)
    rating_sum INTEGER NOT NULL DEFAULT 0 CHECK(rating_sum >= 0),
    
    -- Comment count: denormalized count for display
    -- Number of comments on this plan
    comment_count INTEGER NOT NULL DEFAULT 0 CHECK(comment_count >= 0),
    
    -- Is deleted by admin: soft-delete flag
    -- When true, plan is hidden from listings but preserved for history
    is_deleted_by_admin BOOLEAN NOT NULL DEFAULT false,
    
    -- Created at: UTC timestamp, immutable, auto-set on creation
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Updated at: UTC timestamp, auto-updated on modification or node reordering
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key constraint: author must exist in users table
    FOREIGN KEY (author_id) REFERENCES users(id)
);

-- Create indexes for efficient querying
CREATE INDEX IF NOT EXISTS idx_travel_plans_author_id ON travel_plans(author_id);
CREATE INDEX IF NOT EXISTS idx_travel_plans_status ON travel_plans(status);
CREATE INDEX IF NOT EXISTS idx_travel_plans_destination ON travel_plans(destination);
CREATE INDEX IF NOT EXISTS idx_travel_plans_created_at ON travel_plans(created_at);
CREATE INDEX IF NOT EXISTS idx_travel_plans_status_created_at ON travel_plans(status, created_at);
CREATE INDEX IF NOT EXISTS idx_travel_plans_is_deleted_by_admin ON travel_plans(is_deleted_by_admin);
