package services

import (
	"context"
	"fmt"

	"tll-backend/internal/models"
	"tll-backend/internal/repositories"

	"github.com/google/uuid"
)

// NodeService defines the interface for node-related business logic operations
// Handles both attraction and transition nodes with validations
type NodeService interface {
	// CreateAttractionNode creates a new attraction node with validation
	// Returns the created node ID on success
	CreateAttractionNode(ctx context.Context, createdBy string, detail *models.AttractionNodeDetail) (string, error)

	// CreateTransitionNode creates a new transition node with validation
	// Returns the created node ID on success
	CreateTransitionNode(ctx context.Context, createdBy string, detail *models.TransitionNodeDetail) (string, error)

	// GetNodeByID retrieves a node by its ID with embedded type-specific details
	// Returns a Node with AttractionNodeDetail or TransitionNodeDetail populated based on type
	GetNodeByID(ctx context.Context, nodeID string) (*models.Node, error)

	// ListApprovedNodes retrieves all approved nodes (both attractions and transitions) with pagination
	// Used for public browsing and node selection in plan creation
	// offset: pagination offset, limit: results per page
	ListApprovedNodes(ctx context.Context, offset int, limit int) ([]*models.Node, error)

	// ListApprovedAttractions retrieves all approved attraction nodes with optional category filter and pagination
	// category: optional filter (empty string = all categories)
	// offset: pagination offset, limit: results per page
	ListApprovedAttractions(ctx context.Context, category string, offset int, limit int) ([]*models.Node, []*models.AttractionNodeDetail, error)

	// ListApprovedTransitions retrieves all approved transition nodes with optional mode filter and pagination
	// mode: optional filter (empty string = all modes)
	// offset: pagination offset, limit: results per page
	ListApprovedTransitions(ctx context.Context, mode string, offset int, limit int) ([]*models.Node, []*models.TransitionNodeDetail, error)

	// ListNodesByCreator retrieves all nodes created by a specific user with pagination
	// Includes both approved and unapproved (for user moderation features)
	// offset: pagination offset, limit: results per page
	ListNodesByCreator(ctx context.Context, creatorID string, offset int, limit int) ([]*models.Node, error)

	// ListDraftNodesByCreator retrieves all unapproved (draft) nodes created by a specific user with pagination
	// Used for "My Nodes" tab in node selector
	// offset: pagination offset, limit: results per page
	ListDraftNodesByCreator(ctx context.Context, creatorID string, offset int, limit int) ([]*models.Node, error)

	// ListNodesByType retrieves all approved nodes of a specific type (attraction or transition) with pagination
	// offset: pagination offset, limit: results per page
	ListNodesByType(ctx context.Context, nodeType models.NodeType, offset int, limit int) ([]*models.Node, error)

	// SearchAttractions searches for attractions by name or category (approved only) with pagination
	// Used for finding existing attractions when creating plans
	// offset: pagination offset, limit: results per page
	SearchAttractions(ctx context.Context, query string, offset int, limit int) ([]*models.Node, []*models.AttractionNodeDetail, error)

	// ApproveNode approves a node for public visibility
	// Typically called by admin for moderation of user-created nodes
	ApproveNode(ctx context.Context, nodeID string) error

	// DisapproveNode marks a node as unapproved (for moderation)
	DisapproveNode(ctx context.Context, nodeID string) error

	// DeleteNode permanently removes a node from the system
	DeleteNode(ctx context.Context, nodeID string) error

	// CountApprovedAttractions returns total count of approved attractions
	CountApprovedAttractions(ctx context.Context) (int64, error)

	// CountApprovedTransitions returns total count of approved transitions
	CountApprovedTransitions(ctx context.Context) (int64, error)

	// CountNodesByCreator returns count of nodes created by a user
	CountNodesByCreator(ctx context.Context, creatorID string) (int64, error)

	// CountDraftNodesByCreator returns count of unapproved (draft) nodes created by a user
	CountDraftNodesByCreator(ctx context.Context, creatorID string) (int64, error)
}

// RelationalNodeService implements NodeService with relational database backend
type RelationalNodeService struct {
	nodeRepo repositories.NodeRepository
}

// NewRelationalNodeService creates a new node service
func NewRelationalNodeService(nodeRepo repositories.NodeRepository) NodeService {
	return &RelationalNodeService{
		nodeRepo: nodeRepo,
	}
}

// CreateAttractionNode creates a new attraction node with comprehensive validation
func (s *RelationalNodeService) CreateAttractionNode(ctx context.Context, createdBy string, detail *models.AttractionNodeDetail) (string, error) {
	// Validate input parameters
	if createdBy == "" {
		return "", fmt.Errorf("createdBy cannot be empty")
	}

	if detail == nil {
		return "", fmt.Errorf("attraction detail cannot be nil")
	}

	// Validate attraction detail
	if !detail.Validate() {
		return "", fmt.Errorf("invalid attraction node details: validation failed")
	}

	// Validate category is in allowed list
	validCategories := models.ValidAttractionCategories()
	isValidCategory := false
	for _, cat := range validCategories {
		if cat == detail.Category {
			isValidCategory = true
			break
		}
	}
	if !isValidCategory {
		return "", fmt.Errorf("invalid attraction category: %s", detail.Category)
	}

	// Create node with type discriminator
	node := &models.Node{
		ID:        uuid.New().String(),
		Type:      string(models.NodeTypeAttraction),
		CreatedBy: createdBy,
		// User-created nodes default to is_approved=false, pending admin approval
		IsApproved: false,
	}

	// Create and persist in repository
	nodeID, err := s.nodeRepo.CreateAttractionAndSave(ctx, node, detail)
	if err != nil {
		fmt.Println("err", err)
		return "", fmt.Errorf("failed to create attraction node: %w", err)
	}

	return nodeID, nil
}

// CreateTransitionNode creates a new transition node with comprehensive validation
func (s *RelationalNodeService) CreateTransitionNode(ctx context.Context, createdBy string, detail *models.TransitionNodeDetail) (string, error) {
	// Validate input parameters
	if createdBy == "" {
		return "", fmt.Errorf("createdBy cannot be empty")
	}

	if detail == nil {
		return "", fmt.Errorf("transition detail cannot be nil")
	}

	// Validate transition detail
	if !detail.Validate() {
		return "", fmt.Errorf("invalid transition node details: validation failed")
	}

	// Validate mode is in allowed list
	validModes := models.ValidTransitionModes()
	isValidMode := false
	for _, mode := range validModes {
		if mode == detail.Mode {
			isValidMode = true
			break
		}
	}
	if !isValidMode {
		return "", fmt.Errorf("invalid transition mode: %s", detail.Mode)
	}

	// Create node with type discriminator
	node := &models.Node{
		ID:        uuid.New().String(),
		Type:      string(models.NodeTypeTransition),
		CreatedBy: createdBy,
		// User-created nodes default to is_approved=false, pending admin approval
		IsApproved: false,
	}

	// Create and persist in repository
	nodeID, err := s.nodeRepo.CreateTransitionAndSave(ctx, node, detail)
	if err != nil {
		return "", fmt.Errorf("failed to create transition node: %w", err)
	}

	return nodeID, nil
}

// GetNodeByID retrieves a node and its embedded type-specific details
// Returns a Node with AttractionNodeDetail or TransitionNodeDetail populated
func (s *RelationalNodeService) GetNodeByID(ctx context.Context, nodeID string) (*models.Node, error) {
	if nodeID == "" {
		return nil, fmt.Errorf("nodeID cannot be empty")
	}

	node, err := s.nodeRepo.GetNodeByID(ctx, nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get node: %w", err)
	}

	return node, nil
}

// ListApprovedNodes retrieves all approved nodes (both types) with pagination
func (s *RelationalNodeService) ListApprovedNodes(ctx context.Context, offset int, limit int) ([]*models.Node, error) {
	// Fetch both attractions and transitions
	attractionNodes, _, err := s.nodeRepo.ListApprovedAttractions(ctx, "", offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list approved attractions: %w", err)
	}

	transitionNodes, _, err := s.nodeRepo.ListApprovedTransitions(ctx, "", offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list approved transitions: %w", err)
	}

	// Combine and return all approved nodes
	allNodes := append(attractionNodes, transitionNodes...)
	return allNodes, nil
}

// ListApprovedAttractions retrieves approved attraction nodes with optional category filter and pagination
func (s *RelationalNodeService) ListApprovedAttractions(ctx context.Context, category string, offset int, limit int) ([]*models.Node, []*models.AttractionNodeDetail, error) {
	nodes, details, err := s.nodeRepo.ListApprovedAttractions(ctx, category, offset, limit)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list approved attractions: %w", err)
	}

	return nodes, details, nil
}

// ListApprovedTransitions retrieves approved transition nodes with optional mode filter and pagination
func (s *RelationalNodeService) ListApprovedTransitions(ctx context.Context, mode string, offset int, limit int) ([]*models.Node, []*models.TransitionNodeDetail, error) {
	nodes, details, err := s.nodeRepo.ListApprovedTransitions(ctx, mode, offset, limit)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list approved transitions: %w", err)
	}

	return nodes, details, nil
}

// ListNodesByCreator retrieves all nodes created by a specific user with pagination
func (s *RelationalNodeService) ListNodesByCreator(ctx context.Context, creatorID string, offset int, limit int) ([]*models.Node, error) {
	if creatorID == "" {
		return nil, fmt.Errorf("creatorID cannot be empty")
	}

	nodes, err := s.nodeRepo.ListNodesByCreator(ctx, creatorID, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes by creator: %w", err)
	}

	return nodes, nil
}

// ListDraftNodesByCreator retrieves all unapproved (draft) nodes created by a user with pagination
func (s *RelationalNodeService) ListDraftNodesByCreator(ctx context.Context, creatorID string, offset int, limit int) ([]*models.Node, error) {
	if creatorID == "" {
		return nil, fmt.Errorf("creatorID cannot be empty")
	}

	nodes, err := s.nodeRepo.ListDraftNodesByCreator(ctx, creatorID, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list draft nodes by creator: %w", err)
	}

	return nodes, nil
}

// ListNodesByType retrieves all approved nodes of a specific type with pagination
func (s *RelationalNodeService) ListNodesByType(ctx context.Context, nodeType models.NodeType, offset int, limit int) ([]*models.Node, error) {
	// Validate node type
	if !models.CheckNodeType(nodeType) {
		return nil, fmt.Errorf("invalid node type: %s", nodeType)
	}

	// Fetch nodes based on type
	var nodes []*models.Node
	var err error

	if nodeType == models.NodeTypeAttraction {
		nodes, _, err = s.nodeRepo.ListApprovedAttractions(ctx, "", offset, limit)
	} else if nodeType == models.NodeTypeTransition {
		nodes, _, err = s.nodeRepo.ListApprovedTransitions(ctx, "", offset, limit)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list nodes by type: %w", err)
	}

	return nodes, nil
}

// SearchAttractions searches for attractions by name (approved only) with pagination
func (s *RelationalNodeService) SearchAttractions(ctx context.Context, query string, offset int, limit int) ([]*models.Node, []*models.AttractionNodeDetail, error) {
	if query == "" {
		return nil, nil, fmt.Errorf("search query cannot be empty")
	}

	nodes, details, err := s.nodeRepo.SearchAttractions(ctx, query, offset, limit)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to search attractions: %w", err)
	}

	return nodes, details, nil
}

// ApproveNode approves a node for public visibility
func (s *RelationalNodeService) ApproveNode(ctx context.Context, nodeID string) error {
	if nodeID == "" {
		return fmt.Errorf("nodeID cannot be empty")
	}

	err := s.nodeRepo.ApproveNode(ctx, nodeID)
	if err != nil {
		return fmt.Errorf("failed to approve node: %w", err)
	}

	return nil
}

// DisapproveNode marks a node as unapproved
func (s *RelationalNodeService) DisapproveNode(ctx context.Context, nodeID string) error {
	if nodeID == "" {
		return fmt.Errorf("nodeID cannot be empty")
	}

	err := s.nodeRepo.DisapproveNode(ctx, nodeID)
	if err != nil {
		return fmt.Errorf("failed to disapprove node: %w", err)
	}

	return nil
}

// DeleteNode permanently removes a node from the system
func (s *RelationalNodeService) DeleteNode(ctx context.Context, nodeID string) error {
	if nodeID == "" {
		return fmt.Errorf("nodeID cannot be empty")
	}

	err := s.nodeRepo.DeleteNode(ctx, nodeID)
	if err != nil {
		return fmt.Errorf("failed to delete node: %w", err)
	}

	return nil
}

// CountApprovedAttractions returns total count of approved attractions
func (s *RelationalNodeService) CountApprovedAttractions(ctx context.Context) (int64, error) {
	count, err := s.nodeRepo.CountApprovedAttractions(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count approved attractions: %w", err)
	}

	return count, nil
}

// CountApprovedTransitions returns total count of approved transitions
func (s *RelationalNodeService) CountApprovedTransitions(ctx context.Context) (int64, error) {
	count, err := s.nodeRepo.CountApprovedTransitions(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count approved transitions: %w", err)
	}

	return count, nil
}

// CountNodesByCreator returns count of nodes created by a user
func (s *RelationalNodeService) CountNodesByCreator(ctx context.Context, creatorID string) (int64, error) {
	if creatorID == "" {
		return 0, fmt.Errorf("creatorID cannot be empty")
	}

	count, err := s.nodeRepo.CountNodesByCreator(ctx, creatorID)
	if err != nil {
		return 0, fmt.Errorf("failed to count nodes by creator: %w", err)
	}

	return count, nil
}

// CountDraftNodesByCreator returns count of unapproved (draft) nodes created by a user
func (s *RelationalNodeService) CountDraftNodesByCreator(ctx context.Context, creatorID string) (int64, error) {
	if creatorID == "" {
		return 0, fmt.Errorf("creatorID cannot be empty")
	}

	count, err := s.nodeRepo.CountDraftNodesByCreator(ctx, creatorID)
	if err != nil {
		return 0, fmt.Errorf("failed to count draft nodes by creator: %w", err)
	}

	return count, nil
}
