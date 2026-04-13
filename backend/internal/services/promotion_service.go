package services

import (
	"context"

	"tll-backend/internal/models"
	"tll-backend/internal/repositories"
)

// PromotionService defines promotion request business logic
type PromotionService interface {
	// SubmitRequest submits a promotion request from a user (for role upgrade or plan promotion)
	SubmitRequest(ctx context.Context, userID string, planID *string) (string, error)

	// GetRequest retrieves a specific promotion request
	GetRequest(ctx context.Context, id string) (*models.PromotionRequest, error)

	// ListUserRequests lists all promotion requests from a user (paginated)
	ListUserRequests(ctx context.Context, userID string, offset, limit int) ([]*models.PromotionRequest, int, error)

	// ListPending lists all pending promotion requests for admin review (paginated)
	ListPending(ctx context.Context, offset, limit int) ([]*models.PromotionRequest, int, error)

	// ApproveRequest approves a promotion request and promotes user role if applicable
	ApproveRequest(ctx context.Context, requestID string, adminNotes string) error

	// RejectRequest rejects a promotion request
	RejectRequest(ctx context.Context, requestID string, adminNotes string) error
}

// RelationalPromotionService implements PromotionService with relational database
type RelationalPromotionService struct {
	promotionRepo repositories.PromotionRequestRepository
	userRepo      repositories.UserRepository
	planRepo      repositories.PlanRepository
}

// NewRelationalPromotionService creates a new relational promotion service
func NewRelationalPromotionService(
	promotionRepo repositories.PromotionRequestRepository,
	userRepo repositories.UserRepository,
	planRepo repositories.PlanRepository,
) PromotionService {
	return &RelationalPromotionService{
		promotionRepo: promotionRepo,
		userRepo:      userRepo,
		planRepo:      planRepo,
	}
}

// SubmitRequest submits a promotion request
func (s *RelationalPromotionService) SubmitRequest(ctx context.Context, userID string, planID *string) (string, error) {
	// Verify user exists
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return "", models.ErrNotFound
	}

	// If planID provided, verify plan exists and belongs to user
	if planID != nil && *planID != "" {
		plan, err := s.planRepo.GetPlanByID(ctx, *planID)
		if err != nil || plan == nil {
			return "", models.ErrNotFound
		}

		if plan.AuthorID != userID {
			return "", models.ErrValidation
		}
	}

	// Create promotion request
	request := &models.PromotionRequest{
		UserID: userID,
		PlanID: planID,
		Status: models.PromotionStatusPending.String(),
	}

	return s.promotionRepo.Create(ctx, request)
}

// GetRequest retrieves a promotion request
func (s *RelationalPromotionService) GetRequest(ctx context.Context, id string) (*models.PromotionRequest, error) {
	return s.promotionRepo.GetByID(ctx, id)
}

// ListUserRequests lists user's promotion requests
func (s *RelationalPromotionService) ListUserRequests(ctx context.Context, userID string, offset, limit int) ([]*models.PromotionRequest, int, error) {
	return s.promotionRepo.GetByUserID(ctx, userID, offset, limit)
}

// ListPending lists pending promotion requests for admin review
func (s *RelationalPromotionService) ListPending(ctx context.Context, offset, limit int) ([]*models.PromotionRequest, int, error) {
	return s.promotionRepo.GetPending(ctx, offset, limit)
}

// ApproveRequest approves a promotion request and promotes user if it's a role upgrade
func (s *RelationalPromotionService) ApproveRequest(ctx context.Context, requestID string, adminNotes string) error {
	// Get the request
	request, err := s.promotionRepo.GetByID(ctx, requestID)
	if err != nil || request == nil {
		return models.ErrNotFound
	}

	// Approve the request
	_, err = s.promotionRepo.Approve(ctx, requestID, adminNotes)
	if err != nil {
		return err
	}

	// If no plan specified, this is a role upgrade request - promote user to traveller
	if request.PlanID == nil || *request.PlanID == "" {
		err = s.userRepo.PromoteToTraveller(ctx, request.UserID)
		if err != nil {
			return err
		}
	}

	return nil
}

// RejectRequest rejects a promotion request
func (s *RelationalPromotionService) RejectRequest(ctx context.Context, requestID string, adminNotes string) error {
	request, err := s.promotionRepo.GetByID(ctx, requestID)
	if err != nil || request == nil {
		return models.ErrNotFound
	}

	_, err = s.promotionRepo.Reject(ctx, requestID, adminNotes)
	return err
}
