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
	nodeService services.NodeService
}

// NewPlanController creates a new plan controller
func NewPlanController(planService services.PlanService, nodeService services.NodeService) *PlanController {
	return &PlanController{
		planService: planService,
		nodeService: nodeService,
	}
}

// PlanResponse represents plan data in API response
type PlanResponse struct {
	ID            string  `json:"id"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	Destination   string  `json:"destination"`
	AuthorID      string  `json:"author_id"`
	Status        string  `json:"status"`
	RatingAverage float64 `json:"rating_average"`
	RatingCount   int     `json:"rating_count"`
	CommentCount  int     `json:"comment_count"`
	NodeCount     int     `json:"node_count"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

// PlanNodeResponse represents a node within a plan with enriched details
// Combines PlanNode (plan-specific data) with Node details (type, name, location, etc.)
type PlanNodeResponse struct {
	ID                  string      `json:"id"`                    // PlanNode ID (plan-node association ID)
	SequencePosition    int         `json:"sequence_position"`     // Position in the linked list (1-indexed)
	Type                string      `json:"type"`                  // Node type: "attraction" or "transition"
	Description         *string     `json:"description"`           // Plan-specific notes (optional)
	EstimatedPriceCents *int        `json:"estimated_price_cents"` // Plan-specific cost in cents (optional)
	DurationMinutes     *int        `json:"duration_minutes"`      // Plan-specific duration in minutes (optional)
	Details             interface{} `json:"details"`               // Node detail object (name, location, etc based on type)
}

// NodeDetailRequest represents plan-specific information for a node in a plan
type NodeDetailRequest struct {
	NodeID              string  `json:"node_id" binding:"required"`
	Description         *string `json:"description" binding:"max=500"`         // Plan-specific notes (optional)
	EstimatedPriceCents *int    `json:"estimated_price_cents" binding:"min=0"` // Cost in cents (optional)
	DurationMinutes     *int    `json:"duration_minutes" binding:"min=1"`      // Duration in minutes (optional)
}

// CreatePlanRequest represents plan creation request with node details
type CreatePlanRequest struct {
	Title       string              `json:"title" binding:"required,min=1,max=150"`
	Description string              `json:"description" binding:"max=1000"`
	Destination string              `json:"destination" binding:"required,min=1,max=200"`
	Nodes       []NodeDetailRequest `json:"nodes" binding:"required,min=1"` // Node details with plan-specific info
}

// UpdatePlanRequest represents plan update request
type UpdatePlanRequest struct {
	Title       string `json:"title" binding:"min=1,max=150"`
	Description string `json:"description" binding:"max=1000"`
	Destination string `json:"destination" binding:"min=1,max=200"`
}

// BrowsePlans handles GET /api/v1/plans - browse published plans
// @Summary Browse published travel plans
// @Description Get paginated list of published travel plans with ratings and comments (public endpoint)
// @Tags plans
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Results per page" default(20)
// @Success 200 {object} map[string]interface{} "List of published plans with pagination"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /plans [get]
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
			ID:            plan.ID,
			Title:         plan.Title,
			Description:   plan.Description,
			Destination:   plan.Destination,
			AuthorID:      plan.AuthorID,
			Status:        plan.Status,
			RatingAverage: avgRating,
			RatingCount:   plan.RatingCount,
			CommentCount:  plan.CommentCount,
			CreatedAt:     plan.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:     plan.UpdatedAt.Format("2006-01-02T15:04:05Z"),
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
// @Summary Search published travel plans
// @Description Search published plans by title, description, or destination (public endpoint)
// @Tags plans
// @Produce json
// @Param q query string true "Search query (required)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Results per page" default(20)
// @Success 200 {object} map[string]interface{} "Search results with pagination"
// @Failure 400 {object} middleware.SwaggerErrorResponse "Missing search query"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /plans/search [get]
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
			ID:            plan.ID,
			Title:         plan.Title,
			Description:   plan.Description,
			Destination:   plan.Destination,
			AuthorID:      plan.AuthorID,
			Status:        plan.Status,
			RatingAverage: avgRating,
			RatingCount:   plan.RatingCount,
			CommentCount:  plan.CommentCount,
			CreatedAt:     plan.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:     plan.UpdatedAt.Format("2006-01-02T15:04:05Z"),
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
// @Summary Get travel plan details
// @Description Retrieve detailed information about a published travel plan including nodes (public endpoint)
// @Tags plans
// @Produce json
// @Param id path string true "Plan ID"
// @Success 200 {object} map[string]interface{} "Plan details with enriched nodes"
// @Failure 404 {object} middleware.SwaggerErrorResponse "Plan not found"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /plans/{id} [get]
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

	// Fetch raw plan nodes
	planNodes, _ := pc.planService.GetPlanNodes(c.Request.Context(), plan.ID)

	// Enrich nodes with full node details (type, name, location, etc.)
	enrichedNodes := make([]PlanNodeResponse, 0, len(planNodes))
	for _, pn := range planNodes {
		if pn == nil {
			continue
		}

		// Fetch the full node with preloaded attraction/transition details
		node, err := pc.nodeService.GetNodeByID(c.Request.Context(), pn.NodeID)
		if err != nil || node == nil {
			continue
		}

		// Build enriched response
		enrichedNode := PlanNodeResponse{
			ID:                  pn.ID,
			SequencePosition:    pn.SequencePosition,
			Type:                node.Type,
			Description:         pn.Description,
			EstimatedPriceCents: pn.EstimatedPriceCents,
			DurationMinutes:     pn.DurationMinutes,
		}

		// Add type-specific details
		if node.Type == "attraction" && node.AttractionNodeDetail != nil {
			enrichedNode.Details = map[string]interface{}{
				"name":        node.AttractionNodeDetail.Name,
				"description": node.AttractionNodeDetail.Description,
				"location":    node.AttractionNodeDetail.Location,
				"category":    node.AttractionNodeDetail.Category,
			}
		} else if node.Type == "transition" && node.TransitionNodeDetail != nil {
			enrichedNode.Details = map[string]interface{}{
				"title":       node.TransitionNodeDetail.Title,
				"description": node.TransitionNodeDetail.Description,
			}
		}

		enrichedNodes = append(enrichedNodes, enrichedNode)
	}

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
		"nodes": enrichedNodes,
	}

	middleware.SuccessResponse(c, http.StatusOK, resp)
}

// CreatePlan handles POST /api/v1/plans - create new plan with nodes and plan-specific details
// @Summary Create a new travel plan with nodes and plan-specific information
// @Description Create a travel plan with initial nodes and plan-dependent details (description, price, duration). User must be traveller or admin. Plan status determined by query parameter (draft or published).
// @Tags plans
// @Security Bearer
// @Accept json
// @Produce json
// @Param status query string false "Plan status (draft or published)" default(draft)
// @Param request body CreatePlanRequest true "Plan creation request with nodes and plan-specific details"
// @Success 201 {object} map[string]interface{} "Plan created with ID, details, and nodes added with plan-dependent information"
// @Failure 400 {object} middleware.SwaggerErrorResponse "Validation error or invalid nodes"
// @Failure 401 {object} middleware.SwaggerErrorResponse "Not authenticated"
// @Failure 403 {object} middleware.SwaggerErrorResponse "Only traveller or admin can create plans"
// @Failure 404 {object} middleware.SwaggerErrorResponse "One or more nodes not found"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /plans [post]
func (pc *PlanController) CreatePlan(c *gin.Context) {
	var req CreatePlanRequest

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ValidationErrorResponse(c, "invalid request", nil)
		return
	}

	// Get current user from context (claims are set by middleware)
	userID, ok := utilities.GetUserIDFromContext(c)
	if !ok {
		middleware.AuthErrorResponse(c, "User not authenticated")
		return
	}

	// Get status from query parameter (default: draft)
	status := c.DefaultQuery("status", "draft")
	if status != "draft" && status != "published" {
		middleware.ValidationErrorResponse(c, "status must be 'draft' or 'published'", gin.H{"field": "status"})
		return
	}

	// Create plan model
	plan := &models.TravelPlan{
		Title:       strings.TrimSpace(req.Title),
		Description: strings.TrimSpace(req.Description),
		Destination: strings.TrimSpace(req.Destination),
		AuthorID:    userID,
		Status:      status,
	}

	// Call service to create plan
	planID, err := pc.planService.CreatePlan(c.Request.Context(), plan)
	if err != nil {
		if err == models.ErrValidation {
			middleware.ValidationErrorResponse(c, "Invalid plan data", nil)
			return
		}
		middleware.InternalErrorResponse(c, "Failed to create plan")
		return
	}

	// Add nodes to plan in sequence with plan-specific information
	// Note: sequence_position is 1-indexed, so we pass i+1 (not i)
	for i, nodeDetail := range req.Nodes {
		// Call service to add node with plan-dependent information
		_, err := pc.planService.AddNodeToPlanWithDetails(
			c.Request.Context(),
			planID,
			nodeDetail.NodeID,
			i+1,
			nodeDetail.Description,
			nodeDetail.EstimatedPriceCents,
			nodeDetail.DurationMinutes,
		)
		if err != nil {
			if err == models.ErrNotFound {
				middleware.NotFoundErrorResponse(c, "One or more nodes not found")
				return
			}
			if err == models.ErrValidation {
				middleware.ValidationErrorResponse(c, "Invalid node or node sequence (no consecutive attractions allowed)", nil)
				return
			}
			middleware.InternalErrorResponse(c, "Failed to add node to plan")
			return
		}
	}

	// Fetch complete plan with nodes
	createdPlan, err := pc.planService.GetPlanByID(c.Request.Context(), planID)
	if err != nil {
		middleware.InternalErrorResponse(c, "Failed to fetch created plan")
		return
	}

	if createdPlan == nil {
		middleware.InternalErrorResponse(c, "Created plan not found")
		return
	}

	// Get nodes for response
	nodes, _ := pc.planService.GetPlanNodes(c.Request.Context(), planID)

	// Get additional stats
	avgRating, _ := pc.planService.GetAverageRating(c.Request.Context(), planID)

	resp := gin.H{
		"plan": PlanResponse{
			ID:            createdPlan.ID,
			Title:         createdPlan.Title,
			Description:   createdPlan.Description,
			Destination:   createdPlan.Destination,
			AuthorID:      createdPlan.AuthorID,
			Status:        createdPlan.Status,
			RatingAverage: avgRating,
			RatingCount:   createdPlan.RatingCount,
			CommentCount:  createdPlan.CommentCount,
			NodeCount:     len(nodes),
			CreatedAt:     createdPlan.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:     createdPlan.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		},
		"nodes": nodes,
	}

	middleware.SuccessResponse(c, http.StatusCreated, resp)
}

// PublishPlan handles PATCH /api/v1/plans/:id/publish - publish a draft plan
// @Summary Publish a travel plan
// @Description Change plan status from draft to published. Plan author or admin only.
// @Tags plans
// @Security Bearer
// @Produce json
// @Param id path string true "Plan ID"
// @Success 200 {object} map[string]interface{} "Plan published successfully"
// @Failure 401 {object} middleware.SwaggerErrorResponse "Not authenticated"
// @Failure 403 {object} middleware.SwaggerErrorResponse "Only plan author or admin can publish"
// @Failure 404 {object} middleware.SwaggerErrorResponse "Plan not found"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /plans/{id}/publish [patch]
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

// UpdatePlanNodesRequest represents request to add/reorder/remove nodes from a plan
type UpdatePlanNodesRequest struct {
	Operation string `json:"operation" binding:"required,oneof=add reorder remove"`
	NodeID    string `json:"node_id" binding:"required"`
	Position  int    `json:"position" binding:"min=0"` // Only used for add/reorder
}

// UpdatePlanNodes handles PATCH /api/v1/plans/{id}/nodes - add/reorder/remove nodes from plan
// @Summary Update plan nodes (add/reorder/remove)
// @Description Add, reorder, or remove nodes from a travel plan. Can only be done by plan author or admin.
// @Tags plans
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "Plan ID"
// @Param request body UpdatePlanNodesRequest true "Node operation request"
// @Success 200 {object} map[string]interface{} "Operation completed successfully"
// @Failure 400 {object} middleware.SwaggerErrorResponse "Invalid request or operation"
// @Failure 401 {object} middleware.SwaggerErrorResponse "Not authenticated"
// @Failure 403 {object} middleware.SwaggerErrorResponse "Only plan author or admin can modify plan"
// @Failure 404 {object} middleware.SwaggerErrorResponse "Plan or node not found"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /plans/{id}/nodes [patch]
func (pc *PlanController) UpdatePlanNodes(c *gin.Context) {
	planID := c.Param("id")
	var req UpdatePlanNodesRequest

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

	// Verify plan ownership
	plan, err := pc.planService.GetPlanByID(c.Request.Context(), planID)
	if err != nil || plan == nil {
		middleware.NotFoundErrorResponse(c, "Plan not found")
		return
	}

	if plan.AuthorID != userID {
		middleware.ForbiddenErrorResponse(c, "You do not have permission to modify this plan")
		return
	}

	// Handle operation
	switch req.Operation {
	case "add":
		if req.Position < 0 {
			middleware.ValidationErrorResponse(c, "position must be >= 0", nil)
			return
		}
		_, err := pc.planService.AddNodeToPlan(c.Request.Context(), planID, req.NodeID, req.Position)
		if err != nil {
			if err == models.ErrValidation {
				middleware.ValidationErrorResponse(c, "Invalid node or plan state", nil)
				return
			}
			middleware.InternalErrorResponse(c, "Failed to add node to plan")
			return
		}

	case "reorder":
		if req.Position < 0 {
			middleware.ValidationErrorResponse(c, "position must be >= 0", nil)
			return
		}
		err := pc.planService.ReorderNodeInPlan(c.Request.Context(), planID, req.NodeID, req.Position)
		if err != nil {
			if err == models.ErrValidation {
				middleware.ValidationErrorResponse(c, "Node not found in plan", nil)
				return
			}
			middleware.InternalErrorResponse(c, "Failed to reorder node")
			return
		}

	case "remove":
		err := pc.planService.RemoveNodeFromPlan(c.Request.Context(), planID, req.NodeID)
		if err != nil {
			if err == models.ErrValidation {
				middleware.ValidationErrorResponse(c, "Node not found in plan", nil)
				return
			}
			middleware.InternalErrorResponse(c, "Failed to remove node")
			return
		}

	default:
		middleware.ValidationErrorResponse(c, "unknown operation", nil)
		return
	}

	middleware.SuccessResponse(c, http.StatusOK, gin.H{"message": "Operation completed successfully"})
}
