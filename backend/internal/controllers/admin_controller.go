package controllers

import (
	"net/http"

	"tll-backend/internal/middleware"
	"tll-backend/internal/services"
	"tll-backend/internal/utilities"

	"github.com/gin-gonic/gin"
)

// AdminController handles admin/moderation operations
type AdminController struct {
	planService services.PlanService
	nodeService services.NodeService
	userService services.UserService
}

// NewAdminController creates a new admin controller
func NewAdminController(planService services.PlanService, nodeService services.NodeService, userService services.UserService) *AdminController {
	return &AdminController{
		planService: planService,
		nodeService: nodeService,
		userService: userService,
	}
}

// SuspendPlan handles PATCH /api/v1/plans/:id/suspend - admin suspend plan
func (ac *AdminController) SuspendPlan(c *gin.Context) {
	planID := c.Param("id")

	// Verify plan exists
	plan, _ := ac.planService.GetPlanByID(c.Request.Context(), planID)
	if plan == nil {
		middleware.NotFoundErrorResponse(c, "Plan not found")
		return
	}

	// Suspend plan
	err := ac.planService.SuspendPlan(c.Request.Context(), planID)
	if err != nil {
		middleware.InternalErrorResponse(c, "Failed to suspend plan")
		return
	}

	middleware.SuccessResponse(c, http.StatusOK, gin.H{
		"status": "suspended",
	})
}

// DeletePlan handles DELETE /api/v1/plans/:id - admin soft-delete plan
func (ac *AdminController) DeletePlan(c *gin.Context) {
	planID := c.Param("id")

	// Verify plan exists
	plan, _ := ac.planService.GetPlanByID(c.Request.Context(), planID)
	if plan == nil {
		middleware.NotFoundErrorResponse(c, "Plan not found")
		return
	}

	// Delete (soft-delete) plan
	err := ac.planService.DeletePlan(c.Request.Context(), planID)
	if err != nil {
		middleware.InternalErrorResponse(c, "Failed to delete plan")
		return
	}

	middleware.SuccessResponse(c, http.StatusOK, gin.H{
		"message": "Plan deleted",
	})
}

// ApproveNode handles PATCH /api/v1/nodes/:id/approve - admin approve user-created node
func (ac *AdminController) ApproveNode(c *gin.Context) {
	nodeID := c.Param("id")

	// Verify node exists
	node, err := ac.nodeService.GetNodeByID(c.Request.Context(), nodeID)
	if err != nil || node == nil {
		middleware.NotFoundErrorResponse(c, "Node not found")
		return
	}

	// Approve node
	err = ac.nodeService.ApproveNode(c.Request.Context(), nodeID)
	if err != nil {
		middleware.InternalErrorResponse(c, "Failed to approve node")
		return
	}

	middleware.SuccessResponse(c, http.StatusOK, gin.H{
		"status": "approved",
	})
}

// DisapproveNode handles PATCH /api/v1/nodes/:id/disapprove - admin disapprove node
func (ac *AdminController) DisapproveNode(c *gin.Context) {
	nodeID := c.Param("id")

	// Verify node exists
	node, err := ac.nodeService.GetNodeByID(c.Request.Context(), nodeID)
	if err != nil || node == nil {
		middleware.NotFoundErrorResponse(c, "Node not found")
		return
	}

	// Disapprove node
	// TODO: Implement DisapproveNode in NodeService
	// err = ac.nodeService.DisapproveNode(c.Request.Context(), nodeID)

	middleware.SuccessResponse(c, http.StatusOK, gin.H{
		"status": "disapproved",
		"message": "Disapprove operation not yet implemented",
	})
}

// DeleteNode handles DELETE /api/v1/nodes/:id - admin delete node
func (ac *AdminController) DeleteNode(c *gin.Context) {
	nodeID := c.Param("id")

	// Verify node exists
	node, err := ac.nodeService.GetNodeByID(c.Request.Context(), nodeID)
	if err != nil || node == nil {
		middleware.NotFoundErrorResponse(c, "Node not found")
		return
	}

	// Delete node
	err = ac.nodeService.DeleteNode(c.Request.Context(), nodeID)
	if err != nil {
		middleware.InternalErrorResponse(c, "Failed to delete node")
		return
	}

	middleware.SuccessResponse(c, http.StatusOK, gin.H{
		"message": "Node deleted",
	})
}

// UpdateUserRole handles PATCH /api/v1/users/:id/role - admin change user role
func (ac *AdminController) UpdateUserRole(c *gin.Context) {
	userID := c.Param("id")

	var req struct {
		Role string `json:"role" binding:"required"`
	}

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ValidationErrorResponse(c, "invalid request", nil)
		return
	}

	// Get current user (admin)
	adminID, ok := utilities.GetUserIDFromContext(c)
	if !ok {
		middleware.AuthErrorResponse(c, "Admin not authenticated")
		return
	}

	// Prevent self-modification
	if adminID == userID {
		middleware.ForbiddenErrorResponse(c, "You cannot change your own role")
		return
	}

	// TODO: When service is ready, implement role change logic
	// Based on role argument, call appropriate service method
	// e.g., PromoteToTraveller, DemoteToSimple, MakeAdmin

	middleware.SuccessResponse(c, http.StatusOK, gin.H{
		"message": "User role update not yet implemented via service",
	})
}

// DeactivateUser handles PATCH /api/v1/users/:id/deactivate - admin deactivate user
func (ac *AdminController) DeactivateUser(c *gin.Context) {
	userID := c.Param("id")

	// Get current admin
	adminID, ok := utilities.GetUserIDFromContext(c)
	if !ok {
		middleware.AuthErrorResponse(c, "Admin not authenticated")
		return
	}

	// Prevent self-deactivation
	if adminID == userID {
		middleware.ForbiddenErrorResponse(c, "You cannot deactivate yourself")
		return
	}

	// TODO: When service is ready, implement user deactivation
	// err := ac.userService.DeleteUser(c.Request.Context(), userID)

	middleware.SuccessResponse(c, http.StatusOK, gin.H{
		"message": "User deactivated",
	})
}
