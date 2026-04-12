package repositories

import (
	"context"
	"fmt"

	"tll-backend/internal/database"
	"tll-backend/internal/models"

	"gorm.io/gorm"
)

// RelationalCommentRepository implements CommentRepository using relational database via GORM
type RelationalCommentRepository struct {
	*BaseRepository
}

// NewRelationalCommentRepository creates a new relational database comment repository
func NewRelationalCommentRepository(dbService database.Service) CommentRepository {
	return &RelationalCommentRepository{
		BaseRepository: NewBaseRepository(dbService),
	}
}

// CreateComment creates and persists a new comment
func (r *RelationalCommentRepository) CreateComment(ctx context.Context, comment *models.Comment) (string, error) {
	if err := r.getDB().WithContext(ctx).Create(comment).Error; err != nil {
		return "", fmt.Errorf("failed to create comment: %w", err)
	}

	return comment.ID, nil
}

// GetCommentByID retrieves a comment by ID
func (r *RelationalCommentRepository) GetCommentByID(ctx context.Context, commentID string) (*models.Comment, error) {
	var comment models.Comment

	if err := r.getDB().WithContext(ctx).First(&comment, "id = ?", commentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}

	return &comment, nil
}

// ListCommentsByPlan retrieves all comments for a plan with pagination
func (r *RelationalCommentRepository) ListCommentsByPlan(ctx context.Context, planID string, offset int, limit int) ([]*models.Comment, int, error) {
	var comments []*models.Comment
	var total int64

	query := r.getDB().WithContext(ctx).
		Where("plan_id = ? AND is_deleted_by_admin = ?", planID, false).
		Order("created_at DESC")

	// Get total count
	if err := query.Model(&models.Comment{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count comments: %w", err)
	}

	// Apply pagination
	if offset > 0 {
		query = query.Offset(offset)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&comments).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list comments: %w", err)
	}

	return comments, int(total), nil
}

// UpdateComment updates a comment's text
func (r *RelationalCommentRepository) UpdateComment(ctx context.Context, commentID string, text string) error {
	result := r.getDB().WithContext(ctx).
		Model(&models.Comment{}).
		Where("id = ?", commentID).
		Update("text", text)

	if result.Error != nil {
		return fmt.Errorf("failed to update comment: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("comment not found")
	}

	return nil
}

// DeleteComment soft-deletes a comment
func (r *RelationalCommentRepository) DeleteComment(ctx context.Context, commentID string) error {
	result := r.getDB().WithContext(ctx).
		Model(&models.Comment{}).
		Where("id = ?", commentID).
		Update("is_deleted_by_admin", true)

	if result.Error != nil {
		return fmt.Errorf("failed to delete comment: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("comment not found")
	}

	return nil
}

// CountCommentsByPlan returns the count of non-deleted comments for a plan
func (r *RelationalCommentRepository) CountCommentsByPlan(ctx context.Context, planID string) (int, error) {
	var count int64

	if err := r.getDB().WithContext(ctx).
		Model(&models.Comment{}).
		Where("plan_id = ? AND is_deleted_by_admin = ?", planID, false).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count comments: %w", err)
	}

	return int(count), nil
}
