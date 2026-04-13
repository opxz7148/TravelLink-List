-- Migration: Add new fields to transition_node_details and plan_nodes
-- - transition_node_details: title (service/line identifier), hours_of_operation
-- - plan_nodes: description (plan-specific context), estimated_price_cents (cost per leg)
-- Date: 2026-04-13

-- ============================================================================
-- ALTER TRANSITION_NODE_DETAILS TABLE
-- ============================================================================

-- Add title field to transition_node_details (e.g., "Bus Line 5", "M1 Train", "Walking")
-- This identifies the specific service/line and is immutable
ALTER TABLE transition_node_details ADD COLUMN title TEXT NOT NULL DEFAULT '' CHECK(length(title) <= 200);

-- Add hours_of_operation field (e.g., "Mon-Fri 6:00-23:00, Sat-Sun 7:00-22:00")
-- Applies to the service itself (e.g., public transit schedules)
-- Null for modes without fixed hours (walking, driving)
ALTER TABLE transition_node_details ADD COLUMN hours_of_operation TEXT CHECK(length(hours_of_operation) <= 200);

-- ============================================================================
-- ALTER PLAN_NODES TABLE
-- ============================================================================

-- Add description field for plan-specific context
-- Independent of the node's description, allows authors to add plan-level annotations
-- Example: "Try the house special pasta", "Scenic route with views"
ALTER TABLE plan_nodes ADD COLUMN description TEXT CHECK(length(description) <= 500);

-- Add estimated_price_cents field for cost per node in this plan
-- Stored in cents to avoid floating-point issues (e.g., 1500 = $15.00)
-- Optional for free transitions (walking) or attractions without pricing
-- NULL or 0 means free
ALTER TABLE plan_nodes ADD COLUMN estimated_price_cents INTEGER CHECK(estimated_price_cents IS NULL OR estimated_price_cents > 0);

-- Create indexes for new fields to support queries
CREATE INDEX IF NOT EXISTS idx_plan_nodes_price ON plan_nodes(estimated_price_cents) WHERE estimated_price_cents IS NOT NULL;
