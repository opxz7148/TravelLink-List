package repositories

import (
	"context"

	"tll-backend/internal/models"
)

// PromotionRequestRepository defines methods for promotion request data access
type PromotionRequestRepository interface {
	// Create creates a new promotion request
	Create(ctx context.Context, request *models.PromotionRequest) (string, error)

	// GetByID retrieves a promotion request by ID
	GetByID(ctx context.Context, id string) (*models.PromotionRequest, error)

	// GetByUserID retrieves all promotion requests for a user (paginated)
	GetByUserID(ctx context.Context, userID string, offset, limit int) ([]*models.PromotionRequest, int, error)

	// GetPending retrieves all pending promotion requests (for admin review)
	GetPending(ctx context.Context, offset, limit int) ([]*models.PromotionRequest, int, error)

	// Update updates an existing promotion request
	Update(ctx context.Context, request *models.PromotionRequest) error

	// Approve approves a promotion request and returns updated request
	Approve(ctx context.Context, id string, adminNotes string) (*models.PromotionRequest, error)

	// Reject rejects a promotion request and returns updated request
	Reject(ctx context.Context, id string, adminNotes string) (*models.PromotionRequest, error)

	// Count returns total number of pending promotion requests
	Count(ctx context.Context) (int, error)
}
