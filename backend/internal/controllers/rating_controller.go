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
	planService   services.PlanService
	ratingService services.RatingService
}

// NewRatingController creates a new rating controller
func NewRatingController(planService services.PlanService, ratingService services.RatingService) *RatingController {
	return &RatingController{
		planService:   planService,
		ratingService: ratingService,
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
// @Summary Submit a rating for a travel plan
// @Description Rate a travel plan with 1-5 stars (upsert operation). Authenticated users can rate.
// @Tags ratings
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "Plan ID"
// @Param request body SubmitRatingRequest true "Rating submission request (stars: 1-5)"
// @Success 201 {object} map[string]interface{} "Rating submitted with ID"
// @Failure 400 {object} middleware.SwaggerErrorResponse "Validation error"
// @Failure 401 {object} middleware.SwaggerErrorResponse "Not authenticated"
// @Failure 404 {object} middleware.SwaggerErrorResponse "Plan not found"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /plans/{id}/ratings [post]
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

	// Submit rating via service (upsert)
	ratingID, err := rc.ratingService.SubmitRating(c.Request.Context(), planID, userID, req.Stars)
	if err != nil {
		middleware.InternalErrorResponse(c, "Failed to submit rating")
		return
	}

	middleware.SuccessResponse(c, http.StatusCreated, gin.H{
		"id":      ratingID,
		"message": "Rating submitted successfully",
	})
}

// UpdateRating handles PUT /api/v1/plans/:id/ratings - update user's rating
// @Summary Update a travel plan rating
// @Description Update your rating for a travel plan (1-5 stars). Authenticated users can update their rating.
// @Tags ratings
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "Plan ID"
// @Param request body SubmitRatingRequest true "Rating update request (stars: 1-5)"
// @Success 200 {object} map[string]interface{} "Rating updated with ID"
// @Failure 400 {object} middleware.SwaggerErrorResponse "Validation error"
// @Failure 401 {object} middleware.SwaggerErrorResponse "Not authenticated"
// @Failure 404 {object} middleware.SwaggerErrorResponse "Plan not found"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /plans/{id}/ratings [put]
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

	// Update rating via service (SubmitRating is an upsert operation)
	ratingID, err := rc.ratingService.SubmitRating(c.Request.Context(), planID, userID, req.Stars)
	if err != nil {
		middleware.InternalErrorResponse(c, "Failed to update rating")
		return
	}

	middleware.SuccessResponse(c, http.StatusOK, gin.H{
		"id":      ratingID,
		"message": "Rating updated successfully",
	})
}

// GetUserRating handles GET /api/v1/plans/:id/my-rating - get current user's rating
// @Summary Get your rating for a travel plan
// @Description Retrieve your current rating for a travel plan. Returns null if not rated.
// @Tags ratings
// @Security Bearer
// @Produce json
// @Param id path string true "Plan ID"
// @Success 200 {object} map[string]interface{} "Your rating or null"
// @Failure 401 {object} middleware.SwaggerErrorResponse "Not authenticated"
// @Failure 404 {object} middleware.SwaggerErrorResponse "Plan not found"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /plans/{id}/my-rating [get]
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

	// Get user's rating via service
	rating, err := rc.ratingService.GetRating(c.Request.Context(), planID, userID)
	if err != nil {
		middleware.InternalErrorResponse(c, "Failed to get rating")
		return
	}

	// If no rating exists, return nil response
	if rating == nil {
		middleware.SuccessResponse(c, http.StatusOK, gin.H{
			"rating": nil,
		})
		return
	}

	// Convert to response format
	ratingResponse := RatingResponse{
		ID:        rating.ID,
		PlanID:    rating.PlanID,
		UserID:    rating.UserID,
		Stars:     rating.Stars,
		CreatedAt: rating.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: rating.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	middleware.SuccessResponse(c, http.StatusOK, gin.H{
		"rating": ratingResponse,
	})
}

// GetPlanRatingStats handles GET /api/v1/plans/:id/ratings - get rating statistics
// @Summary Get rating statistics for a travel plan
// @Description Retrieve aggregate rating statistics for a travel plan (public endpoint)
// @Tags ratings
// @Produce json
// @Param id path string true "Plan ID"
// @Success 200 {object} map[string]interface{} "Rating statistics (average, count, sum)"
// @Failure 404 {object} middleware.SwaggerErrorResponse "Plan not found"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /plans/{id}/ratings [get]
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
