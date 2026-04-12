package controllers

import (
	"net/http"

	"tll-backend/internal/middleware"
	"tll-backend/internal/services"
	"tll-backend/internal/utilities"

	"github.com/gin-gonic/gin"
)

// RatingController handles rating operations
type RatingController struct {
	planService services.PlanService
	// RatingService will be added when implemented
	// ratingService services.RatingService
}

// NewRatingController creates a new rating controller
func NewRatingController(planService services.PlanService) *RatingController {
	return &RatingController{
		planService: planService,
	}
}

// SubmitRatingRequest represents rating submission request
type SubmitRatingRequest struct {
	Stars int `json:"stars" binding:"required,min=1,max=5"`
}

// RatingResponse represents rating data in API response
type RatingResponse struct {
	ID        string `json:"id"`
	PlanID    string `json:"plan_id"`
	UserID    string `json:"user_id"`
	Stars     int    `json:"stars"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// SubmitRating handles POST /api/v1/plans/:id/ratings - submit or update rating
func (rc *RatingController) SubmitRating(c *gin.Context) {
	planID := c.Param("id")
	var req SubmitRatingRequest

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ValidationErrorResponse(c, "invalid request", nil)
		return
	}

	// Validate stars value
	if req.Stars < 1 || req.Stars > 5 {
		middleware.ValidationErrorResponse(c, "Stars must be between 1 and 5", gin.H{"field": "stars"})
		return
	}

	// Get current user
	_, ok := utilities.GetUserIDFromContext(c)
	if !ok {
		middleware.AuthErrorResponse(c, "User not authenticated")
		return
	}

	// Verify plan exists
	plan, _ := rc.planService.GetPlanByID(c.Request.Context(), planID)
	if plan == nil {
		middleware.NotFoundErrorResponse(c, "Plan not found")
		return
	}

	// TODO: When RatingService is implemented, call it here
	// ratingID, err := rc.ratingService.SubmitRating(c.Request.Context(), planID, userID, req.Stars)

	middleware.SuccessResponse(c, http.StatusCreated, gin.H{
		"message": "Rating submitted successfully (service not yet implemented)",
	})
}

// UpdateRating handles PUT /api/v1/plans/:id/ratings - update user's rating
func (rc *RatingController) UpdateRating(c *gin.Context) {
	planID := c.Param("id")
	var req SubmitRatingRequest

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ValidationErrorResponse(c, "invalid request", nil)
		return
	}

	// Validate stars value
	if req.Stars < 1 || req.Stars > 5 {
		middleware.ValidationErrorResponse(c, "Stars must be between 1 and 5", gin.H{"field": "stars"})
		return
	}

	// Get current user
	userID, ok := utilities.GetUserIDFromContext(c)
	if !ok {
		middleware.AuthErrorResponse(c, "User not authenticated")
		return
	}

	// Verify plan exists
	plan, _ := rc.planService.GetPlanByID(c.Request.Context(), planID)
	if plan == nil {
		middleware.NotFoundErrorResponse(c, "Plan not found")
		return
	}

	// TODO: When RatingService is implemented
	// err := rc.ratingService.UpdateRating(c.Request.Context(), planID, userID, req.Stars)

	_ = userID // Suppress unused warning

	middleware.SuccessResponse(c, http.StatusOK, gin.H{
		"message": "Rating updated successfully (service not yet implemented)",
	})
}

// GetUserRating handles GET /api/v1/plans/:id/my-rating - get current user's rating
func (rc *RatingController) GetUserRating(c *gin.Context) {
	planID := c.Param("id")

	// Get current user
	userID, ok := utilities.GetUserIDFromContext(c)
	if !ok {
		middleware.AuthErrorResponse(c, "User not authenticated")
		return
	}

	// Verify plan exists
	plan, _ := rc.planService.GetPlanByID(c.Request.Context(), planID)
	if plan == nil {
		middleware.NotFoundErrorResponse(c, "Plan not found")
		return
	}

	// TODO: When RatingService is implemented
	// rating, err := rc.ratingService.GetUserRating(c.Request.Context(), planID, userID)

	_ = userID // Suppress unused warning

	middleware.SuccessResponse(c, http.StatusOK, gin.H{
		"rating": nil,
		"message": "Get user rating feature not yet implemented",
	})
}

// GetPlanRatingStats handles GET /api/v1/plans/:id/ratings - get rating statistics
func (rc *RatingController) GetPlanRatingStats(c *gin.Context) {
	planID := c.Param("id")

	// Verify plan exists
	plan, _ := rc.planService.GetPlanByID(c.Request.Context(), planID)
	if plan == nil {
		middleware.NotFoundErrorResponse(c, "Plan not found")
		return
	}

	// Get average rating from service
	avgRating, _ := rc.planService.GetAverageRating(c.Request.Context(), planID)

	middleware.SuccessResponse(c, http.StatusOK, gin.H{
		"plan_id":        planID,
		"rating_average": avgRating,
		"rating_count":   plan.RatingCount,
		"rating_sum":     plan.RatingSum,
	})
}
