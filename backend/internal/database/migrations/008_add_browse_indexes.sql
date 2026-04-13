-- Migration 008: Add browse performance indexes
-- Optimizes queries for:
-- - ListPublishedPlans (filtered by status, ordered by created_at)
-- - SearchPlans (filtered by destination)

-- Index for browse performance (status, created_at)
CREATE INDEX IF NOT EXISTS idx_travel_plans_status_created_at 
ON travel_plans(status, created_at DESC);

-- Index for destination search
CREATE INDEX IF NOT EXISTS idx_travel_plans_destination 
ON travel_plans(destination);

-- Index for user's plans
CREATE INDEX IF NOT EXISTS idx_travel_plans_author_id 
ON travel_plans(author_id);

-- Index for plan nodes linked list traversal
CREATE INDEX IF NOT EXISTS idx_plan_nodes_plan_sequence 
ON plan_nodes(plan_id, sequence_position);

-- Index for node approval filtering
CREATE INDEX IF NOT EXISTS idx_nodes_type_approved 
ON nodes(type, is_approved);
