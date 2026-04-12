package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"tll-backend/internal/middleware"
	"tll-backend/internal/models"
	"tll-backend/internal/services"
	"tll-backend/internal/utilities"

	"github.com/gin-gonic/gin"
)

// PlanController handles travel plan operations
type PlanController struct {
	planService services.PlanService
}

// NewPlanController creates a new plan controller
func NewPlanController(planService services.PlanService) *PlanController {
	return &PlanController{
		planService: planService,
	}
}

// PlanResponse represents plan data in API response
type PlanResponse struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Destination    string `json:"destination"`
	AuthorID       string `json:"author_id"`
	Status         string `json:"status"`
	RatingAverage  float64 `json:"rating_average"`
	RatingCount    int    `json:"rating_count"`
	CommentCount   int    `json:"comment_count"`
	NodeCount      int    `json:"node_count"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// CreatePlanRequest represents plan creation request
type CreatePlanRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=150"`
	Description string `json:"description" binding:"max=1000"`
	Destination string `json:"destination" binding:"required,min=1,max=200"`
}

// UpdatePlanRequest represents plan update request
type UpdatePlanRequest struct {
	Title       string `json:"title" binding:"min=1,max=150"`
	Description string `json:"description" binding:"max=1000"`
	Destination string `json:"destination" binding:"min=1,max=200"`
}

// BrowsePlans handles GET /api/v1/plans - browse published plans
func (pc *PlanController) BrowsePlans(c *gin.Context) {
	// Parse query parameters
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Fetch plans
	plans, totalCount, err := pc.planService.ListPublishedPlans(c.Request.Context(), offset, limit)
	if err != nil {
		middleware.InternalErrorResponse(c, "Failed to fetch plans")
		return
	}

	// Convert to response format
	planResponses := make([]PlanResponse, len(plans))
	for i, plan := range plans {
		avgRating, _ := pc.planService.GetAverageRating(c.Request.Context(), plan.ID)
		planResponses[i] = PlanResponse{
			ID:           plan.ID,
			Title:        plan.Title,
			Description:  plan.Description,
			Destination:  plan.Destination,
			AuthorID:     plan.AuthorID,
			Status:       plan.Status,
			RatingAverage: avgRating,
			RatingCount:  plan.RatingCount,
			CommentCount: plan.CommentCount,
			CreatedAt:    plan.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:    plan.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}

		// Get node count
		nodeCount, _ := pc.planService.CountPlanNodes(c.Request.Context(), plan.ID)
		planResponses[i].NodeCount = nodeCount
	}

	resp := gin.H{
		"plans": planResponses,
		"pagination": gin.H{
			"current_page": page,
			"total_pages":  (totalCount + limit - 1) / limit,
			"total_items":  totalCount,
			"limit":        limit,
		},
	}

	middleware.SuccessResponse(c, http.StatusOK, resp)
}

// SearchPlans handles GET /api/v1/plans/search - search published plans
func (pc *PlanController) SearchPlans(c *gin.Context) {
	query := c.Query("q")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	if strings.TrimSpace(query) == "" {
		middleware.ValidationErrorResponse(c, "search query required", gin.H{"field": "q"})
		return
	}

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Fetch search results
	plans, totalCount, err := pc.planService.SearchPlans(c.Request.Context(), query, offset, limit)
	if err != nil {
		middleware.InternalErrorResponse(c, "Failed to search plans")
		return
	}

	// Convert to response format
	planResponses := make([]PlanResponse, len(plans))
	for i, plan := range plans {
		avgRating, _ := pc.planService.GetAverageRating(c.Request.Context(), plan.ID)
		planResponses[i] = PlanResponse{
			ID:           plan.ID,
			Title:        plan.Title,
			Description:  plan.Description,
			Destination:  plan.Destination,
			AuthorID:     plan.AuthorID,
			Status:       plan.Status,
			RatingAverage: avgRating,
			RatingCount:  plan.RatingCount,
			CommentCount: plan.CommentCount,
			CreatedAt:    plan.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:    plan.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}

		// Get node count
		nodeCount, _ := pc.planService.CountPlanNodes(c.Request.Context(), plan.ID)
		planResponses[i].NodeCount = nodeCount
	}

	resp := gin.H{
		"plans": planResponses,
		"pagination": gin.H{
			"current_page": page,
			"total_pages":  (totalCount + limit - 1) / limit,
			"total_items":  totalCount,
			"limit":        limit,
		},
	}

	middleware.SuccessResponse(c, http.StatusOK, resp)
}

// GetPlanDetails handles GET /api/v1/plans/:id - get plan details
func (pc *PlanController) GetPlanDetails(c *gin.Context) {
	planID := c.Param("id")

	// Fetch plan
	plan, err := pc.planService.GetPlanByID(c.Request.Context(), planID)
	if err != nil {
		middleware.InternalErrorResponse(c, "Failed to fetch plan")
		return
	}

	if plan == nil {
		middleware.NotFoundErrorResponse(c, "Plan not found")
		return
	}

	// Get additional stats
	avgRating, _ := pc.planService.GetAverageRating(c.Request.Context(), plan.ID)
	nodeCount, _ := pc.planService.CountPlanNodes(c.Request.Context(), plan.ID)

	// Fetch nodes
	nodes, _ := pc.planService.GetPlanNodes(c.Request.Context(), plan.ID)

	resp := gin.H{
		"plan": PlanResponse{
			ID:            plan.ID,
			Title:         plan.Title,
			Description:   plan.Description,
			Destination:   plan.Destination,
			AuthorID:      plan.AuthorID,
			Status:        plan.Status,
			RatingAverage: avgRating,
			RatingCount:   plan.RatingCount,
			CommentCount:  plan.CommentCount,
			NodeCount:     nodeCount,
			CreatedAt:     plan.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:     plan.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		},
		"nodes": nodes,
	}

	middleware.SuccessResponse(c, http.StatusOK, resp)
}

// CreatePlan handles POST /api/v1/plans - create new draft plan
func (pc *PlanController) CreatePlan(c *gin.Context) {
	var req CreatePlanRequest

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ValidationErrorResponse(c, "invalid request", nil)
		return
	}

	// Get current user from context
	userID, ok := utilities.GetUserIDFromContext(c)
	if !ok {
		middleware.AuthErrorResponse(c, "User not authenticated")
		return
	}

	// Create plan model
	plan := &models.TravelPlan{
		Title:       strings.TrimSpace(req.Title),
		Description: strings.TrimSpace(req.Description),
		Destination: strings.TrimSpace(req.Destination),
		AuthorID:    userID,
		Status:      models.TravelPlanStatusDraft.String(),
	}

	// Call service
	planID, err := pc.planService.CreatePlan(c.Request.Context(), plan)
	if err != nil {
		if err == models.ErrValidation {
			middleware.ValidationErrorResponse(c, "Invalid plan data", nil)
			return
		}
		middleware.InternalErrorResponse(c, "Failed to create plan")
		return
	}

	resp := gin.H{
		"plan_id": planID,
		"status":  models.TravelPlanStatusDraft.String(),
	}

	middleware.SuccessResponse(c, http.StatusCreated, resp)
}

// PublishPlan handles PATCH /api/v1/plans/:id/publish - publish a draft plan
func (pc *PlanController) PublishPlan(c *gin.Context) {
	planID := c.Param("id")

	// Get current user
	userID, ok := utilities.GetUserIDFromContext(c)
	if !ok {
		middleware.AuthErrorResponse(c, "User not authenticated")
		return
	}

	// Verify ownership
	plan, _ := pc.planService.GetPlanByID(c.Request.Context(), planID)
	if plan == nil {
		middleware.NotFoundErrorResponse(c, "Plan not found")
		return
	}

	if plan.AuthorID != userID {
		middleware.ForbiddenErrorResponse(c, "You do not have permission to publish this plan")
		return
	}

	// Publish plan
	err := pc.planService.PublishPlan(c.Request.Context(), planID)
	if err != nil {
		if err == models.ErrValidation {
			middleware.ValidationErrorResponse(c, "Plan is not in draft status", nil)
			return
		}
		middleware.InternalErrorResponse(c, "Failed to publish plan")
		return
	}

	middleware.SuccessResponse(c, http.StatusOK, gin.H{"status": "published"})
}
