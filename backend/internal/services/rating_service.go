package services

import (
	"context"
	"fmt"

	"tll-backend/internal/models"
	"tll-backend/internal/repositories"

	"github.com/google/uuid"
)

// RatingService defines the interface for rating-related business logic operations
type RatingService interface {
	// SubmitRating creates or updates a user's rating for a travel plan
	// Validates stars (1-5) and plan existence
	// If user already rated, updates their existing rating
	// Updates TravelPlan.rating_count and rating_sum denormalizations
	// Returns the rating ID on success
	SubmitRating(ctx context.Context, planID string, userID string, stars int) (string, error)

	// GetRating retrieves a user's rating for a specific plan
	// If no rating exists, returns nil (not an error)
	GetRating(ctx context.Context, planID string, userID string) (*models.Rating, error)

	// ListRatingsByPlan retrieves all ratings for a plan with pagination
	// Ordered by created_at DESC (newest first)
	// offset: 0-based pagination offset, limit: max results per page
	ListRatingsByPlan(ctx context.Context, planID string, offset int, limit int) ([]*models.Rating, int, error)

	// DeleteRating removes a user's rating for a plan
	// Only the rater or admin can delete
	// Updates TravelPlan.rating_count and rating_sum denormalizations
	// Returns error if not found or authorization fails
	DeleteRating(ctx context.Context, planID string, userID string, requestUserID string, userRole models.UserRole) error

	// GetAverageRating returns the average rating for a plan
	// Returns 0 if no ratings exist
	GetAverageRating(ctx context.Context, planID string) (float64, error)

	// GetRatingCount returns the count of ratings for a plan
	GetRatingCount(ctx context.Context, planID string) (int, error)

	// GetRatingSum returns the sum of all ratings for a plan
	// Used for denormalization purposes
	GetRatingSum(ctx context.Context, planID string) (int, error)
}

// RelationalRatingService implements RatingService with relational database backend
type RelationalRatingService struct {
	ratingRepo repositories.RatingRepository
	planRepo   repositories.PlanRepository
}

// NewRelationalRatingService creates a new rating service
func NewRelationalRatingService(
	ratingRepo repositories.RatingRepository,
	planRepo repositories.PlanRepository,
) RatingService {
	return &RelationalRatingService{
		ratingRepo: ratingRepo,
		planRepo:   planRepo,
	}
}

// SubmitRating creates or updates a rating (upsert operation)
func (s *RelationalRatingService) SubmitRating(ctx context.Context, planID string, userID string, stars int) (string, error) {
	// Validate inputs
	if planID == "" {
		return "", fmt.Errorf("%w: planID cannot be empty", models.ErrValidation)
	}

	if userID == "" {
		return "", fmt.Errorf("%w: userID cannot be empty", models.ErrValidation)
	}

	if stars < 1 || stars > 5 {
		return "", fmt.Errorf("%w: rating stars must be between 1 and 5", models.ErrValidation)
	}

	// Verify plan exists
	plan, err := s.planRepo.GetPlanByID(ctx, planID)
	if err != nil {
		return "", fmt.Errorf("failed to verify plan: %w", err)
	}

	if plan == nil {
		return "", fmt.Errorf("%w: plan not found", models.ErrNotFound)
	}

	// Check if existing rating
	existingRating, err := s.ratingRepo.GetRating(ctx, planID, userID)
	if err != nil {
		return "", fmt.Errorf("failed to check existing rating: %w", err)
	}

	// Determine if this is a new rating or update
	isNewRating := existingRating == nil
	oldStars := 0
	if existingRating != nil {
		oldStars = existingRating.Stars
	}

	// Create or update rating
	rating := &models.Rating{
		ID:     uuid.New().String(),
		PlanID: planID,
		UserID: userID,
		Stars:  stars,
	}

	if err := rating.Validate(); err != nil {
		return "", err
	}

	// Submit via repository (upsert)
	ratingID, err := s.ratingRepo.SubmitRating(ctx, rating)
	if err != nil {
		return "", err
	}

	// Update plan denormalization
	if isNewRating {
		// New rating: increment count, add stars to sum
		if err := s.planRepo.IncrementRatingCount(ctx, planID, stars); err != nil {
			fmt.Printf("warning: failed to increment rating count: %v\n", err)
		}
	} else if oldStars != stars {
		// Existing rating changed: update sum (stars - oldStars)
		diff := int64(stars - oldStars)
		if err := s.planRepo.AddRatingSum(ctx, planID, diff); err != nil {
			fmt.Printf("warning: failed to update rating sum: %v\n", err)
		}
	}
	// If stars unchanged, no denormalization update needed

	return ratingID, nil
}

// GetRating retrieves a user's rating
func (s *RelationalRatingService) GetRating(ctx context.Context, planID string, userID string) (*models.Rating, error) {
	if planID == "" {
		return nil, fmt.Errorf("%w: planID cannot be empty", models.ErrValidation)
	}

	if userID == "" {
		return nil, fmt.Errorf("%w: userID cannot be empty", models.ErrValidation)
	}

	rating, err := s.ratingRepo.GetRating(ctx, planID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get rating: %w", err)
	}

	// Return nil if not found (not an error case)
	return rating, nil
}

// ListRatingsByPlan retrieves ratings for a plan with pagination
func (s *RelationalRatingService) ListRatingsByPlan(ctx context.Context, planID string, offset int, limit int) ([]*models.Rating, int, error) {
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

	ratings, total, err := s.ratingRepo.ListRatingsByPlan(ctx, planID, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list ratings: %w", err)
	}

	return ratings, total, nil
}

// DeleteRating deletes a user's rating with authorization
func (s *RelationalRatingService) DeleteRating(ctx context.Context, planID string, userID string, requestUserID string, userRole models.UserRole) error {
	if planID == "" {
		return fmt.Errorf("%w: planID cannot be empty", models.ErrValidation)
	}

	if userID == "" {
		return fmt.Errorf("%w: userID cannot be empty", models.ErrValidation)
	}

	if requestUserID == "" {
		return fmt.Errorf("%w: requestUserID cannot be empty", models.ErrValidation)
	}

	// Get rating
	rating, err := s.ratingRepo.GetRating(ctx, planID, userID)
	if err != nil {
		return fmt.Errorf("failed to get rating: %w", err)
	}

	if rating == nil {
		return fmt.Errorf("%w: rating not found", models.ErrNotFound)
	}

	// Authorization: only the rater or admin can delete
	isRater := rating.UserID == requestUserID
	isAdmin := userRole == models.RoleAdmin

	if !isRater && !isAdmin {
		return fmt.Errorf("%w: unauthorized to delete rating", models.ErrUnauthorized)
	}

	// Delete rating
	if err := s.ratingRepo.DeleteRating(ctx, planID, userID); err != nil {
		return fmt.Errorf("failed to delete rating: %w", err)
	}

	// Update plan denormalization (decrement count, subtract stars from sum)
	if err := s.planRepo.DecrementRatingCount(ctx, planID, rating.Stars); err != nil {
		fmt.Printf("warning: failed to decrement rating count: %v\n", err)
	}

	return nil
}

// GetAverageRating returns the average rating for a plan
func (s *RelationalRatingService) GetAverageRating(ctx context.Context, planID string) (float64, error) {
	if planID == "" {
		return 0, fmt.Errorf("%w: planID cannot be empty", models.ErrValidation)
	}

	avg, err := s.ratingRepo.GetAverageRating(ctx, planID)
	if err != nil {
		return 0, fmt.Errorf("failed to get average rating: %w", err)
	}

	return avg, nil
}

// GetRatingCount returns the count of ratings
func (s *RelationalRatingService) GetRatingCount(ctx context.Context, planID string) (int, error) {
	if planID == "" {
		return 0, fmt.Errorf("%w: planID cannot be empty", models.ErrValidation)
	}

	count, err := s.ratingRepo.GetRatingCount(ctx, planID)
	if err != nil {
		return 0, fmt.Errorf("failed to get rating count: %w", err)
	}

	return count, nil
}

// GetRatingSum returns the sum of ratings
func (s *RelationalRatingService) GetRatingSum(ctx context.Context, planID string) (int, error) {
	if planID == "" {
		return 0, fmt.Errorf("%w: planID cannot be empty", models.ErrValidation)
	}

	sum, err := s.ratingRepo.GetRatingSum(ctx, planID)
	if err != nil {
		return 0, fmt.Errorf("failed to get rating sum: %w", err)
	}

	return sum, nil
}
