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
	Mode                     string `json:"mode" binding:"required"`
	EstimatedDurationMinutes int    `json:"estimated_duration_minutes" binding:"required,gt=0"`
	RouteNotes               string `json:"route_notes" binding:"max=500"`
	EstimatedDistanceKm      float64 `json:"estimated_distance_km"`
}

// ListNodes handles GET /api/v1/nodes - list available nodes
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
		Name:                          req.Name,
		Category:                      req.Category,
		Location:                      req.Location,
		Description:                   req.Description,
		ContactInfo:                   req.ContactInfo,
		HoursOfOperation:              req.HoursOfOperation,
		EstimatedVisitDurationMinutes: &req.EstimatedVisitDurationMin,
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
func (nc *NodeController) CreateTransitionNode(c *gin.Context) {
	var req CreateTransitionRequest

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

	// Build transition detail object
	distance := req.EstimatedDistanceKm
	detail := &models.TransitionNodeDetail{
		Mode:                    req.Mode,
		EstimatedDurationMinutes: req.EstimatedDurationMinutes,
		RouteNotes:              req.RouteNotes,
		EstimatedDistanceKm:     &distance,
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
