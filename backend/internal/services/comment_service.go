package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"tll-backend/internal/models"
	"tll-backend/internal/repositories"

	"github.com/google/uuid"
)

// CommentService defines the interface for comment-related business logic operations
type CommentService interface {
	// CreateComment creates a new comment on a travel plan
	// Validates comment content and plan existence
	// Returns the created comment ID on success
	CreateComment(ctx context.Context, planID string, authorID string, text string) (string, error)

	// GetCommentByID retrieves a comment by its ID
	// Returns nil if not found
	GetCommentByID(ctx context.Context, commentID string) (*models.Comment, error)

	// ListCommentsByPlan retrieves all non-deleted comments for a travel plan with pagination
	// Ordered by created_at DESC (newest first)
	// offset: 0-based pagination offset, limit: max results per page
	ListCommentsByPlan(ctx context.Context, planID string, offset int, limit int) ([]*models.Comment, int, error)

	// UpdateComment updates a comment's text
	// Only the comment author or admin can update
	// Cannot update comments older than 30 days
	// Returns error if not found or authorization fails
	UpdateComment(ctx context.Context, commentID string, userID string, userRole models.UserRole, text string) error

	// DeleteComment soft-deletes a comment (only author or admin)
	// Updates TravelPlan.comment_count denormalization
	// Returns error if comment not found or authorization fails
	DeleteComment(ctx context.Context, commentID string, userID string, userRole models.UserRole) error

	// GetCommentCount returns the count of non-deleted comments for a plan
	GetCommentCount(ctx context.Context, planID string) (int, error)
}

// RelationalCommentService implements CommentService with relational database backend
type RelationalCommentService struct {
	commentRepo repositories.CommentRepository
	planRepo    repositories.PlanRepository
}

// NewRelationalCommentService creates a new comment service
func NewRelationalCommentService(
	commentRepo repositories.CommentRepository,
	planRepo repositories.PlanRepository,
) CommentService {
	return &RelationalCommentService{
		commentRepo: commentRepo,
		planRepo:    planRepo,
	}
}

// CreateComment creates a new comment with validation
func (s *RelationalCommentService) CreateComment(ctx context.Context, planID string, authorID string, text string) (string, error) {
	// Validate inputs
	if planID == "" {
		return "", fmt.Errorf("%w: planID cannot be empty", models.ErrValidation)
	}

	if authorID == "" {
		return "", fmt.Errorf("%w: authorID cannot be empty", models.ErrValidation)
	}

	// Trim and validate text
	text = strings.TrimSpace(text)
	if text == "" || len(text) < 1 || len(text) > 1000 {
		return "", fmt.Errorf("%w: comment text must be 1-1000 characters", models.ErrValidation)
	}

	// Verify plan exists
	plan, err := s.planRepo.GetPlanByID(ctx, planID)
	if err != nil {
		return "", fmt.Errorf("failed to verify plan: %w", err)
	}

	if plan == nil {
		return "", fmt.Errorf("%w: plan not found", models.ErrNotFound)
	}

	// Create comment
	comment := &models.Comment{
		ID:       uuid.New().String(),
		PlanID:   planID,
		AuthorID: authorID,
		Text:     text,
	}

	if err := comment.Validate(); err != nil {
		return "", err
	}

	// Create via repository
	commentID, err := s.commentRepo.CreateComment(ctx, comment)
	if err != nil {
		return "", err
	}

	// Update plan denormalization (increment comment_count)
	if err := s.planRepo.IncrementCommentCount(ctx, planID); err != nil {
		// Log but don't fail - comment created successfully
		fmt.Printf("warning: failed to update plan comment count: %v\n", err)
	}

	return commentID, nil
}

// GetCommentByID retrieves a comment
func (s *RelationalCommentService) GetCommentByID(ctx context.Context, commentID string) (*models.Comment, error) {
	if commentID == "" {
		return nil, fmt.Errorf("%w: commentID cannot be empty", models.ErrValidation)
	}

	comment, err := s.commentRepo.GetCommentByID(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}

	return comment, nil
}

// ListCommentsByPlan retrieves comments for a plan with pagination
func (s *RelationalCommentService) ListCommentsByPlan(ctx context.Context, planID string, offset int, limit int) ([]*models.Comment, int, error) {
	if planID == "" {
		return nil, 0, fmt.Errorf("%w: planID cannot be empty", models.ErrValidation)
	}

	// Validate pagination
	if offset < 0 {
		offset = 0
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}

	comments, total, err := s.commentRepo.ListCommentsByPlan(ctx, planID, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list comments: %w", err)
	}

	return comments, total, nil
}

// UpdateComment updates a comment's text with authorization
func (s *RelationalCommentService) UpdateComment(ctx context.Context, commentID string, userID string, userRole models.UserRole, text string) error {
	if commentID == "" {
		return fmt.Errorf("%w: commentID cannot be empty", models.ErrValidation)
	}

	if userID == "" {
		return fmt.Errorf("%w: userID cannot be empty", models.ErrValidation)
	}

	// Trim and validate text
	text = strings.TrimSpace(text)
	if text == "" || len(text) < 1 || len(text) > 1000 {
		return fmt.Errorf("%w: comment text must be 1-1000 characters", models.ErrValidation)
	}

	// Get existing comment
	comment, err := s.commentRepo.GetCommentByID(ctx, commentID)
	if err != nil {
		return fmt.Errorf("failed to get comment: %w", err)
	}

	if comment == nil {
		return fmt.Errorf("%w: comment not found", models.ErrNotFound)
	}

	// Authorization: only author or admin can update
	isAuthor := comment.AuthorID == userID
	isAdmin := userRole == models.RoleAdmin

	if !isAuthor && !isAdmin {
		return fmt.Errorf("%w: unauthorized to update comment", models.ErrUnauthorized)
	}

	// Check if comment too old (30+ days) - only enforce for author
	if isAuthor && !isAdmin {
		age := time.Since(comment.CreatedAt)
		if age > 30*24*time.Hour {
			return fmt.Errorf("%w: cannot update comment older than 30 days", models.ErrValidation)
		}
	}

	// Update comment
	if err := s.commentRepo.UpdateComment(ctx, commentID, text); err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}

	return nil
}

// DeleteComment soft-deletes a comment with authorization
func (s *RelationalCommentService) DeleteComment(ctx context.Context, commentID string, userID string, userRole models.UserRole) error {
	if commentID == "" {
		return fmt.Errorf("%w: commentID cannot be empty", models.ErrValidation)
	}

	if userID == "" {
		return fmt.Errorf("%w: userID cannot be empty", models.ErrValidation)
	}

	// Get comment
	comment, err := s.commentRepo.GetCommentByID(ctx, commentID)
	if err != nil {
		return fmt.Errorf("failed to get comment: %w", err)
	}

	if comment == nil {
		return fmt.Errorf("%w: comment not found", models.ErrNotFound)
	}

	// Authorization: only author or admin can delete
	isAuthor := comment.AuthorID == userID
	isAdmin := userRole == models.RoleAdmin

	if !isAuthor && !isAdmin {
		return fmt.Errorf("%w: unauthorized to delete comment", models.ErrUnauthorized)
	}

	// Delete comment
	if err := s.commentRepo.DeleteComment(ctx, commentID); err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	// Update plan denormalization (decrement comment_count)
	if err := s.planRepo.DecrementCommentCount(ctx, comment.PlanID); err != nil {
		// Log but don't fail - comment deleted successfully
		fmt.Printf("warning: failed to update plan comment count: %v\n", err)
	}

	return nil
}

// GetCommentCount returns the count of comments for a plan
func (s *RelationalCommentService) GetCommentCount(ctx context.Context, planID string) (int, error) {
	if planID == "" {
		return 0, fmt.Errorf("%w: planID cannot be empty", models.ErrValidation)
	}

	count, err := s.commentRepo.CountCommentsByPlan(ctx, planID)
	if err != nil {
		return 0, fmt.Errorf("failed to count comments: %w", err)
	}

	return count, nil
}
