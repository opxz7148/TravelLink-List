package controllers

import (
	"net/http"
	"strconv"

	"tll-backend/internal/middleware"
	"tll-backend/internal/models"
	"tll-backend/internal/services"
	"tll-backend/internal/utilities"

	"github.com/gin-gonic/gin"
)

// PromotionController handles promotion request operations
type PromotionController struct {
	promotionService services.PromotionService
	userService      services.UserService
}

// NewPromotionController creates a new promotion controller
func NewPromotionController(promotionService services.PromotionService, userService services.UserService) *PromotionController {
	return &PromotionController{
		promotionService: promotionService,
		userService:      userService,
	}
}

// SubmitPromotionRequest represents a promotion request submission
type SubmitPromotionRequest struct {
	PlanID *string `json:"plan_id"` // Optional: if provided, request is for plan promotion; if null, request is for user role upgrade
}

// PromotionRequestResponse represents promotion request in API response
type PromotionRequestResponse struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	PlanID     *string `json:"plan_id"`
	Status     string `json:"status"`
	AdminNotes string `json:"admin_notes"`
	CreatedAt  string `json:"created_at"`
	ReviewedAt *string `json:"reviewed_at"`
}

// SubmitRequest handles POST /api/v1/promotions/request - submit promotion request
// @Summary Submit a promotion request
// @Description Submit a request to upgrade to traveller role or for plan promotion. Simple users can submit requests.
// @Tags promotions
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body SubmitPromotionRequest true "Promotion request submission"
// @Success 201 {object} map[string]interface{} "Request submitted successfully"
// @Failure 400 {object} middleware.SwaggerErrorResponse "Validation error"
// @Failure 401 {object} middleware.SwaggerErrorResponse "Not authenticated"
// @Failure 404 {object} middleware.SwaggerErrorResponse "Plan not found or user not found"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /promotions/request [post]
func (pc *PromotionController) SubmitRequest(c *gin.Context) {
	var req SubmitPromotionRequest

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ValidationErrorResponse(c, "invalid request", nil)
		return
	}

	// Get current user
	userID, ok := utilities.GetUserIDFromContext(c)
	if !ok {
		middleware.AuthErrorResponse(c, "User not authenticated")
		return
	}

	// Submit request
	requestID, err := pc.promotionService.SubmitRequest(c.Request.Context(), userID, req.PlanID)
	if err != nil {
		if err == models.ErrNotFound {
			middleware.NotFoundErrorResponse(c, "Plan not found")
			return
		}
		if err == models.ErrValidation {
			middleware.ValidationErrorResponse(c, "Plan must be authored by you", nil)
			return
		}
		middleware.InternalErrorResponse(c, "Failed to submit promotion request")
		return
	}

	middleware.SuccessResponse(c, http.StatusCreated, gin.H{
		"request_id": requestID,
		"status":     "pending",
	})
}

// ListMyRequests handles GET /api/v1/promotions/my-requests - list user's promotion requests
// @Summary List my promotion requests
// @Description Get all promotion requests submitted by the current user.
// @Tags promotions
// @Security Bearer
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Results per page" default(50)
// @Success 200 {object} map[string]interface{} "List of promotion requests"
// @Failure 401 {object} middleware.SwaggerErrorResponse "Not authenticated"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /promotions/my-requests [get]
func (pc *PromotionController) ListMyRequests(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "50")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}

	offset := (page - 1) * limit

	// Get current user
	userID, ok := utilities.GetUserIDFromContext(c)
	if !ok {
		middleware.AuthErrorResponse(c, "User not authenticated")
		return
	}

	// List requests
	requests, total, err := pc.promotionService.ListUserRequests(c.Request.Context(), userID, offset, limit)
	if err != nil {
		middleware.InternalErrorResponse(c, "Failed to fetch promotion requests")
		return
	}

	// Convert to response format
	responses := make([]PromotionRequestResponse, len(requests))
	for i, req := range requests {
		reviewedAtStr := ""
		if req.ReviewedAt != nil {
			reviewedAtStr = req.ReviewedAt.String()
		}
		responses[i] = PromotionRequestResponse{
			ID:         req.ID,
			UserID:     req.UserID,
			PlanID:     req.PlanID,
			Status:     req.Status,
			AdminNotes: req.AdminNotes,
			CreatedAt:  req.CreatedAt.String(),
			ReviewedAt: &reviewedAtStr,
		}
	}

	middleware.SuccessResponse(c, http.StatusOK, gin.H{
		"requests": responses,
		"pagination": gin.H{
			"current_page": page,
			"limit":        limit,
			"total":        total,
		},
	})
}

// GetRequestStatus handles GET /api/v1/promotions/request/:id - get promotion request status
// @Summary Get promotion request status
// @Description Get the status of a specific promotion request.
// @Tags promotions
// @Security Bearer
// @Produce json
// @Param id path string true "Request ID"
// @Success 200 {object} PromotionRequestResponse "Promotion request details"
// @Failure 401 {object} middleware.SwaggerErrorResponse "Not authenticated"
// @Failure 404 {object} middleware.SwaggerErrorResponse "Request not found"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /promotions/request/{id} [get]
func (pc *PromotionController) GetRequestStatus(c *gin.Context) {
	requestID := c.Param("id")

	// Get current user
	userID, ok := utilities.GetUserIDFromContext(c)
	if !ok {
		middleware.AuthErrorResponse(c, "User not authenticated")
		return
	}

	// Get request
	request, err := pc.promotionService.GetRequest(c.Request.Context(), requestID)
	if err != nil || request == nil {
		middleware.NotFoundErrorResponse(c, "Promotion request not found")
		return
	}

	// Verify ownership (user can only see their own requests)
	if request.UserID != userID {
		middleware.ForbiddenErrorResponse(c, "You can only view your own requests")
		return
	}

	reviewedAtStr := ""
	if request.ReviewedAt != nil {
		reviewedAtStr = request.ReviewedAt.String()
	}

	response := PromotionRequestResponse{
		ID:         request.ID,
		UserID:     request.UserID,
		PlanID:     request.PlanID,
		Status:     request.Status,
		AdminNotes: request.AdminNotes,
		CreatedAt:  request.CreatedAt.String(),
		ReviewedAt: &reviewedAtStr,
	}

	middleware.SuccessResponse(c, http.StatusOK, gin.H{"request": response})
}
