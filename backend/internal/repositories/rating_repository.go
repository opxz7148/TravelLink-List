package repositories

import (
	"context"

	"tll-backend/internal/models"
)

// RatingRepository defines the interface for rating data access operations
type RatingRepository interface {
	// SubmitRating creates or updates a rating (upsert operation)
	// If rating already exists for (planID, userID), it updates the stars value
	// Returns the rating ID on success
	SubmitRating(ctx context.Context, rating *models.Rating) (string, error)

	// GetRating retrieves a rating by plan and user
	// Returns nil if not found (not an error)
	GetRating(ctx context.Context, planID string, userID string) (*models.Rating, error)

	// ListRatingsByPlan retrieves all ratings for a specific plan with pagination
	// Ordered by created_at DESC (newest first)
	// offset: 0-based pagination offset
	// limit: maximum results to return
	// Returns ratings, total count, error
	ListRatingsByPlan(ctx context.Context, planID string, offset int, limit int) ([]*models.Rating, int, error)

	// DeleteRating removes a rating
	// Returns error if rating doesn't exist
	DeleteRating(ctx context.Context, planID string, userID string) error

	// GetAverageRating calculates the average rating for a plan
	// Returns 0.0 if no ratings exist
	GetAverageRating(ctx context.Context, planID string) (float64, error)

	// GetRatingCount returns the total count of ratings for a plan
	GetRatingCount(ctx context.Context, planID string) (int, error)

	// GetRatingSum returns the sum of all star ratings for a plan
	// Used with GetRatingCount to calculate average
	GetRatingSum(ctx context.Context, planID string) (int, error)
}
