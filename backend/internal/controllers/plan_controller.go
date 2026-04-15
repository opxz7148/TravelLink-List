package controllers

import (
	"fmt"
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
	userService services.UserService
}

// NewPlanController creates a new plan controller
func NewPlanController(planService services.PlanService, nodeService services.NodeService, userService services.UserService) *PlanController {
	return &PlanController{
		planService: planService,
		nodeService: nodeService,
		userService: userService,
	}
}

// PlanResponse represents plan data in API response
type PlanResponse struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Destination   string `json:"destination"`
	AuthorID      string `json:"author_id"`
	Author        *AuthorResponse `json:"author,omitempty"`
	Status        string `json:"status"`
	RatingAverage float64 `json:"rating_average"`
	RatingSum     int `json:"rating_sum"`
	RatingCount   int `json:"rating_count"`
	CommentCount  int `json:"comment_count"`
	NodeCount     int `json:"node_count"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// AuthorResponse represents minimal author information
type AuthorResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
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

// autoApproveUserDraftNodes auto-approves draft nodes created by the user in a plan
// Call this when a plan is published (either via PublishPlan endpoint or CreatePlan with status=published)
func (pc *PlanController) autoApproveUserDraftNodes(c *gin.Context, planID, userID string) {
	fmt.Println("\n=== AUTO-APPROVAL DEBUG ===")
	fmt.Println("Plan ID:", planID)
	fmt.Println("User ID:", userID)

	planNodes, err := pc.planService.GetPlanNodes(c.Request.Context(), planID)
	if err != nil {
		fmt.Println("❌ ERROR fetching plan nodes:", err)
		c.Error(err)
		return
	}

	fmt.Println("✓ FETCHED plan nodes, count:", len(planNodes))

	if len(planNodes) == 0 {
		fmt.Println("⚠️  No nodes found in plan!")
		fmt.Println("=== AUTO-APPROVAL DEBUG END ===\n")
		return
	}

	// Auto-approve draft nodes created by this user
	for i, planNode := range planNodes {
		fmt.Println("\n--- Checking PlanNode #" + fmt.Sprint(i+1) + " ---")

		if planNode == nil {
			fmt.Println("⚠️  PlanNode is nil, skipping")
			continue
		}

		fmt.Println("PlanNode ID:", planNode.ID)
		fmt.Println("PlanNode.NodeID:", planNode.NodeID)

		if planNode.NodeID == "" {
			fmt.Println("⚠️  PlanNode.NodeID is empty, skipping")
			continue
		}

		node, err := pc.nodeService.GetNodeByID(c.Request.Context(), planNode.NodeID)
		if err != nil {
			fmt.Println("❌ ERROR fetching node by ID:", planNode.NodeID, "Error:", err)
			continue
		}
		if node == nil {
			fmt.Println("❌ Node not found for ID:", planNode.NodeID)
			continue
		}

		fmt.Println("✓ Node found:")
		fmt.Println("  - Node.ID:", node.ID)
		fmt.Println("  - Node.CreatedBy:", node.CreatedBy)
		fmt.Println("  - Node.IsApproved:", node.IsApproved)
		fmt.Println("  - User.ID:", userID)
		fmt.Println("  - CreatedBy == UserID?", node.CreatedBy == userID)
		fmt.Println("  - IsApproved == false?", !node.IsApproved)

		// Auto-approve if created by user and not yet approved
		if node.CreatedBy == userID && !node.IsApproved {
			fmt.Println("✓ CONDITION MET - calling ApproveNode...")
			if err := pc.nodeService.ApproveNode(c.Request.Context(), node.ID); err != nil {
				fmt.Println("❌ ERROR approving node:", node.ID, "Error:", err)
				c.Error(err)
			} else {
				fmt.Println("✓ NODE APPROVED SUCCESSFULLY:", node.ID)
			}
		} else {
			if node.CreatedBy != userID {
				fmt.Println("⚠️  CONDITION FAILED: CreatedBy does NOT match UserID")
			}
			if node.IsApproved {
				fmt.Println("⚠️  CONDITION FAILED: Node is already approved")
			}
		}
	}

	fmt.Println("=== AUTO-APPROVAL DEBUG END ===\n")
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
			RatingSum:     plan.RatingSum,
			RatingCount:   plan.RatingCount,
			CommentCount:  plan.CommentCount,
			CreatedAt:     plan.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:     plan.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}

		// Get node count
		nodeCount, _ := pc.planService.CountPlanNodes(c.Request.Context(), plan.ID)
		planResponses[i].NodeCount = nodeCount

		// Get author information
		if author, err := pc.userService.GetUserByID(c.Request.Context(), plan.AuthorID); err == nil && author != nil {
			planResponses[i].Author = &AuthorResponse{
				ID:       author.ID,
				Username: author.Username,
			}
		}
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
			RatingSum:     plan.RatingSum,
			RatingCount:   plan.RatingCount,
			CommentCount:  plan.CommentCount,
			CreatedAt:     plan.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:     plan.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}

		// Get node count
		nodeCount, _ := pc.planService.CountPlanNodes(c.Request.Context(), plan.ID)
		planResponses[i].NodeCount = nodeCount

		// Get author information
		if author, err := pc.userService.GetUserByID(c.Request.Context(), plan.AuthorID); err == nil && author != nil {
			planResponses[i].Author = &AuthorResponse{
				ID:       author.ID,
				Username: author.Username,
			}
		}
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
// @Description Retrieve detailed information about a travel plan including nodes. Published plans are public. Draft plans require authentication and user must be owner or admin.
// @Tags plans
// @Produce json
// @Param id path string true "Plan ID"
// @Success 200 {object} map[string]interface{} "Plan details with enriched nodes"
// @Failure 403 {object} middleware.SwaggerErrorResponse "Access denied to draft plan"
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

	// Authorization check for draft plans
	if plan.Status == models.TravelPlanStatusDraft.String() {
		// Get current user ID from context (may not be authenticated)
		userID, ok := utilities.GetUserIDFromContext(c)
		fmt.Println("====================", userID, ok)

		// If not authenticated, deny access
		if !ok {
			middleware.ForbiddenErrorResponse(c, "Access denied to draft plan")
			return
		}

		// Get user role for admin check
		userRole, _ := utilities.GetUserRoleFromContext(c)

		// Only owner or admin can view draft plans
		if plan.AuthorID != userID && userRole != models.RoleAdmin.String() {
			middleware.ForbiddenErrorResponse(c, "Access denied to draft plan")
			return
		}
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
			RatingSum:     plan.RatingSum,
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

	// Auto-approve user's draft nodes if plan is created with status=published
	if status == "published" {
		pc.autoApproveUserDraftNodes(c, planID, userID)
	}

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
			RatingSum:     createdPlan.RatingSum,
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
// If user is traveller or admin, user-created draft nodes are automatically approved.
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

	// Auto-approve user's draft nodes when plan is published
	pc.autoApproveUserDraftNodes(c, planID, userID)

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

// GetUserPlans handles GET /api/v1/users/me/plans - get current user's plans (draft and published)
// @Summary Get user's travel plans
// @Description Get all travel plans created by the authenticated user (both draft and published)
// @Tags plans
// @Security Bearer
// @Produce json
// @Param offset query int false "Pagination offset" default(0)
// @Param limit query int false "Results per page" default(20)
// @Success 200 {object} map[string]interface{} "User's plans with pagination"
// @Failure 401 {object} middleware.SwaggerErrorResponse "Not authenticated"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /users/me/plans [get]
func (pc *PlanController) GetUserPlans(c *gin.Context) {
	// Get current user
	userID, ok := utilities.GetUserIDFromContext(c)
	if !ok {
		middleware.AuthErrorResponse(c, "User not authenticated")
		return
	}

	// Parse pagination parameters
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "20")

	offset, _ := strconv.Atoi(offsetStr)
	limit, _ := strconv.Atoi(limitStr)

	if offset < 0 {
		offset = 0
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Fetch user's plans via service
	plans, totalCount, err := pc.planService.GetPlansByAuthor(c.Request.Context(), userID, offset, limit)
	if err != nil {
		middleware.InternalErrorResponse(c, "Failed to fetch user plans")
		return
	}

	// Convert to response format
	planResponses := make([]gin.H, len(plans))
	for i, plan := range plans {
		avgRating, _ := pc.planService.GetAverageRating(c.Request.Context(), plan.ID)
		nodeCount, _ := pc.planService.CountPlanNodes(c.Request.Context(), plan.ID)

		planResponses[i] = gin.H{
			"id":             plan.ID,
			"title":          plan.Title,
			"description":    plan.Description,
			"destination":    plan.Destination,
			"author_id":      plan.AuthorID,
			"status":         plan.Status,
			"rating_average": avgRating,
			"rating_count":   plan.RatingCount,
			"comment_count":  plan.CommentCount,
			"node_count":     nodeCount,
			"created_at":     plan.CreatedAt.Format("2006-01-02T15:04:05Z"),
			"updated_at":     plan.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	middleware.SuccessResponse(c, http.StatusOK, gin.H{
		"plans": planResponses,
		"pagination": gin.H{
			"offset":      offset,
			"limit":       limit,
			"total_items": totalCount,
		},
	})
}

// DeletePlan handles DELETE /api/v1/plans/:id - delete a plan
// @Summary Delete a travel plan
// @Description Delete (soft-delete) a travel plan. Plan author or admin only.
// @Tags plans
// @Security Bearer
// @Produce json
// @Param id path string true "Plan ID"
// @Success 200 {object} map[string]interface{} "Plan deleted successfully"
// @Failure 401 {object} middleware.SwaggerErrorResponse "Not authenticated"
// @Failure 403 {object} middleware.SwaggerErrorResponse "Only plan author or admin can delete"
// @Failure 404 {object} middleware.SwaggerErrorResponse "Plan not found"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /plans/{id} [delete]
func (pc *PlanController) DeletePlan(c *gin.Context) {
	planID := c.Param("id")

	// Get current user
	userID, ok := utilities.GetUserIDFromContext(c)
	if !ok {
		middleware.AuthErrorResponse(c, "User not authenticated")
		return
	}

	// Get user role
	userRole, _ := utilities.GetUserRoleFromContext(c)

	// Verify plan exists
	plan, err := pc.planService.GetPlanByID(c.Request.Context(), planID)
	if err != nil {
		middleware.InternalErrorResponse(c, "Failed to fetch plan")
		return
	}

	if plan == nil {
		middleware.NotFoundErrorResponse(c, "Plan not found")
		return
	}

	// Check authorization: plan owner OR admin
	isOwner := plan.AuthorID == userID
	isAdmin := userRole == string(models.RoleAdmin)

	if !isOwner && !isAdmin {
		middleware.ForbiddenErrorResponse(c, "You do not have permission to delete this plan")
		return
	}

	// Delete plan via service
	err = pc.planService.DeletePlan(c.Request.Context(), planID)
	if err != nil {
		middleware.InternalErrorResponse(c, "Failed to delete plan")
		return
	}

	middleware.SuccessResponse(c, http.StatusOK, gin.H{"message": "Plan deleted successfully"})
}

// EditPlan handles PATCH /api/v1/plans/:id - edit plan and replace all nodes
// @Summary Edit a travel plan with node replacement
// @Description Update plan details and replace all associated nodes. Plan author or admin only.
// @Tags plans
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "Plan ID"
// @Param request body CreatePlanRequest true "Plan edit request with full node replacement"
// @Success 200 {object} map[string]interface{} "Updated plan with nodes"
// @Failure 400 {object} middleware.SwaggerErrorResponse "Validation error"
// @Failure 401 {object} middleware.SwaggerErrorResponse "Not authenticated"
// @Failure 403 {object} middleware.SwaggerErrorResponse "Permission denied"
// @Failure 404 {object} middleware.SwaggerErrorResponse "Plan not found"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /plans/{id} [patch]
func (pc *PlanController) EditPlan(c *gin.Context) {
	planID := c.Param("id")
	var req CreatePlanRequest

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

	// Get user role
	userRole, _ := utilities.GetUserRoleFromContext(c)

	// Verify plan exists and ownership
	plan, err := pc.planService.GetPlanByID(c.Request.Context(), planID)
	if err != nil {
		middleware.InternalErrorResponse(c, "Failed to fetch plan")
		return
	}

	if plan == nil {
		middleware.NotFoundErrorResponse(c, "Plan not found")
		return
	}

	// Check authorization: plan owner OR admin
	isOwner := plan.AuthorID == userID
	isAdmin := userRole == string(models.RoleAdmin)

	if !isOwner && !isAdmin {
		middleware.ForbiddenErrorResponse(c, "You do not have permission to edit this plan")
		return
	}

	// Update plan metadata
	plan.Title = strings.TrimSpace(req.Title)
	plan.Description = strings.TrimSpace(req.Description)
	plan.Destination = strings.TrimSpace(req.Destination)

	if err := pc.planService.UpdatePlan(c.Request.Context(), plan); err != nil {
		middleware.InternalErrorResponse(c, "Failed to update plan")
		return
	}

	// Delete all existing nodes for this plan
	if err := pc.planService.DeleteAllPlanNodes(c.Request.Context(), planID); err != nil {
		middleware.InternalErrorResponse(c, "Failed to delete existing nodes")
		return
	}

	// Add new nodes from request (1-indexed)
	for i, nodeDetail := range req.Nodes {
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

	// Fetch updated plan with nodes
	updatedPlan, err := pc.planService.GetPlanByID(c.Request.Context(), planID)
	if err != nil {
		middleware.InternalErrorResponse(c, "Failed to fetch updated plan")
		return
	}

	if updatedPlan == nil {
		middleware.InternalErrorResponse(c, "Updated plan not found")
		return
	}

	// Get nodes for response
	nodes, _ := pc.planService.GetPlanNodes(c.Request.Context(), planID)

	// Get additional stats
	avgRating, _ := pc.planService.GetAverageRating(c.Request.Context(), planID)

	resp := gin.H{
		"plan": PlanResponse{
			ID:            updatedPlan.ID,
			Title:         updatedPlan.Title,
			Description:   updatedPlan.Description,
			Destination:   updatedPlan.Destination,
			AuthorID:      updatedPlan.AuthorID,
			Status:        updatedPlan.Status,
			RatingAverage: avgRating,
			RatingSum:     updatedPlan.RatingSum,
			RatingCount:   updatedPlan.RatingCount,
			CommentCount:  updatedPlan.CommentCount,
			NodeCount:     len(nodes),
			CreatedAt:     updatedPlan.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:     updatedPlan.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		},
		"nodes": nodes,
	}

	middleware.SuccessResponse(c, http.StatusOK, resp)
}
