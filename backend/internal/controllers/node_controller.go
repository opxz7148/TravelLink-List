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

// NodeController handles node operations (attractions and transitions)
type NodeController struct {
	nodeService services.NodeService
}

// NewNodeController creates a new node controller
func NewNodeController(nodeService services.NodeService) *NodeController {
	return &NodeController{
		nodeService: nodeService,
	}
}

// CreateAttractionRequest represents attraction creation request
type CreateAttractionRequest struct {
	Name                      string `json:"name" binding:"required,min=1,max=200"`
	Category                  string `json:"category" binding:"required"`
	Location                  string `json:"location" binding:"required,min=1,max=300"`
	Description               string `json:"description" binding:"max=1000"`
	ContactInfo               string `json:"contact_info" binding:"max=200"`
	HoursOfOperation          string `json:"hours_of_operation" binding:"max=200"`
	EstimatedVisitDurationMin int    `json:"estimated_visit_duration_minutes"`
}

// CreateTransitionRequest represents transition creation request
type CreateTransitionRequest struct {
	Title            string `json:"title" binding:"required,max=200"`
	Mode             string `json:"mode" binding:"required"`
	Description      string `json:"description" binding:"max=1000"`
	HoursOfOperation string `json:"hours_of_operation" binding:"max=200"`
	RouteNotes       string `json:"route_notes" binding:"max=500"`
}

// ListNodes handles GET /api/v1/nodes - list available nodes
// @Summary List available nodes with embedded details
// @Description Get paginated list of approved nodes (attractions and/or transitions) with embedded type-specific details. Each node includes complete information: base properties and either attraction or transition details populated. Can filter by type. (public endpoint)
// @Tags nodes
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Results per page" default(50)
// @Param type query string false "Node type filter" Enums(attraction,transition)
// @Success 200 {object} map[string]interface{} "Array of nodes with embedded details and pagination metadata (total_items, total_pages, current_page, items_per_page)"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /nodes [get]
func (nc *NodeController) ListNodes(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "50")
	nodeType := c.Query("type") // attraction | transition

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}

	offset := (page - 1) * limit

	// Fetch nodes based on type
	var nodes []*models.Node
	var err error

	if nodeType == "attraction" {
		// Empty string for category = all categories
		nodes, _, err = nc.nodeService.ListApprovedAttractions(c.Request.Context(), "", offset, limit)
	} else if nodeType == "transition" {
		// Empty string for mode = all modes
		nodes, _, err = nc.nodeService.ListApprovedTransitions(c.Request.Context(), "", offset, limit)
	} else {
		nodes, err = nc.nodeService.ListApprovedNodes(c.Request.Context(), offset, limit)
	}

	if err != nil {
		middleware.InternalErrorResponse(c, "Failed to fetch nodes")
		return
	}

	middleware.SuccessResponse(c, http.StatusOK, gin.H{
		"nodes": nodes,
		"pagination": gin.H{
			"current_page": page,
			"limit":        limit,
		},
	})
}

// GetNodeDetail handles GET /api/v1/nodes/:id - get node details
// @Summary Get node details with embedded type-specific information
// @Description Retrieve complete node information including: base properties (id, type, created_by, is_approved, created_at), and embedded type-specific details. For attractions: name, category, location, description, contact_info, hours_of_operation. For transitions: title, mode, description, hours_of_operation, route_notes. The "attraction" field is populated for attraction nodes, "transition" for transition nodes (null otherwise). (public endpoint)
// @Tags nodes
// @Produce json
// @Param id path string true "Node ID"
// @Success 200 {object} map[string]interface{} "Complete node with embedded details. For attractions: {id, type, created_by, is_approved, attraction: {node_id, name, category, location, ...}, transition: null}. For transitions: similar with transition object populated."
// @Failure 404 {object} middleware.SwaggerErrorResponse "Node not found"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /nodes/{id} [get]
func (nc *NodeController) GetNodeDetail(c *gin.Context) {
	nodeID := c.Param("id")

	// Fetch node
	node, err := nc.nodeService.GetNodeByID(c.Request.Context(), nodeID)
	if err != nil {
		middleware.InternalErrorResponse(c, "Failed to fetch node")
		return
	}

	if node == nil {
		middleware.NotFoundErrorResponse(c, "Node not found")
		return
	}

	middleware.SuccessResponse(c, http.StatusOK, gin.H{"node": node})
}

// CreateAttractionNode handles POST /api/v1/nodes/attraction - create attraction node
// @Summary Create an attraction node
// @Description Create a new user-generated attraction node. Requires admin approval before public use. Traveller or admin only.
// @Tags nodes
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body CreateAttractionRequest true "Attraction creation request"
// @Success 201 {object} map[string]interface{} "Attraction node created (pending approval)"
// @Failure 400 {object} middleware.SwaggerErrorResponse "Validation error"
// @Failure 401 {object} middleware.SwaggerErrorResponse "Not authenticated"
// @Failure 403 {object} middleware.SwaggerErrorResponse "Only traveller or admin can create nodes"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /nodes/attraction [post]
func (nc *NodeController) CreateAttractionNode(c *gin.Context) {
	var req CreateAttractionRequest

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

	// Build attraction detail object
	detail := &models.AttractionNodeDetail{
		Name:             req.Name,
		Category:         req.Category,
		Location:         req.Location,
		Description:      req.Description,
		ContactInfo:      req.ContactInfo,
		HoursOfOperation: req.HoursOfOperation,
	}

	// Call service
	nodeID, err := nc.nodeService.CreateAttractionNode(c.Request.Context(), userID, detail)
	if err != nil {
		if err == models.ErrValidation {
			middleware.ValidationErrorResponse(c, "Invalid attraction data", nil)
			return
		}
		middleware.InternalErrorResponse(c, "Failed to create attraction")
		return
	}

	middleware.SuccessResponse(c, http.StatusCreated, gin.H{
		"node_id": nodeID,
		"type":    "attraction",
	})
}

// CreateTransitionNode handles POST /api/v1/nodes/transition - create transition node
// @Summary Create a transition node
// @Description Create a new user-generated transition node (journey between attractions). Requires admin approval before public use. Traveller or admin only.
// @Description Transition nodes represent movement between attractions with modes like walking, car, bus, train, bike, taxi, or flight.
// @Description Fields: title (service/line identifier like "Bus Line 5"), mode (transportation type), hours (operating hours if applicable), description, route_notes, distance
// @Tags nodes
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body CreateTransitionRequest true "Transition creation request"
// @Success 201 {object} map[string]interface{} "Transition node created (pending approval)"
// @Failure 400 {object} middleware.SwaggerErrorResponse "Validation error"
// @Failure 401 {object} middleware.SwaggerErrorResponse "Not authenticated"
// @Failure 403 {object} middleware.SwaggerErrorResponse "Only traveller or admin can create nodes"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /nodes/transition [post]
func (nc *NodeController) CreateTransitionNode(c *gin.Context) {
	var req CreateTransitionRequest

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ValidationErrorResponse(c, err.Error(), nil)
		return
	}

	// Get current user
	userID, ok := utilities.GetUserIDFromContext(c)
	if !ok {
		middleware.AuthErrorResponse(c, "User not authenticated")
		return
	}

	// Build transition detail object
	var hoursOfOp *string
	if req.HoursOfOperation != "" {
		hoursOfOp = &req.HoursOfOperation
	}
	detail := &models.TransitionNodeDetail{
		Title:            req.Title,
		Mode:             req.Mode,
		Description:      req.Description,
		HoursOfOperation: hoursOfOp,
		RouteNotes:       req.RouteNotes,
	}

	// Call service
	nodeID, err := nc.nodeService.CreateTransitionNode(c.Request.Context(), userID, detail)
	if err != nil {
		if err == models.ErrValidation {
			middleware.ValidationErrorResponse(c, "Invalid transition data", nil)
			return
		}
		middleware.InternalErrorResponse(c, "Failed to create transition")
		return
	}

	middleware.SuccessResponse(c, http.StatusCreated, gin.H{
		"node_id": nodeID,
		"type":    "transition",
	})
}
