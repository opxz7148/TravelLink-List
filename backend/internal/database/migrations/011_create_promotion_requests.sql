-- Migration 010: Create promotion_requests table
-- Represents promotions requests for user role upgrades or plan visibility

CREATE TABLE IF NOT EXISTS promotion_requests (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL UNIQUE,
    plan_id TEXT,
    status TEXT NOT NULL DEFAULT 'pending',
    admin_notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    reviewed_at TIMESTAMP,
    UNIQUE(user_id, plan_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (plan_id) REFERENCES travel_plans(id) ON DELETE SET NULL,
    CHECK (status IN ('pending', 'approved', 'rejected'))
);

-- Indexes for queries
CREATE INDEX IF NOT EXISTS idx_promotion_requests_user_id 
ON promotion_requests(user_id);

CREATE INDEX IF NOT EXISTS idx_promotion_requests_plan_id 
ON promotion_requests(plan_id);

CREATE INDEX IF NOT EXISTS idx_promotion_requests_status 
ON promotion_requests(status);

CREATE INDEX IF NOT EXISTS idx_promotion_requests_created_at 
ON promotion_requests(created_at DESC);
