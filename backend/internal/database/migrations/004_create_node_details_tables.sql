-- Migration: Create attraction and transition node detail tables
-- Stores subtype-specific fields for polymorphic nodes

-- Attraction node details table
CREATE TABLE IF NOT EXISTS attraction_node_details (
    -- Foreign key to nodes table
    node_id TEXT PRIMARY KEY,
    
    -- Name of the attraction (required, max 200 chars)
    name TEXT NOT NULL CHECK(length(name) > 0 AND length(name) <= 200),
    
    -- Category of attraction (fixed set of values)
    category TEXT NOT NULL CHECK(category IN (
        'tourist_attraction', 'restaurant', 'hotel', 'museum', 
        'park', 'shopping', 'entertainment', 'other'
    )),
    
    -- Location/address (required, max 300 chars)
    location TEXT NOT NULL CHECK(length(location) > 0 AND length(location) <= 300),
    
    -- Detailed description (optional, max 1000 chars)
    description TEXT CHECK(length(description) <= 1000),
    
    -- Contact information: phone, email, or website (optional, max 200 chars)
    contact_info TEXT CHECK(length(contact_info) <= 200),
    
    -- Operating hours - human readable or ISO 8601 (optional, max 200 chars)
    hours_of_operation TEXT CHECK(length(hours_of_operation) <= 200),
    
    -- Estimated time to spend at this location in minutes (optional)
    estimated_visit_duration_minutes INTEGER CHECK(estimated_visit_duration_minutes > 0),
    
    -- Timestamp
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key constraint
    FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE CASCADE ON UPDATE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_attraction_category ON attraction_node_details(category);

-- Transition node details table
CREATE TABLE IF NOT EXISTS transition_node_details (
    -- Foreign key to nodes table
    node_id TEXT PRIMARY KEY,
    
    -- Mode of transportation (required)
    mode TEXT NOT NULL CHECK(mode IN (
        'walking', 'car', 'bus', 'train', 'bike', 'taxi', 'flight', 'other'
    )),
    
    -- Estimated duration in minutes (required, > 0)
    estimated_duration_minutes INTEGER NOT NULL CHECK(estimated_duration_minutes > 0),
    
    -- Route notes/directions (optional, max 500 chars)
    route_notes TEXT CHECK(length(route_notes) <= 500),
    
    -- Estimated distance in kilometers (optional, >= 0)
    estimated_distance_km REAL CHECK(estimated_distance_km >= 0),
    
    -- Timestamp
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key constraint
    FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE CASCADE ON UPDATE CASCADE
);
