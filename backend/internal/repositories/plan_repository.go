package repositories

import (
	"context"

	"tll-backend/internal/models"
)

// PlanRepository defines the interface for travel plan data access operations
// Handles CRUD operations on TravelPlan and PlanNode entities
type PlanRepository interface {
	// TravelPlan operations

	// CreatePlan creates and persists a new travel plan
	// Returns the created plan ID on success
	CreatePlan(ctx context.Context, plan *models.TravelPlan) (string, error)

	// GetPlanByID retrieves a travel plan by its ID
	// Returns nil plan if not found (not an error condition)
	GetPlanByID(ctx context.Context, planID string) (*models.TravelPlan, error)

	// GetPlansByAuthor retrieves all travel plans authored by a specific user
	// Includes draft, published, and suspended plans
	// offset: pagination offset (0-based)
	// limit: maximum results to return
	GetPlansByAuthor(ctx context.Context, authorID string, offset, limit int) ([]*models.TravelPlan, error)

	// ListPublishedPlans retrieves all published travel plans (not deleted by admin)
	// Paginated for efficient browsing
	// offset: pagination offset (0-based)
	// limit: maximum results to return
	ListPublishedPlans(ctx context.Context, offset, limit int) ([]*models.TravelPlan, error)

	// SearchPlans searches travel plans by destination or title (published only)
	// query: search term (matched against destination and title)
	// offset: pagination offset (0-based)
	// limit: maximum results to return
	SearchPlans(ctx context.Context, query string, offset, limit int) ([]*models.TravelPlan, error)

	// UpdatePlan updates an existing travel plan
	// Updates: title, description, destination, status, or denormalized counts
	// Returns error if plan doesn't exist
	UpdatePlan(ctx context.Context, plan *models.TravelPlan) error

	// DeletePlan performs a soft-delete by setting is_deleted_by_admin to true
	// Plan remains in database but hidden from listings
	DeletePlan(ctx context.Context, planID string) error

	// CountPlans returns the total number of published plans (not deleted)
	CountPlans(ctx context.Context) (int64, error)

	// PlanNode operations (linked list structure)

	// AddNodeToPlan adds a node to a travel plan at a specific sequence position
	// If sequence_position is provided:
	//   - Inserts node at that position
	//   - Atomically updates all positions >= target to shift by +1
	// If sequence_position is 0:
	//   - Appends node to end of plan (max position + 1)
	// Returns the newly created PlanNode ID on success
	// Returns error if plan/node doesn't exist or node already in plan
	AddNodeToPlan(ctx context.Context, planNode *models.PlanNode) (string, error)

	// GetPlanNodes retrieves all nodes in a travel plan, ordered by sequence_position
	// Returns slice of PlanNode objects in sequence order
	GetPlanNodes(ctx context.Context, planID string) ([]*models.PlanNode, error)

	// GetNodePosition retrieves a specific node's position in a plan
	// Returns sequence_position and error status
	// Returns error if node not in plan
	GetNodePosition(ctx context.Context, planID, nodeID string) (int, error)

	// ReorderNodes updates the sequence position of a node in a plan
	// Atomically adjusts affected nodes to maintain contiguous sequence
	// Returns error if node not in plan or position would create gaps
	ReorderNodes(ctx context.Context, planID, nodeID string, newPosition int) error

	// RemoveNodeFromPlan removes a node from a plan
	// Atomically decrements all subsequent positions to maintain contiguous sequence
	// Returns error if node not in plan
	RemoveNodeFromPlan(ctx context.Context, planID, nodeID string) error

	// RemoveNodeFromPlanAtPosition removes a node at a specific sequence position
	// Atomically decrements all subsequent positions to maintain contiguous sequence
	// Returns error if no node at that position
	RemoveNodeFromPlanAtPosition(ctx context.Context, planID string, position int) error

	// CountPlanNodes returns the number of nodes in a plan
	CountPlanNodes(ctx context.Context, planID string) (int, error)

	// IncrementRatingCount atomically increments the rating_count and rating_sum denormalized fields
	// Used transactionally when a new rating is created
	IncrementRatingCount(ctx context.Context, planID string, stars int) error

	// DecrementRatingCount atomically decrements the rating_count and rating_sum denormalized fields
	// Used transactionally when a rating is deleted
	DecrementRatingCount(ctx context.Context, planID string, stars int) error

	// IncrementCommentCount atomically increments the comment_count denormalized field
	// Used transactionally when a new comment is created
	IncrementCommentCount(ctx context.Context, planID string) error

	// DecrementCommentCount atomically decrements the comment_count denormalized field
	// Used transactionally when a comment is deleted
	DecrementCommentCount(ctx context.Context, planID string) error

	// GetAverageRating calculates the average rating for a plan
	// Returns average (0.0 if no ratings) and error status
	GetAverageRating(ctx context.Context, planID string) (float64, error)

	// AddRatingSum atomically adds to the rating_sum denormalized field
	// Used for adjusting rating_sum when ratings are updated
	// diff can be positive (add) or negative (subtract)
	AddRatingSum(ctx context.Context, planID string, diff int64) error
}
