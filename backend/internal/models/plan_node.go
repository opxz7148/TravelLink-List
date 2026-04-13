package models

import "time"

// PlanNode represents the linked list structure of a TravelPlan.
// It is a junction table that associates nodes with travel plans in a specific order.
// Follows the specification from data-model.md with all validation rules.
//
// Each PlanNode is a link in the linked list, pointed to by the next node via sequence_position.
// The sequence ensures plans maintain a coherent order of attractions and transitions.
// Duration fields allow planners to customize how long to spend at each node (per-plan basis).
type PlanNode struct {
	// id (UUID, primary key)
	// Unique identifier for this plan-node association
	ID string `gorm:"primaryKey;type:TEXT" json:"id"`

	// plan_id (UUID, foreign key to TravelPlan)
	// The travel plan this node belongs to
	PlanID string `gorm:"type:TEXT;not null;index:idx_plan_node_seq,priority:1" json:"plan_id"`

	// node_id (UUID, foreign key to Node)
	// The node being added to the travel plan
	NodeID string `gorm:"type:TEXT;not null;index:idx_plan_node_node" json:"node_id"`

	// sequence_position (integer, 1-indexed, required)
	// Position in the linked list (1 = first, 2 = second, etc.)
	// Must be contiguous (no gaps) within a plan
	SequencePosition int `gorm:"type:INTEGER;not null;index:idx_plan_node_seq,priority:2" json:"sequence_position"`

	// description (string, max 500 chars, optional)
	// Plan-specific context and annotations for this node in this plan
	// Independent of node description, allows additional plan-level notes
	// Example: "Try the house special pasta", "Scenic route with views"
	Description *string `gorm:"type:TEXT;size:500" json:"description"`

	// estimated_price_cents (integer, optional, in cents)
	// Cost for this leg in the journey (cents to avoid floating-point issues)
	// 1500 = $15.00, NULL or 0 = free
	// Allows calculating total plan cost from sum of leg prices
	EstimatedPriceCents *int `gorm:"type:INTEGER" json:"estimated_price_cents"`

	// duration_minutes (integer, optional, in minutes)
	// Plan-specific duration for this node in the context of this plan.
	// For attractions: how long the planner wants to spend at the location (e.g., 90 minutes)
	// For transitions: how long the planner expects travel to take (e.g., 30 minutes)
	// Allows flexibility - same node can have different durations across different plans.
	// If not set, frontend can display node's default duration values.
	DurationMinutes *int `gorm:"type:INTEGER" json:"duration_minutes"`

	// created_at (timestamp, UTC, immutable)
	// When this node was added to the plan
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
}

// TableName specifies the database table name for PlanNode
func (PlanNode) TableName() string {
	return "plan_nodes"
}

// Validate performs domain validation on the PlanNode
// Note: Full validation (consecutive attractions, sequence integrity) must be done at service layer
// This validates only the individual field constraints
func (pn *PlanNode) Validate() error {
	// PlanID must not be empty
	if pn.PlanID == "" {
		return ErrValidation
	}

	// NodeID must not be empty
	if pn.NodeID == "" {
		return ErrValidation
	}

	// SequencePosition must be positive (1-indexed, no 0 or negative)
	if pn.SequencePosition <= 0 {
		return ErrValidation
	}

	return nil
}

// PlanNodeList represents a slice of PlanNode for convenience operations
type PlanNodeList []*PlanNode

// Len returns the number of plan nodes
func (pnl PlanNodeList) Len() int {
	return len(pnl)
}

// Less returns true if plan node at index i should sort before plan node at index j
// Used for sorting by sequence position
func (pnl PlanNodeList) Less(i, j int) bool {
	if pnl[i] == nil || pnl[j] == nil {
		return pnl[i] != nil
	}
	return pnl[i].SequencePosition < pnl[j].SequencePosition
}

// Swap swaps the plan nodes at indexes i and j
func (pnl PlanNodeList) Swap(i, j int) {
	pnl[i], pnl[j] = pnl[j], pnl[i]
}

// SequenceGapsExist checks if there are gaps in the sequence positions
// Valid sequence should be: 1, 2, 3, ..., N with no gaps
func (pnl PlanNodeList) SequenceGapsExist() bool {
	if len(pnl) == 0 {
		return false
	}

	// Create a map of all positions for fast lookup
	positions := make(map[int]bool)
	for _, pn := range pnl {
		if pn != nil {
			positions[pn.SequencePosition] = true
		}
	}

	// Check if we have all positions from 1 to len(pnl)
	for i := 1; i <= len(pnl); i++ {
		if !positions[i] {
			return true
		}
	}

	return false
}

// MaxSequencePosition returns the maximum sequence position in the list
func (pnl PlanNodeList) MaxSequencePosition() int {
	max := 0
	for _, pn := range pnl {
		if pn != nil && pn.SequencePosition > max {
			max = pn.SequencePosition
		}
	}
	return max
}
