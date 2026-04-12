package services

import (
	"context"
	"time"

	"tll-backend/internal/models"
	"tll-backend/internal/repositories"

	"github.com/google/uuid"
)

// PlanService defines the interface for travel plan business logic operations
// Handles creation, modification, querying, and lifecycle management of travel plans
// Also manages the linked-list structure of nodes within plans
type PlanService interface {
	// TravelPlan CRUD Operations

	// CreatePlan creates a new travel plan
	// Validates plan according to data-model.md rules
	// Returns plan ID on success
	// Returns ErrValidation if validation fails
	CreatePlan(ctx context.Context, plan *models.TravelPlan) (string, error)

	// GetPlanByID retrieves a travel plan by its ID
	// Does NOT perform authorization - caller must check CanBeViewedBy
	// Returns nil plan if not found
	GetPlanByID(ctx context.Context, planID string) (*models.TravelPlan, error)

	// GetPlansByAuthor retrieves all travel plans authored by a specific user
	// Includes draft, published, and suspended plans
	// Paginated: offset (0-based), limit (max results)
	GetPlansByAuthor(ctx context.Context, authorID string, offset, limit int) ([]*models.TravelPlan, int, error)

	// ListPublishedPlans retrieves all published travel plans (excluding admin-deleted)
	// Used for browsing and discovery
	// Paginated: offset (0-based), limit (max results)
	// Returns plans, total count, error
	ListPublishedPlans(ctx context.Context, offset, limit int) ([]*models.TravelPlan, int, error)

	// SearchPlans searches travel plans by destination or title (published only)
	// Used for user search/discovery
	// Paginated: offset (0-based), limit (max results)
	// Returns plans, total count, error
	SearchPlans(ctx context.Context, query string, offset, limit int) ([]*models.TravelPlan, int, error)

	// UpdatePlan updates an existing travel plan
	// Validates updated fields
	// Does NOT perform authorization - caller must check CanBeEditedBy
	// Returns error if plan doesn't exist
	UpdatePlan(ctx context.Context, plan *models.TravelPlan) error

	// PublishPlan transitions a draft plan to published status
	// Only callable on draft plans
	// Returns error if plan not found or not draft
	PublishPlan(ctx context.Context, planID string) error

	// SuspendPlan transitions any plan to suspended status (admin action)
	// Admin can suspend any plan
	// Returns error if plan not found
	SuspendPlan(ctx context.Context, planID string) error

	// DeletePlan performs soft-delete (sets is_deleted_by_admin = true)
	// Plan hidden from listings but preserved for history
	// Returns error if plan not found
	DeletePlan(ctx context.Context, planID string) error

	// Linked List Node Operations

	// AddNodeToPlan adds a node to a travel plan at a specific sequence position
	// If position <= 0: appends to end
	// If position > 0: inserts at position, atomically shifts subsequent nodes
	// Validates:
	//   - Plan exists
	//   - Node exists
	//   - Node not already in plan
	//   - No consecutive attractions (transitions allowed consecutively)
	// Returns new PlanNode ID on success
	// Returns ErrValidation if validation fails
	AddNodeToPlan(ctx context.Context, planID, nodeID string, position int) (string, error)

	// GetPlanNodes retrieves all nodes in a travel plan with their details
	// Ordered by sequence_position (1..N)
	// Includes both Node and NodeDetail information for display
	// Returns slice of plans nodes (may be empty if no nodes), error
	GetPlanNodes(ctx context.Context, planID string) ([]*models.PlanNode, error)

	// GetPlanNodeDetails retrieves nodes with their full details (name, category, etc.)
	// Used for display purposes
	// Returns map of PlanNode -> full node details (attraction or transition)
	GetPlanNodeDetails(ctx context.Context, planID string) (map[string]interface{}, error)

	// ReorderNodeInPlan reorders a node to a new sequence position within a plan
	// Atomically updates all affected position values
	// Validates new position doesn't violate linked list constraints
	// Returns error if node not in plan or position invalid
	ReorderNodeInPlan(ctx context.Context, planID, nodeID string, newPosition int) error

	// RemoveNodeFromPlan removes a node from a travel plan
	// Atomically decrements positions of all subsequent nodes
	// Returns error if node not in plan
	RemoveNodeFromPlan(ctx context.Context, planID, nodeID string) error

	// ValidatePlanNodeSequence validates that a plan's node sequence has no gaps and follows rules
	// Checks:
	//   - All positions from 1..N are present
	//   - No consecutive attraction nodes
	// Used for data integrity validation
	// Returns error if sequence invalid
	ValidatePlanNodeSequence(ctx context.Context, planID string) error

	// Planning Utilities

	// GetPlanStatistics retrieves aggregate statistics about a plan
	// Includes: node count, average rating, comment count
	GetPlanStatistics(ctx context.Context, planID string) (map[string]interface{}, error)

	// CanUserEditPlan checks if a user can edit a specific plan
	// Returns true if user is author and plan is draft, or user is admin
	CanUserEditPlan(ctx context.Context, userID string, userRole models.UserRole, planID string) (bool, error)

	// CanUserViewPlan checks if a user can view a specific plan
	// Returns true if plan is published, or user is author/admin
	CanUserViewPlan(ctx context.Context, userID string, userRole models.UserRole, planID string) (bool, error)

	// CountPublishedPlans returns total count of published plans (not admin-deleted)
	CountPublishedPlans(ctx context.Context) (int64, error)

	// GetAverageRating returns the average rating for a plan
	// Returns 0.0 if no ratings exist
	GetAverageRating(ctx context.Context, planID string) (float64, error)

	// CountPlanNodes returns the total number of nodes in a plan
	CountPlanNodes(ctx context.Context, planID string) (int, error)
}

// RelationalPlanService implements PlanService with relational database backend
// Handles all business logic for travel plan management, including linked-list operations
// and denormalization consistency
type RelationalPlanService struct {
	planRepo repositories.PlanRepository
	nodeRepo repositories.NodeRepository
}

// NewRelationalPlanService creates a new plan service
func NewRelationalPlanService(planRepo repositories.PlanRepository, nodeRepo repositories.NodeRepository) PlanService {
	return &RelationalPlanService{
		planRepo: planRepo,
		nodeRepo: nodeRepo,
	}
}

// ============================================================================
// TravelPlan CRUD Operations
// ============================================================================

// CreatePlan creates a new travel plan with validation
func (s *RelationalPlanService) CreatePlan(ctx context.Context, plan *models.TravelPlan) (string, error) {
	// Validate plan
	if err := plan.Validate(); err != nil {
		return "", err
	}

	// Generate ID if not provided
	if plan.ID == "" {
		plan.ID = uuid.New().String()
	}

	// Set default status if not provided
	if plan.Status == "" {
		plan.Status = models.TravelPlanStatusDraft.String()
	}

	// Set timestamps
	now := time.Now().UTC()
	plan.CreatedAt = now
	plan.UpdatedAt = now

	// Persist to database
	return s.planRepo.CreatePlan(ctx, plan)
}

// GetPlanByID retrieves a travel plan by its ID
// Returns nil plan if not found
func (s *RelationalPlanService) GetPlanByID(ctx context.Context, planID string) (*models.TravelPlan, error) {
	return s.planRepo.GetPlanByID(ctx, planID)
}

// GetPlansByAuthor retrieves all travel plans authored by a specific user
// Includes draft, published, and suspended plans
func (s *RelationalPlanService) GetPlansByAuthor(ctx context.Context, authorID string, offset, limit int) ([]*models.TravelPlan, int, error) {
	plans, err := s.planRepo.GetPlansByAuthor(ctx, authorID, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	// Get total count for pagination
	count, err := s.planRepo.CountPlans(ctx)
	if err != nil {
		return nil, 0, err
	}

	return plans, int(count), nil
}

// ListPublishedPlans retrieves all published travel plans (excluding admin-deleted)
func (s *RelationalPlanService) ListPublishedPlans(ctx context.Context, offset, limit int) ([]*models.TravelPlan, int, error) {
	plans, err := s.planRepo.ListPublishedPlans(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	// Get total count for pagination
	count, err := s.planRepo.CountPlans(ctx)
	if err != nil {
		return nil, 0, err
	}

	return plans, int(count), nil
}

// SearchPlans searches travel plans by destination or title (published only)
func (s *RelationalPlanService) SearchPlans(ctx context.Context, query string, offset, limit int) ([]*models.TravelPlan, int, error) {
	plans, err := s.planRepo.SearchPlans(ctx, query, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	// Get total count for pagination
	count, err := s.planRepo.CountPlans(ctx)
	if err != nil {
		return nil, 0, err
	}

	return plans, int(count), nil
}

// UpdatePlan updates an existing travel plan
// Validates updated fields
func (s *RelationalPlanService) UpdatePlan(ctx context.Context, plan *models.TravelPlan) error {
	// Validate plan
	if err := plan.Validate(); err != nil {
		return err
	}

	// Update timestamp
	now := time.Now().UTC()
	plan.UpdatedAt = now

	return s.planRepo.UpdatePlan(ctx, plan)
}

// PublishPlan transitions a draft plan to published status
func (s *RelationalPlanService) PublishPlan(ctx context.Context, planID string) error {
	plan, err := s.planRepo.GetPlanByID(ctx, planID)
	if err != nil {
		return err
	}

	if plan == nil {
		return models.ErrNotFound
	}

	// Can only publish draft plans
	if plan.Status != models.TravelPlanStatusDraft.String() {
		return models.ErrValidation
	}

	// Update status
	plan.Status = models.TravelPlanStatusPublished.String()
	plan.UpdatedAt = time.Now().UTC()

	return s.planRepo.UpdatePlan(ctx, plan)
}

// SuspendPlan transitions any plan to suspended status (admin action)
func (s *RelationalPlanService) SuspendPlan(ctx context.Context, planID string) error {
	plan, err := s.planRepo.GetPlanByID(ctx, planID)
	if err != nil {
		return err
	}

	if plan == nil {
		return models.ErrNotFound
	}

	// Can suspend any plan
	plan.Status = models.TravelPlanStatusSuspended.String()
	plan.UpdatedAt = time.Now().UTC()

	return s.planRepo.UpdatePlan(ctx, plan)
}

// DeletePlan performs soft-delete (sets is_deleted_by_admin = true)
func (s *RelationalPlanService) DeletePlan(ctx context.Context, planID string) error {
	plan, err := s.planRepo.GetPlanByID(ctx, planID)
	if err != nil {
		return err
	}

	if plan == nil {
		return models.ErrNotFound
	}

	// Just set the flag; repository handles the update
	return s.planRepo.DeletePlan(ctx, planID)
}

// ============================================================================
// Linked List Node Operations
// ============================================================================

// AddNodeToPlan adds a node to a travel plan at a specific sequence position
// Validates no consecutive attractions before adding
func (s *RelationalPlanService) AddNodeToPlan(ctx context.Context, planID, nodeID string, position int) (string, error) {
	// Validate plan exists
	plan, err := s.planRepo.GetPlanByID(ctx, planID)
	if err != nil {
		return "", err
	}

	if plan == nil {
		return "", models.ErrNotFound
	}

	// Validate node exists
	detail, err := s.nodeRepo.GetNodeByID(ctx, nodeID)
	if err != nil {
		return "", err
	}

	if detail == nil {
		return "", models.ErrNotFound
	}

	// Get the base node from detail for type checking
	var nodeType string
	switch d := detail.(type) {
	case *models.AttractionNodeDetail:
		if d.Node != nil {
			nodeType = d.Node.Type
		}
	case *models.TransitionNodeDetail:
		if d.Node != nil {
			nodeType = d.Node.Type
		}
	}

	// Create PlanNode with auto-generated ID
	planNode := &models.PlanNode{
		ID:               uuid.New().String(),
		PlanID:           planID,
		NodeID:           nodeID,
		SequencePosition: position,
		CreatedAt:        time.Now().UTC(),
	}

	// Validate attraction sequence (no consecutive attractions)
	// This is done after getting current nodes
	if nodeType == string(models.NodeTypeAttraction) {
		existingNodes, err := s.planRepo.GetPlanNodes(ctx, planID)
		if err != nil {
			return "", err
		}

		// Check if adding this node would create consecutive attractions
		for _, existingNode := range existingNodes {
			existingDetail, err := s.nodeRepo.GetNodeByID(ctx, existingNode.NodeID)
			if err != nil {
				return "", err
			}

			// Extract node type from detail
			var existingNodeType string
			switch d := existingDetail.(type) {
			case *models.AttractionNodeDetail:
				if d.Node != nil {
					existingNodeType = d.Node.Type
				}
			case *models.TransitionNodeDetail:
				if d.Node != nil {
					existingNodeType = d.Node.Type
				}
			}

			// If the last node in the plan is also an attraction, reject
			if existingDetail != nil && existingNodeType == string(models.NodeTypeAttraction) {
				if position <= 0 || position > len(existingNodes)+1 {
					// Appending to attraction, or inserting at end - would create consecutive
					return "", models.ErrValidation
				}
			}
		}
	}

	// Add node to plan via repository
	return s.planRepo.AddNodeToPlan(ctx, planNode)
}

// GetPlanNodes retrieves all nodes in a travel plan with their details
// Ordered by sequence_position
func (s *RelationalPlanService) GetPlanNodes(ctx context.Context, planID string) ([]*models.PlanNode, error) {
	return s.planRepo.GetPlanNodes(ctx, planID)
}

// GetPlanNodeDetails retrieves nodes with their full details
// Returns map with enriched node information
func (s *RelationalPlanService) GetPlanNodeDetails(ctx context.Context, planID string) (map[string]interface{}, error) {
	planNodes, err := s.planRepo.GetPlanNodes(ctx, planID)
	if err != nil {
		return nil, err
	}

	// Build details map: map[nodeID] -> full node data
	details := make(map[string]interface{})

	for _, pn := range planNodes {
		detail, err := s.nodeRepo.GetNodeByID(ctx, pn.NodeID)
		if err != nil {
			return nil, err
		}

		if detail != nil {
			// Include both planNode and detail
			details[pn.NodeID] = map[string]interface{}{
				"planNode": pn,
				"node":     detail,
			}
		}
	}

	return details, nil
}

// ReorderNodeInPlan reorders a node to a new sequence position
// Atomically updates all affected positions
func (s *RelationalPlanService) ReorderNodeInPlan(ctx context.Context, planID, nodeID string, newPosition int) error {
	return s.planRepo.ReorderNodes(ctx, planID, nodeID, newPosition)
}

// RemoveNodeFromPlan removes a node from a travel plan
// Atomically decrements positions of subsequent nodes
func (s *RelationalPlanService) RemoveNodeFromPlan(ctx context.Context, planID, nodeID string) error {
	return s.planRepo.RemoveNodeFromPlan(ctx, planID, nodeID)
}

// ValidatePlanNodeSequence validates that a plan's node sequence has no gaps
// Checks:
//   - All positions from 1..N are present
//   - No consecutive attraction nodes
func (s *RelationalPlanService) ValidatePlanNodeSequence(ctx context.Context, planID string) error {
	planNodes, err := s.planRepo.GetPlanNodes(ctx, planID)
	if err != nil {
		return err
	}

	if len(planNodes) == 0 {
		return nil // Empty plan is valid
	}

	// Check for contiguous sequence (1..N)
	for i, pn := range planNodes {
		if pn.SequencePosition != i+1 {
			return models.ErrValidation // Gap in sequence
		}
	}

	// Check for consecutive attractions
	for i := 0; i < len(planNodes)-1; i++ {
		currentDetail, err := s.nodeRepo.GetNodeByID(ctx, planNodes[i].NodeID)
		if err != nil {
			return err
		}

		nextDetail, err := s.nodeRepo.GetNodeByID(ctx, planNodes[i+1].NodeID)
		if err != nil {
			return err
		}

		// Extract node types from details
		var currentNodeType string
		var nextNodeType string

		switch d := currentDetail.(type) {
		case *models.AttractionNodeDetail:
			if d.Node != nil {
				currentNodeType = d.Node.Type
			}
		case *models.TransitionNodeDetail:
			if d.Node != nil {
				currentNodeType = d.Node.Type
			}
		}

		switch d := nextDetail.(type) {
		case *models.AttractionNodeDetail:
			if d.Node != nil {
				nextNodeType = d.Node.Type
			}
		case *models.TransitionNodeDetail:
			if d.Node != nil {
				nextNodeType = d.Node.Type
			}
		}

		// Two consecutive attractions is invalid
		if currentDetail != nil && nextDetail != nil &&
			currentNodeType == string(models.NodeTypeAttraction) &&
			nextNodeType == string(models.NodeTypeAttraction) {
			return models.ErrValidation
		}
	}

	return nil
}

// ============================================================================
// Planning Utilities
// ============================================================================

// GetPlanStatistics retrieves aggregate statistics about a plan
func (s *RelationalPlanService) GetPlanStatistics(ctx context.Context, planID string) (map[string]interface{}, error) {
	plan, err := s.planRepo.GetPlanByID(ctx, planID)
	if err != nil {
		return nil, err
	}

	if plan == nil {
		return nil, models.ErrNotFound
	}

	// Get node count
	nodeCount, err := s.planRepo.CountPlanNodes(ctx, planID)
	if err != nil {
		return nil, err
	}

	// Get average rating
	avgRating, err := s.planRepo.GetAverageRating(ctx, planID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"node_count":    nodeCount,
		"avg_rating":    avgRating,
		"comment_count": plan.CommentCount,
		"rating_count":  plan.RatingCount,
	}, nil
}

// CanUserEditPlan checks if a user can edit a specific plan
func (s *RelationalPlanService) CanUserEditPlan(ctx context.Context, userID string, userRole models.UserRole, planID string) (bool, error) {
	plan, err := s.planRepo.GetPlanByID(ctx, planID)
	if err != nil {
		return false, err
	}

	if plan == nil {
		return false, models.ErrNotFound
	}

	return plan.CanBeEditedBy(userID, userRole), nil
}

// CanUserViewPlan checks if a user can view a specific plan
func (s *RelationalPlanService) CanUserViewPlan(ctx context.Context, userID string, userRole models.UserRole, planID string) (bool, error) {
	plan, err := s.planRepo.GetPlanByID(ctx, planID)
	if err != nil {
		return false, err
	}

	if plan == nil {
		return false, models.ErrNotFound
	}

	return plan.CanBeViewedBy(userID, userRole), nil
}

// CountPublishedPlans returns total count of published plans (not admin-deleted)
func (s *RelationalPlanService) CountPublishedPlans(ctx context.Context) (int64, error) {
	return s.planRepo.CountPlans(ctx)
}

// GetAverageRating returns the average rating for a plan
func (s *RelationalPlanService) GetAverageRating(ctx context.Context, planID string) (float64, error) {
	plan, err := s.planRepo.GetPlanByID(ctx, planID)
	if err != nil {
		return 0.0, err
	}

	if plan == nil {
		return 0.0, models.ErrNotFound
	}

	// Calculate average from denormalized fields (RatingSum / RatingCount)
	if plan.RatingCount == 0 {
		return 0.0, nil
	}

	return float64(plan.RatingSum) / float64(plan.RatingCount), nil
}

// CountPlanNodes returns the total number of nodes in a plan
func (s *RelationalPlanService) CountPlanNodes(ctx context.Context, planID string) (int, error) {
	nodes, err := s.planRepo.GetPlanNodes(ctx, planID)
	if err != nil {
		return 0, err
	}

	return len(nodes), nil
}
