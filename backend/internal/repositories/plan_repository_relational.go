package repositories

import (
	"context"

	"tll-backend/internal/database"
	"tll-backend/internal/models"

	"gorm.io/gorm"
)

// RelationalPlanRepository implements PlanRepository using GORM for relational database access
// Handles TravelPlan and PlanNode persistence with SQLite3
type RelationalPlanRepository struct {
	*BaseRepository
}

// NewRelationalPlanRepository creates a new relational database plan repository
// Takes a database.Service which provides access to the underlying GORM DB instance
func NewRelationalPlanRepository(dbService database.Service) PlanRepository {
	return &RelationalPlanRepository{
		BaseRepository: NewBaseRepository(dbService),
	}
}

// ============================================================================
// TravelPlan CRUD Operations
// ============================================================================

// CreatePlan creates and persists a new travel plan into the database
// Returns the created plan ID on success
func (r *RelationalPlanRepository) CreatePlan(ctx context.Context, plan *models.TravelPlan) (string, error) {
	if err := r.getDB().WithContext(ctx).Create(plan).Error; err != nil {
		return "", err
	}
	return plan.ID, nil
}

// GetPlanByID retrieves a travel plan by its ID
// Returns nil plan if not found (not an error condition)
func (r *RelationalPlanRepository) GetPlanByID(ctx context.Context, planID string) (*models.TravelPlan, error) {
	var plan models.TravelPlan
	if err := r.FindFirst(ctx, &plan, "id = ?", planID); err != nil {
		if r.isRecordNotFound(err) {
			return nil, nil // Not found is not an error for this interface
		}
		return nil, err
	}
	return &plan, nil
}

// GetPlansByAuthor retrieves all travel plans authored by a specific user
// Excludes deleted plans (only draft, published, and suspended)
// offset: pagination offset (0-based)
// limit: maximum results to return
func (r *RelationalPlanRepository) GetPlansByAuthor(ctx context.Context, authorID string, offset, limit int) ([]*models.TravelPlan, error) {
	var plans []*models.TravelPlan
	query := r.getDB().WithContext(ctx).
		Where("author_id = ? AND status != ?", authorID, models.TravelPlanStatusDeleted.String()).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit)

	if err := query.Find(&plans).Error; err != nil {
		return nil, err
	}

	if plans == nil {
		plans = make([]*models.TravelPlan, 0)
	}

	return plans, nil
}

// ListPublishedPlans retrieves all published travel plans (excludes deleted)
// Paginated for efficient browsing
// offset: pagination offset (0-based)
// limit: maximum results to return
func (r *RelationalPlanRepository) ListPublishedPlans(ctx context.Context, offset, limit int) ([]*models.TravelPlan, error) {
	var plans []*models.TravelPlan
	query := r.getDB().WithContext(ctx).
		Where("status = ? AND status != ?", models.TravelPlanStatusPublished.String(), models.TravelPlanStatusDeleted.String()).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit)

	if err := query.Find(&plans).Error; err != nil {
		return nil, err
	}

	if plans == nil {
		plans = make([]*models.TravelPlan, 0)
	}

	return plans, nil
}

// SearchPlans searches travel plans by destination or title (published only)
// query: search term (matched against destination and title)
// offset: pagination offset (0-based)
// limit: maximum results to return
func (r *RelationalPlanRepository) SearchPlans(ctx context.Context, query string, offset, limit int) ([]*models.TravelPlan, error) {
	var plans []*models.TravelPlan
	searchPattern := "%" + query + "%"

	dbQuery := r.getDB().WithContext(ctx).
		Where("status = ? AND status != ?", models.TravelPlanStatusPublished.String(), models.TravelPlanStatusDeleted.String()).
		Where("destination LIKE ? OR title LIKE ?", searchPattern, searchPattern).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit)

	if err := dbQuery.Find(&plans).Error; err != nil {
		return nil, err
	}

	if plans == nil {
		plans = make([]*models.TravelPlan, 0)
	}

	return plans, nil
}

// UpdatePlan updates an existing travel plan
// Updates: title, description, destination, status, or denormalized counts
// Returns error if plan doesn't exist
func (r *RelationalPlanRepository) UpdatePlan(ctx context.Context, plan *models.TravelPlan) error {
	if err := r.getDB().WithContext(ctx).Model(plan).Updates(plan).Error; err != nil {
		return err
	}
	return nil
}

// DeletePlan performs a soft-delete by setting status to "deleted"
// Plan remains in database but never included in API responses
func (r *RelationalPlanRepository) DeletePlan(ctx context.Context, planID string) error {
	return r.UpdateField(ctx, &models.TravelPlan{}, "status", models.TravelPlanStatusDeleted.String(), "id = ?", planID)
}

// CountPlans returns the total number of published plans (excludes deleted)
func (r *RelationalPlanRepository) CountPlans(ctx context.Context) (int64, error) {
	var count int64
	if err := r.getDB().WithContext(ctx).
		Model(&models.TravelPlan{}).
		Where("status = ? AND status != ?", models.TravelPlanStatusPublished.String(), models.TravelPlanStatusDeleted.String()).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ============================================================================
// PlanNode Operations (Linked List Structure)
// ============================================================================

// AddNodeToPlan adds a node to a travel plan at a specific sequence position
// If sequence_position <= 0: appends node to end of plan (max position + 1)
// If sequence_position > 0: inserts node at that position, atomically updates all positions >= target to shift by +1
// Returns the newly created PlanNode ID on success
// Returns error if plan/node doesn't exist or node already in plan
func (r *RelationalPlanRepository) AddNodeToPlan(ctx context.Context, planNode *models.PlanNode) (string, error) {
	var resultID string
	err := r.getDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Check if node already exists in plan (composite unique constraint violation prevention)
		var existing models.PlanNode
		existsErr := tx.Where("plan_id = ? AND node_id = ?", planNode.PlanID, planNode.NodeID).
			First(&existing).Error
		if existsErr == nil {
			// Node already in plan
			return gorm.ErrDuplicatedKey
		}
		if !r.isRecordNotFound(existsErr) {
			return existsErr
		}

		// If position <= 0, append to end (get max position + 1)
		if planNode.SequencePosition <= 0 {
			var maxPos int
			if err := tx.Model(&models.PlanNode{}).
				Where("plan_id = ?", planNode.PlanID).
				Select("COALESCE(MAX(sequence_position), 0)").
				Scan(&maxPos).Error; err != nil {
				return err
			}
			planNode.SequencePosition = maxPos + 1
		} else {
			// Insert at position: shift all positions >= target by +1
			if err := tx.Model(&models.PlanNode{}).
				Where("plan_id = ? AND sequence_position >= ?", planNode.PlanID, planNode.SequencePosition).
				Update("sequence_position", gorm.Expr("sequence_position + 1")).Error; err != nil {
				return err
			}
		}

		// Create the new plan node
		if err := tx.Create(planNode).Error; err != nil {
			return err
		}

		resultID = planNode.ID
		return nil
	})
	if err != nil {
		return "", err
	}
	return resultID, nil
}

// GetPlanNodes retrieves all nodes in a travel plan, ordered by sequence_position
// Returns slice of PlanNode objects in sequence order
func (r *RelationalPlanRepository) GetPlanNodes(ctx context.Context, planID string) ([]*models.PlanNode, error) {
	var planNodes []*models.PlanNode
	query := r.getDB().WithContext(ctx).
		Where("plan_id = ?", planID).
		Order("sequence_position ASC")

	if err := query.Find(&planNodes).Error; err != nil {
		return nil, err
	}

	if planNodes == nil {
		planNodes = make([]*models.PlanNode, 0)
	}

	return planNodes, nil
}

// GetNodePosition retrieves a specific node's position in a plan
// Returns sequence_position and error status
// Returns error if node not in plan
func (r *RelationalPlanRepository) GetNodePosition(ctx context.Context, planID, nodeID string) (int, error) {
	var position int
	if err := r.getDB().WithContext(ctx).
		Model(&models.PlanNode{}).
		Where("plan_id = ? AND node_id = ?", planID, nodeID).
		Select("sequence_position").
		Scan(&position).Error; err != nil {
		return 0, r.convertNotFoundError(err)
	}
	return position, nil
}

// ReorderNodes updates the sequence position of a node in a plan
// Atomically adjusts affected nodes to maintain contiguous sequence
// Returns error if node not in plan or position would create gaps
func (r *RelationalPlanRepository) ReorderNodes(ctx context.Context, planID, nodeID string, newPosition int) error {
	return r.getDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get current position of the node
		var currentPos int
		if err := tx.Model(&models.PlanNode{}).
			Where("plan_id = ? AND node_id = ?", planID, nodeID).
			Select("sequence_position").
			Scan(&currentPos).Error; err != nil {
			return r.convertNotFoundError(err)
		}

		if currentPos == newPosition {
			return nil // No change needed
		}

		// If moving down (increase position), shift intermediate nodes up
		if newPosition > currentPos {
			if err := tx.Model(&models.PlanNode{}).
				Where("plan_id = ? AND sequence_position > ? AND sequence_position <= ?", planID, currentPos, newPosition).
				Update("sequence_position", gorm.Expr("sequence_position - 1")).Error; err != nil {
				return err
			}
		} else {
			// If moving up (decrease position), shift intermediate nodes down
			if err := tx.Model(&models.PlanNode{}).
				Where("plan_id = ? AND sequence_position >= ? AND sequence_position < ?", planID, newPosition, currentPos).
				Update("sequence_position", gorm.Expr("sequence_position + 1")).Error; err != nil {
				return err
			}
		}

		// Update the node's position
		if err := tx.Model(&models.PlanNode{}).
			Where("plan_id = ? AND node_id = ?", planID, nodeID).
			Update("sequence_position", newPosition).Error; err != nil {
			return err
		}

		return nil
	})
}

// RemoveNodeFromPlan removes a node from a plan
// Atomically decrements all subsequent positions to maintain contiguous sequence
// Returns error if node not in plan
func (r *RelationalPlanRepository) RemoveNodeFromPlan(ctx context.Context, planID, nodeID string) error {
	return r.getDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get the position of the node to be deleted
		var position int
		if err := tx.Model(&models.PlanNode{}).
			Where("plan_id = ? AND node_id = ?", planID, nodeID).
			Select("sequence_position").
			Scan(&position).Error; err != nil {
			return r.convertNotFoundError(err)
		}

		// Delete the node
		if result := tx.Where("plan_id = ? AND node_id = ?", planID, nodeID).Delete(&models.PlanNode{}); result.Error != nil {
			return result.Error
		}

		// Decrement all positions > removed position
		if err := tx.Model(&models.PlanNode{}).
			Where("plan_id = ? AND sequence_position > ?", planID, position).
			Update("sequence_position", gorm.Expr("sequence_position - 1")).Error; err != nil {
			return err
		}

		return nil
	})
}

// RemoveNodeFromPlanAtPosition removes a node at a specific sequence position
// Atomically decrements all subsequent positions to maintain contiguous sequence
// Returns error if no node at that position
func (r *RelationalPlanRepository) RemoveNodeFromPlanAtPosition(ctx context.Context, planID string, position int) error {
	return r.getDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete the node at the specified position
		result := tx.Where("plan_id = ? AND sequence_position = ?", planID, position).Delete(&models.PlanNode{})
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return models.ErrNotFound
		}

		// Decrement all positions > removed position
		if err := tx.Model(&models.PlanNode{}).
			Where("plan_id = ? AND sequence_position > ?", planID, position).
			Update("sequence_position", gorm.Expr("sequence_position - 1")).Error; err != nil {
			return err
		}

		return nil
	})
}

// CountPlanNodes returns the number of nodes in a plan
func (r *RelationalPlanRepository) CountPlanNodes(ctx context.Context, planID string) (int, error) {
	var count int64
	if err := r.getDB().WithContext(ctx).
		Model(&models.PlanNode{}).
		Where("plan_id = ?", planID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

// ============================================================================
// Denormalized Count Operations
// ============================================================================

// IncrementRatingCount atomically increments the rating_count and rating_sum denormalized fields
// Used transactionally when a new rating is created
func (r *RelationalPlanRepository) IncrementRatingCount(ctx context.Context, planID string, stars int) error {
	if err := r.getDB().WithContext(ctx).
		Model(&models.TravelPlan{}).
		Where("id = ?", planID).
		Updates(map[string]interface{}{
			"rating_count": gorm.Expr("rating_count + 1"),
			"rating_sum":   gorm.Expr("rating_sum + ?", stars),
		}).Error; err != nil {
		return err
	}
	return nil
}

// DecrementRatingCount atomically decrements the rating_count and rating_sum denormalized fields
// Used transactionally when a rating is deleted
func (r *RelationalPlanRepository) DecrementRatingCount(ctx context.Context, planID string, stars int) error {
	if err := r.getDB().WithContext(ctx).
		Model(&models.TravelPlan{}).
		Where("id = ?", planID).
		Updates(map[string]interface{}{
			"rating_count": gorm.Expr("rating_count - 1"),
			"rating_sum":   gorm.Expr("rating_sum - ?", stars),
		}).Error; err != nil {
		return err
	}
	return nil
}

// IncrementCommentCount atomically increments the comment_count denormalized field
// Used transactionally when a new comment is created
func (r *RelationalPlanRepository) IncrementCommentCount(ctx context.Context, planID string) error {
	return r.UpdateField(ctx, &models.TravelPlan{}, "comment_count", gorm.Expr("comment_count + 1"), "id = ?", planID)
}

// DecrementCommentCount atomically decrements the comment_count denormalized field
// Used transactionally when a comment is deleted
func (r *RelationalPlanRepository) DecrementCommentCount(ctx context.Context, planID string) error {
	return r.UpdateField(ctx, &models.TravelPlan{}, "comment_count", gorm.Expr("comment_count - 1"), "id = ?", planID)
}

// GetAverageRating calculates the average rating for a plan
// Returns average (0.0 if no ratings) and error status
func (r *RelationalPlanRepository) GetAverageRating(ctx context.Context, planID string) (float64, error) {
	var plan models.TravelPlan
	if err := r.getDB().WithContext(ctx).
		Model(&models.TravelPlan{}).
		Where("id = ?", planID).
		Select("rating_count", "rating_sum").
		Scan(&plan).Error; err != nil {
		return 0.0, r.convertNotFoundError(err)
	}

	if plan.RatingCount == 0 {
		return 0.0, nil
	}

	average := float64(plan.RatingSum) / float64(plan.RatingCount)
	return average, nil
}

// AddRatingSum atomically adjusts the rating_sum denormalized field
// Used for adjusting rating_sum when ratings are updated
// diff can be positive (add) or negative (subtract)
func (r *RelationalPlanRepository) AddRatingSum(ctx context.Context, planID string, diff int64) error {
	if err := r.getDB().WithContext(ctx).
		Model(&models.TravelPlan{}).
		Where("id = ?", planID).
		Update("rating_sum", gorm.Expr("rating_sum + ?", diff)).Error; err != nil {
		return err
	}
	return nil
}
