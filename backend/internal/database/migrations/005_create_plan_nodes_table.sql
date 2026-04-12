-- Migration: Create plan_nodes table (linked list association)
-- Junction table representing the linked list structure of a TravelPlan
-- Each node maintains a sequence position that ensures contiguous ordering

CREATE TABLE IF NOT EXISTS plan_nodes (
    -- Primary key: UUID string
    id TEXT PRIMARY KEY,
    
    -- Plan ID: foreign key to travel_plans table
    -- The travel plan this node belongs to
    plan_id TEXT NOT NULL,
    
    -- Node ID: foreign key to nodes table
    -- The node being added to the travel plan
    node_id TEXT NOT NULL,
    
    -- Sequence position: 1-indexed position in the linked list
    -- Position must be unique within a plan and form contiguous sequence (1..N with no gaps)
    sequence_position INTEGER NOT NULL CHECK(sequence_position > 0),
    
    -- Created at: UTC timestamp, immutable, auto-set on creation
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key constraints
    FOREIGN KEY (plan_id) REFERENCES travel_plans(id) ON DELETE CASCADE,
    FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE CASCADE,
    
    -- Composite unique constraints
    -- 1. Each node can appear only once in a plan
    UNIQUE(plan_id, node_id),
    -- 2. Sequence position must be unique within each plan (no duplicates)
    UNIQUE(plan_id, sequence_position)
);

-- Create indexes for efficient querying and linked list traversal
-- Plan lookup with ordering by sequence
CREATE INDEX IF NOT EXISTS idx_plan_nodes_plan_id ON plan_nodes(plan_id);
CREATE INDEX IF NOT EXISTS idx_plan_nodes_plan_seq ON plan_nodes(plan_id, sequence_position);
CREATE INDEX IF NOT EXISTS idx_plan_nodes_node_id ON plan_nodes(node_id);
CREATE INDEX IF NOT EXISTS idx_plan_nodes_created_at ON plan_nodes(created_at);
