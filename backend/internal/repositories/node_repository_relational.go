package repositories

import (
	"context"

	"tll-backend/internal/database"
	"tll-backend/internal/models"

	"gorm.io/gorm"
)

// RelationalNodeRepository implements NodeRepository using relational database via GORM
type RelationalNodeRepository struct {
	*BaseRepository
}

// NewRelationalNodeRepository creates a new relational database node repository
func NewRelationalNodeRepository(dbService database.Service) NodeRepository {
	return &RelationalNodeRepository{
		BaseRepository: NewBaseRepository(dbService),
	}
}

// CreateNodeAndSave creates and persists a node with its type-specific details in a transaction
// Accepts either AttractionNodeDetail or TransitionNodeDetail as details parameter
func (r *RelationalNodeRepository) CreateNodeAndSave(ctx context.Context, node *models.Node, detail interface{}) (string, error) {
	// Use transaction to ensure both node and detail are created together
	result := r.getDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create the base node
		if err := tx.Create(node).Error; err != nil {
			return err
		}

		// Set the node ID in the detail based on its type
		switch d := detail.(type) {
		case *models.AttractionNodeDetail:
			d.NodeID = node.ID
		case *models.TransitionNodeDetail:
			d.NodeID = node.ID
		}

		// Create the detail
		if err := tx.Create(detail).Error; err != nil {
			return err
		}

		return nil
	})

	if result != nil {
		return "", result
	}

	return node.ID, nil
}

// CreateAttractionAndSave creates and persists a new attraction node (delegates to CreateNodeAndSave)
func (r *RelationalNodeRepository) CreateAttractionAndSave(ctx context.Context, node *models.Node, detail *models.AttractionNodeDetail) (string, error) {
	return r.CreateNodeAndSave(ctx, node, detail)
}

// CreateTransitionAndSave creates and persists a new transition node (delegates to CreateNodeAndSave)
func (r *RelationalNodeRepository) CreateTransitionAndSave(ctx context.Context, node *models.Node, detail *models.TransitionNodeDetail) (string, error) {
	return r.CreateNodeAndSave(ctx, node, detail)
}

// GetNodeByID retrieves a node by ID with its embedded type-specific details
// Preloads both AttractionNodeDetail and TransitionNodeDetail (only one will be populated)
func (r *RelationalNodeRepository) GetNodeByID(ctx context.Context, nodeID string) (*models.Node, error) {
	var node models.Node

	// Load node with both possible detail relationships
	if err := r.getDB().WithContext(ctx).
		Preload("AttractionNodeDetail").
		Preload("TransitionNodeDetail").
		First(&node, "id = ?", nodeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &node, nil
}

// GetNodeDetailByID retrieves a node with its type-specific details in a generic way
// Uses embedded details via Preload - optimized for single query
func (r *RelationalNodeRepository) GetNodeDetailByID(ctx context.Context, nodeID string, detail models.NodeDetail) (*models.Node, models.NodeDetail, error) {
	var node models.Node
	nodeType := detail.GetNodeType()

	// Get the node with type checking and preload appropriate detail
	if nodeType == models.NodeTypeAttraction {
		if err := r.getDB().WithContext(ctx).
			Preload("AttractionNodeDetail").
			First(&node, "id = ? AND type = ?", nodeID, nodeType).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, nil, nil
			}
			return nil, nil, err
		}
		return &node, node.AttractionNodeDetail, nil
	} else if nodeType == models.NodeTypeTransition {
		if err := r.getDB().WithContext(ctx).
			Preload("TransitionNodeDetail").
			First(&node, "id = ? AND type = ?", nodeID, nodeType).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, nil, nil
			}
			return nil, nil, err
		}
		return &node, node.TransitionNodeDetail, nil
	}

	return nil, nil, nil
}

// GetAttractionByID retrieves an attraction node with its details (delegates to GetNodeDetailByID)
func (r *RelationalNodeRepository) GetAttractionByID(ctx context.Context, nodeID string) (*models.Node, *models.AttractionNodeDetail, error) {
	detail := &models.AttractionNodeDetail{}
	node, detailResult, err := r.GetNodeDetailByID(ctx, nodeID, detail)
	if err != nil || node == nil {
		return nil, nil, err
	}
	if detailResult == nil {
		return node, nil, nil
	}
	return node, detailResult.(*models.AttractionNodeDetail), nil
}

// GetTransitionByID retrieves a transition node with its details (delegates to GetNodeDetailByID)
func (r *RelationalNodeRepository) GetTransitionByID(ctx context.Context, nodeID string) (*models.Node, *models.TransitionNodeDetail, error) {
	detail := &models.TransitionNodeDetail{}
	node, detailResult, err := r.GetNodeDetailByID(ctx, nodeID, detail)
	if err != nil || node == nil {
		return nil, nil, err
	}
	if detailResult == nil {
		return node, nil, nil
	}
	return node, detailResult.(*models.TransitionNodeDetail), nil
}

// ListApprovedAttractions retrieves all approved attractions with optional category filter and pagination
// Preloads embedded AttractionNodeDetail for efficient single query
func (r *RelationalNodeRepository) ListApprovedAttractions(ctx context.Context, category string, offset int, limit int) ([]*models.Node, []*models.AttractionNodeDetail, error) {
	var nodes []*models.Node

	query := r.getDB().WithContext(ctx).
		Preload("AttractionNodeDetail").
		Where("type = ? AND is_approved = ?", models.NodeTypeAttraction, true).
		Order("created_at DESC")

	if offset > 0 {
		query = query.Offset(offset)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&nodes).Error; err != nil {
		return nil, nil, err
	}

	// Filter by category and build details list
	var details []*models.AttractionNodeDetail
	var filteredNodes []*models.Node

	for _, node := range nodes {
		if node.AttractionNodeDetail != nil {
			if category == "" || node.AttractionNodeDetail.Category == category {
				details = append(details, node.AttractionNodeDetail)
				filteredNodes = append(filteredNodes, node)
			}
		}
	}

	return filteredNodes, details, nil
}

// ListApprovedTransitions retrieves all approved transitions with optional mode filter and pagination
// Preloads embedded TransitionNodeDetail for efficient single query
func (r *RelationalNodeRepository) ListApprovedTransitions(ctx context.Context, mode string, offset int, limit int) ([]*models.Node, []*models.TransitionNodeDetail, error) {
	var nodes []*models.Node

	query := r.getDB().WithContext(ctx).
		Preload("TransitionNodeDetail").
		Where("type = ? AND is_approved = ?", models.NodeTypeTransition, true).
		Order("created_at DESC")

	if offset > 0 {
		query = query.Offset(offset)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&nodes).Error; err != nil {
		return nil, nil, err
	}

	// Filter by mode and build details list
	var details []*models.TransitionNodeDetail
	var filteredNodes []*models.Node

	for _, node := range nodes {
		if node.TransitionNodeDetail != nil {
			if mode == "" || node.TransitionNodeDetail.Mode == mode {
				details = append(details, node.TransitionNodeDetail)
				filteredNodes = append(filteredNodes, node)
			}
		}
	}

	return filteredNodes, details, nil
}

// SearchAttractions searches for attractions by name (approved only) with pagination
// Uses embedded AttractionNodeDetail preload for efficient joins
func (r *RelationalNodeRepository) SearchAttractions(ctx context.Context, searchQuery string, offset int, limit int) ([]*models.Node, []*models.AttractionNodeDetail, error) {
	var nodes []*models.Node

	// Search by matching attraction name via embedded detail relationship
	searchPattern := "%" + searchQuery + "%"
	q := r.getDB().WithContext(ctx).
		Joins("JOIN attraction_node_details ON nodes.id = attraction_node_details.node_id").
		Preload("AttractionNodeDetail").
		Where("nodes.type = ? AND nodes.is_approved = ? AND attraction_node_details.name LIKE ?",
			models.NodeTypeAttraction, true, searchPattern).
		Order("nodes.created_at DESC")

	if offset > 0 {
		q = q.Offset(offset)
	}

	if limit > 0 {
		q = q.Limit(limit)
	}

	if err := q.Find(&nodes).Error; err != nil {
		return nil, nil, err
	}

	// Build details list from embedded data
	var details []*models.AttractionNodeDetail
	for _, node := range nodes {
		if node.AttractionNodeDetail != nil {
			details = append(details, node.AttractionNodeDetail)
		}
	}

	return nodes, details, nil
}

// ListNodesByCreator retrieves all nodes created by a user with pagination
// Preloads all detail types for convenient access
func (r *RelationalNodeRepository) ListNodesByCreator(ctx context.Context, creatorID string, offset int, limit int) ([]*models.Node, error) {
	var nodes []*models.Node

	query := r.getDB().WithContext(ctx).
		Preload("AttractionNodeDetail").
		Preload("TransitionNodeDetail").
		Where("created_by = ?", creatorID).
		Order("created_at DESC")

	if offset > 0 {
		query = query.Offset(offset)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&nodes).Error; err != nil {
		return nil, err
	}

	return nodes, nil
}

// ApproveNode approves a node
func (r *RelationalNodeRepository) ApproveNode(ctx context.Context, nodeID string) error {
	result := r.getDB().WithContext(ctx).Model(&models.Node{}).
		Where("id = ?", nodeID).
		Update("is_approved", true)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return models.ErrNotFound
	}

	return nil
}

// DisapproveNode marks a node as unapproved
func (r *RelationalNodeRepository) DisapproveNode(ctx context.Context, nodeID string) error {
	result := r.getDB().WithContext(ctx).Model(&models.Node{}).
		Where("id = ?", nodeID).
		Update("is_approved", false)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return models.ErrNotFound
	}

	return nil
}

// DeleteNode permanently removes a node (cascades to details)
func (r *RelationalNodeRepository) DeleteNode(ctx context.Context, nodeID string) error {
	result := r.getDB().WithContext(ctx).Where("id = ?", nodeID).Delete(&models.Node{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return models.ErrNotFound
	}

	return nil
}

// CountApprovedByType returns count of approved nodes of a specific type
func (r *RelationalNodeRepository) CountApprovedByType(ctx context.Context, nodeType models.NodeType) (int64, error) {
	var count int64
	if err := r.getDB().WithContext(ctx).Model(&models.Node{}).
		Where("type = ? AND is_approved = ?", nodeType, true).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountApprovedAttractions returns count of approved attractions (delegates to CountApprovedByType)
func (r *RelationalNodeRepository) CountApprovedAttractions(ctx context.Context) (int64, error) {
	return r.CountApprovedByType(ctx, models.NodeTypeAttraction)
}

// CountApprovedTransitions returns count of approved transitions (delegates to CountApprovedByType)
func (r *RelationalNodeRepository) CountApprovedTransitions(ctx context.Context) (int64, error) {
	return r.CountApprovedByType(ctx, models.NodeTypeTransition)
}

// CountNodesByCreator returns count of nodes by creator
func (r *RelationalNodeRepository) CountNodesByCreator(ctx context.Context, creatorID string) (int64, error) {
	var count int64
	if err := r.getDB().WithContext(ctx).Model(&models.Node{}).
		Where("created_by = ?", creatorID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Helper function to remove a node from slice by ID
func removeNodeByID(nodes []*models.Node, nodeID string) []*models.Node {
	for i, node := range nodes {
		if node.ID == nodeID {
			return append(nodes[:i], nodes[i+1:]...)
		}
	}
	return nodes
}
