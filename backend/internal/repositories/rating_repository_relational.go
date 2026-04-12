package repositories

import (
	"context"
	"fmt"

	"tll-backend/internal/database"
	"tll-backend/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// RelationalRatingRepository implements RatingRepository using relational database via GORM
type RelationalRatingRepository struct {
	*BaseRepository
}

// NewRelationalRatingRepository creates a new relational database rating repository
func NewRelationalRatingRepository(dbService database.Service) RatingRepository {
	return &RelationalRatingRepository{
		BaseRepository: NewBaseRepository(dbService),
	}
}

// SubmitRating creates or updates a rating (upsert)
func (r *RelationalRatingRepository) SubmitRating(ctx context.Context, rating *models.Rating) (string, error) {
	// Use upsert: if (planID, userID) exists, update stars; otherwise insert
	result := r.getDB().WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "plan_id"}, {Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"stars", "updated_at"}),
		}).
		Create(rating)

	if result.Error != nil {
		return "", fmt.Errorf("failed to submit rating: %w", result.Error)
	}

	return rating.ID, nil
}

// GetRating retrieves a rating by plan and user
func (r *RelationalRatingRepository) GetRating(ctx context.Context, planID string, userID string) (*models.Rating, error) {
	var rating models.Rating

	if err := r.getDB().WithContext(ctx).
		First(&rating, "plan_id = ? AND user_id = ?", planID, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get rating: %w", err)
	}

	return &rating, nil
}

// ListRatingsByPlan retrieves all ratings for a plan with pagination
func (r *RelationalRatingRepository) ListRatingsByPlan(ctx context.Context, planID string, offset int, limit int) ([]*models.Rating, int, error) {
	var ratings []*models.Rating
	var total int64

	query := r.getDB().WithContext(ctx).
		Where("plan_id = ?", planID).
		Order("created_at DESC")

	// Get total count
	if err := query.Model(&models.Rating{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count ratings: %w", err)
	}

	// Apply pagination
	if offset > 0 {
		query = query.Offset(offset)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&ratings).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list ratings: %w", err)
	}

	return ratings, int(total), nil
}

// DeleteRating removes a rating
func (r *RelationalRatingRepository) DeleteRating(ctx context.Context, planID string, userID string) error {
	result := r.getDB().WithContext(ctx).
		Where("plan_id = ? AND user_id = ?", planID, userID).
		Delete(&models.Rating{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete rating: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("rating not found")
	}

	return nil
}

// GetAverageRating calculates the average rating for a plan
func (r *RelationalRatingRepository) GetAverageRating(ctx context.Context, planID string) (float64, error) {
	var avg float64

	if err := r.getDB().WithContext(ctx).
		Model(&models.Rating{}).
		Where("plan_id = ?", planID).
		Select("COALESCE(AVG(stars), 0.0)").
		Row().Scan(&avg); err != nil {
		return 0.0, fmt.Errorf("failed to get average rating: %w", err)
	}

	return avg, nil
}

// GetRatingCount returns the count of ratings for a plan
func (r *RelationalRatingRepository) GetRatingCount(ctx context.Context, planID string) (int, error) {
	var count int64

	if err := r.getDB().WithContext(ctx).
		Model(&models.Rating{}).
		Where("plan_id = ?", planID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count ratings: %w", err)
	}

	return int(count), nil
}

// GetRatingSum returns the sum of all star ratings for a plan
func (r *RelationalRatingRepository) GetRatingSum(ctx context.Context, planID string) (int, error) {
	var sum int

	if err := r.getDB().WithContext(ctx).
		Model(&models.Rating{}).
		Where("plan_id = ?", planID).
		Select("COALESCE(SUM(stars), 0)").
		Row().Scan(&sum); err != nil {
		return 0, fmt.Errorf("failed to get rating sum: %w", err)
	}

	return sum, nil
}
