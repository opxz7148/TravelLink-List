package repositories

import (
	"context"
	"time"

	"tll-backend/internal/database"
	"tll-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RelationalPromotionRequestRepository implements PromotionRequestRepository for relational databases
type RelationalPromotionRequestRepository struct {
	db *gorm.DB
}

// NewRelationalPromotionRequestRepository creates a new relational promotion request repository
func NewRelationalPromotionRequestRepository(dbService database.Service) PromotionRequestRepository {
	return &RelationalPromotionRequestRepository{db: dbService.GetDB()}
}

// Create creates a new promotion request
func (r *RelationalPromotionRequestRepository) Create(ctx context.Context, request *models.PromotionRequest) (string, error) {
	if !request.Validate() {
		return "", models.ErrValidation
	}

	// Generate UUID if not set
	if request.ID == "" {
		request.ID = uuid.New().String()
	}

	// Set status to pending if not set
	if request.Status == "" {
		request.Status = models.PromotionStatusPending.String()
	}

	if err := r.db.WithContext(ctx).Create(request).Error; err != nil {
		return "", err
	}

	return request.ID, nil
}

// GetByID retrieves a promotion request by ID
func (r *RelationalPromotionRequestRepository) GetByID(ctx context.Context, id string) (*models.PromotionRequest, error) {
	var request models.PromotionRequest
	if err := r.db.WithContext(ctx).Preload("User").Preload("Plan").First(&request, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &request, nil
}

// GetByUserID retrieves all promotion requests for a user (paginated)
func (r *RelationalPromotionRequestRepository) GetByUserID(ctx context.Context, userID string, offset, limit int) ([]*models.PromotionRequest, int, error) {
	var requests []*models.PromotionRequest
	var count int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&models.PromotionRequest{}).
		Where("user_id = ?", userID).
		Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Plan").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&requests).Error; err != nil {
		return nil, 0, err
	}

	return requests, int(count), nil
}

// GetPending retrieves all pending promotion requests (for admin review)
func (r *RelationalPromotionRequestRepository) GetPending(ctx context.Context, offset, limit int) ([]*models.PromotionRequest, int, error) {
	var requests []*models.PromotionRequest
	var count int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&models.PromotionRequest{}).
		Where("status = ?", models.PromotionStatusPending.String()).
		Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Plan").
		Where("status = ?", models.PromotionStatusPending.String()).
		Order("created_at ASC").
		Offset(offset).
		Limit(limit).
		Find(&requests).Error; err != nil {
		return nil, 0, err
	}

	return requests, int(count), nil
}

// Update updates an existing promotion request
func (r *RelationalPromotionRequestRepository) Update(ctx context.Context, request *models.PromotionRequest) error {
	if !request.Validate() {
		return models.ErrValidation
	}

	result := r.db.WithContext(ctx).Model(request).Updates(request)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return models.ErrNotFound
	}

	return nil
}

// Approve approves a promotion request
func (r *RelationalPromotionRequestRepository) Approve(ctx context.Context, id string, adminNotes string) (*models.PromotionRequest, error) {
	request, err := r.GetByID(ctx, id)
	if err != nil || request == nil {
		return nil, models.ErrNotFound
	}

	if !request.CanTransitionTo(models.PromotionStatusApproved.String()) {
		return nil, models.ErrValidation
	}

	now := time.Now()
	request.Status = models.PromotionStatusApproved.String()
	request.AdminNotes = adminNotes
	request.ReviewedAt = &now

	if err := r.Update(ctx, request); err != nil {
		return nil, err
	}

	return request, nil
}

// Reject rejects a promotion request
func (r *RelationalPromotionRequestRepository) Reject(ctx context.Context, id string, adminNotes string) (*models.PromotionRequest, error) {
	request, err := r.GetByID(ctx, id)
	if err != nil || request == nil {
		return nil, models.ErrNotFound
	}

	if !request.CanTransitionTo(models.PromotionStatusRejected.String()) {
		return nil, models.ErrValidation
	}

	now := time.Now()
	request.Status = models.PromotionStatusRejected.String()
	request.AdminNotes = adminNotes
	request.ReviewedAt = &now

	if err := r.Update(ctx, request); err != nil {
		return nil, err
	}

	return request, nil
}

// Count returns total number of pending promotion requests
func (r *RelationalPromotionRequestRepository) Count(ctx context.Context) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.PromotionRequest{}).
		Where("status = ?", models.PromotionStatusPending.String()).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}
