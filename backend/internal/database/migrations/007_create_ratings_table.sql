-- Migration 007: Create ratings table
-- Represents user ratings of travel plans
-- Supports one rating per user per plan (updates replace previous rating)
-- Maintains denormalization: TravelPlan.rating_count and rating_sum updated on create/update/delete

CREATE TABLE IF NOT EXISTS ratings (
    id TEXT PRIMARY KEY,
    plan_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    stars INTEGER NOT NULL CHECK (stars >= 1 AND stars <= 5),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    
    -- Foreign key constraints
    FOREIGN KEY (plan_id) REFERENCES travel_plans(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    
    -- Unique constraint: one rating per user per plan
    UNIQUE(plan_id, user_id)
);

-- Indexes for common queries
CREATE INDEX IF NOT EXISTS idx_ratings_plan_id ON ratings(plan_id);
CREATE INDEX IF NOT EXISTS idx_ratings_user_id ON ratings(user_id);
CREATE INDEX IF NOT EXISTS idx_ratings_created_at ON ratings(created_at);
