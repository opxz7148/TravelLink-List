package repositories

import (
	"context"

	"tll-backend/internal/models"
)

// NodeRepository defines the interface for node data access operations
// Handles both attraction and transition nodes with polymorphic queries
type NodeRepository interface {
	// CreateAttractionAndSave creates and persists a new attraction node with its details
	// Returns the created node ID on success
	CreateAttractionAndSave(ctx context.Context, node *models.Node, detail *models.AttractionNodeDetail) (string, error)

	// CreateTransitionAndSave creates and persists a new transition node with its details
	// Returns the created node ID on success
	CreateTransitionAndSave(ctx context.Context, node *models.Node, detail *models.TransitionNodeDetail) (string, error)

	// GetNodeByID retrieves a node by its ID with its type-specific details
	// Returns the detail node which contains the embedded base node
	GetNodeByID(ctx context.Context, nodeID string) (interface{}, error)

	// GetAttractionByID retrieves an attraction node with its details
	// Returns nil if node doesn't exist or is not an attraction
	GetAttractionByID(ctx context.Context, nodeID string) (*models.Node, *models.AttractionNodeDetail, error)

	// GetTransitionByID retrieves a transition node with its details
	// Returns nil if node doesn't exist or is not a transition
	GetTransitionByID(ctx context.Context, nodeID string) (*models.Node, *models.TransitionNodeDetail, error)

	// ListApprovedAttractions retrieves all approved attraction nodes with pagination
	// Optional filter: category (empty string = all categories)
	// Pagination: offset and limit
	ListApprovedAttractions(ctx context.Context, category string, offset int, limit int) ([]*models.Node, []*models.AttractionNodeDetail, error)

	// ListApprovedTransitions retrieves all approved transition nodes with pagination
	// Optional filter: mode (empty string = all modes)
	// Pagination: offset and limit
	ListApprovedTransitions(ctx context.Context, mode string, offset int, limit int) ([]*models.Node, []*models.TransitionNodeDetail, error)

	// SearchAttractions searches for attractions by name or category (approved only) with pagination
	// query: search term (matched against name)
	// offset: pagination offset
	// limit: maximum results to return
	SearchAttractions(ctx context.Context, query string, offset int, limit int) ([]*models.Node, []*models.AttractionNodeDetail, error)

	// ListNodesByCreator retrieves all nodes created by a specific user with pagination
	// Includes both approved and unapproved nodes
	// offset and limit for pagination
	ListNodesByCreator(ctx context.Context, creatorID string, offset int, limit int) ([]*models.Node, error)

	// ApproveNode approves a node for public visibility
	// Returns error if node doesn't exist or is already approved
	ApproveNode(ctx context.Context, nodeID string) error

	// DisapproveNode marks a node as unapproved (for moderation)
	// Returns error if node doesn't exist
	DisapproveNode(ctx context.Context, nodeID string) error

	// DeleteNode permanently removes a node
	// Cascades to detail tables automatically via foreign key
	DeleteNode(ctx context.Context, nodeID string) error

	// CountApprovedAttractions returns total count of approved attractions
	CountApprovedAttractions(ctx context.Context) (int64, error)

	// CountApprovedTransitions returns total count of approved transitions
	CountApprovedTransitions(ctx context.Context) (int64, error)

	// CountNodesByCreator returns count of nodes created by a user
	CountNodesByCreator(ctx context.Context, creatorID string) (int64, error)
}
