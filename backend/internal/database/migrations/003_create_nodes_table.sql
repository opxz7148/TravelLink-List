-- Migration: Create nodes table (base table for polymorphic inheritance)
-- This table stores both attraction and transition nodes with a type discriminator
-- Specific details are stored in separate detail tables

CREATE TABLE IF NOT EXISTS nodes (
    -- Primary key: UUID string
    id TEXT PRIMARY KEY,
    
    -- Type discriminator: 'attraction' or 'transition'
    type TEXT NOT NULL CHECK(type IN ('attraction', 'transition')),
    
    -- Foreign key to users table (creator for moderation context)
    created_by TEXT NOT NULL,
    
    -- Whether node is approved by admin
    -- User-created nodes start as false, system nodes start as true
    is_approved BOOLEAN NOT NULL DEFAULT true,
    
    -- Timestamps
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- Updated only when admin changes approval status
    updated_at DATETIME,
    
    -- Indexes for efficient querying
    FOREIGN KEY (created_by) REFERENCES users(id)
);

-- Create indexes for efficient filtering
CREATE INDEX IF NOT EXISTS idx_nodes_type ON nodes(type);
CREATE INDEX IF NOT EXISTS idx_nodes_created_by ON nodes(created_by);
CREATE INDEX IF NOT EXISTS idx_nodes_is_approved ON nodes(is_approved);
CREATE INDEX IF NOT EXISTS idx_nodes_created_at ON nodes(created_at);
