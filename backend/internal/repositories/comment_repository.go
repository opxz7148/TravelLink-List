package repositories

import (
	"context"

	"tll-backend/internal/models"
)

// CommentRepository defines the interface for comment data access operations
type CommentRepository interface {
	// CreateComment creates and persists a new comment
	// Returns the created comment ID on success
	CreateComment(ctx context.Context, comment *models.Comment) (string, error)

	// GetCommentByID retrieves a comment by its ID
	// Returns nil if not found (not an error)
	GetCommentByID(ctx context.Context, commentID string) (*models.Comment, error)

	// ListCommentsByPlan retrieves all comments for a specific plan with pagination
	// Ordered by created_at DESC (newest first)
	// offset: 0-based pagination offset
	// limit: maximum results to return
	// Returns comments, total count, error
	ListCommentsByPlan(ctx context.Context, planID string, offset int, limit int) ([]*models.Comment, int, error)

	// UpdateComment updates comment text
	// Only allows updating text field
	// Returns error if comment doesn't exist
	UpdateComment(ctx context.Context, commentID string, text string) error

	// DeleteComment soft-deletes a comment (sets is_deleted_by_admin = true)
	// Returns error if comment doesn't exist
	DeleteComment(ctx context.Context, commentID string) error

	// CountCommentsByPlan returns the total count of non-deleted comments for a plan
	CountCommentsByPlan(ctx context.Context, planID string) (int, error)
}
