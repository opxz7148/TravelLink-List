package controllers

import (
	"net/http"
	"strconv"

	"tll-backend/internal/logger"
	"tll-backend/internal/middleware"
	"tll-backend/internal/models"
	"tll-backend/internal/services"
	"tll-backend/internal/utilities"

	"github.com/gin-gonic/gin"
)

// AdminController handles admin/moderation operations
type AdminController struct {
	planService      services.PlanService
	nodeService      services.NodeService
	userService      services.UserService
	promotionService services.PromotionService
}

// NewAdminController creates a new admin controller
func NewAdminController(planService services.PlanService, nodeService services.NodeService, userService services.UserService, promotionService services.PromotionService) *AdminController {
	return &AdminController{
		planService:      planService,
		nodeService:      nodeService,
		userService:      userService,
		promotionService: promotionService,
	}
}

// SuspendPlan handles PATCH /api/v1/plans/:id/suspend - admin suspend plan
// @Summary Suspend a travel plan
// @Description Admin can suspend a travel plan (admin only)
// @Tags admin
// @Security Bearer
// @Produce json
// @Param id path string true "Plan ID"
// @Success 200 {object} map[string]string "Plan suspended successfully"
// @Failure 404 {object} middleware.SwaggerErrorResponse "Plan not found"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /plans/{id}/suspend [patch]
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
// @Summary Delete a travel plan
// @Description Admin can soft-delete a travel plan (admin only)
// @Tags admin
// @Security Bearer
// @Produce json
// @Param id path string true "Plan ID"
// @Success 200 {object} map[string]string "Plan deleted successfully"
// @Failure 404 {object} middleware.SwaggerErrorResponse "Plan not found"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /plans/{id} [delete]
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
// @Summary Approve a user-created node
// @Description Admin can approve a user-created node for inclusion in travel plans (admin only)
// @Tags admin
// @Security Bearer
// @Produce json
// @Param id path string true "Node ID"
// @Success 200 {object} map[string]string "Node approved successfully"
// @Failure 404 {object} middleware.SwaggerErrorResponse "Node not found"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /nodes/{id}/approve [patch]
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
// @Summary Disapprove a node
// @Description Admin can disapprove a node, removing it from travel plans (admin only)
// @Tags admin
// @Security Bearer
// @Produce json
// @Param id path string true "Node ID"
// @Success 200 {object} map[string]interface{} "Node disapproved successfully"
// @Failure 404 {object} middleware.SwaggerErrorResponse "Node not found"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /nodes/{id}/disapprove [patch]
func (ac *AdminController) DisapproveNode(c *gin.Context) {
	nodeID := c.Param("id")

	// Verify node exists
	node, err := ac.nodeService.GetNodeByID(c.Request.Context(), nodeID)
	if err != nil || node == nil {
		middleware.NotFoundErrorResponse(c, "Node not found")
		return
	}

	// Disapprove node via service
	err = ac.nodeService.DisapproveNode(c.Request.Context(), nodeID)
	if err != nil {
		middleware.InternalErrorResponse(c, "Failed to disapprove node")
		return
	}

	middleware.SuccessResponse(c, http.StatusOK, gin.H{
		"status":  "disapproved",
		"message": "Node disapproved successfully",
	})
}

// DeleteNode handles DELETE /api/v1/nodes/:id - admin delete node
// @Summary Delete a node
// @Description Admin can delete a node (admin only)
// @Tags admin
// @Security Bearer
// @Produce json
// @Param id path string true "Node ID"
// @Success 200 {object} map[string]string "Node deleted successfully"
// @Failure 404 {object} middleware.SwaggerErrorResponse "Node not found"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /nodes/{id} [delete]
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
// @Summary Change user role
// @Description Admin can change a user's role (simple/traveller/admin). Admin cannot change admin's own role. (admin only)
// @Tags admin
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body map[string]string true "Role change request (role: simple|traveller|admin)"
// @Success 200 {object} map[string]interface{} "User role updated successfully"
// @Failure 400 {object} middleware.SwaggerErrorResponse "Invalid role or cannot change own role"
// @Failure 401 {object} middleware.SwaggerErrorResponse "Not authenticated"
// @Failure 403 {object} middleware.SwaggerErrorResponse "Cannot change own role"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /users/{id}/role [patch]
func (ac *AdminController) UpdateUserRole(c *gin.Context) {
	userID := c.Param("id")
	log := logger.GetLogger("AdminController")

	var req struct {
		Role string `json:"role" binding:"required"`
	}

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		log.ValidationError(c.Request.Context(), map[string]string{"request": "Invalid request format"})
		middleware.ValidationErrorResponse(c, "invalid request", nil)
		return
	}

	// Get current user (admin)
	adminID, ok := utilities.GetUserIDFromContext(c)
	if !ok {
		log.AuthorizationError(c.Request.Context(), "", "Role update attempted without authentication")
		middleware.AuthErrorResponse(c, "Admin not authenticated")
		return
	}

	// Prevent self-modification
	if adminID == userID {
		log.AuthorizationError(c.Request.Context(), adminID, "Attempted to change own role")
		middleware.ForbiddenErrorResponse(c, "You cannot change your own role")
		return
	}

	// Validate role
	if !models.CheckRole(models.UserRole(req.Role)) {
		log.ValidationError(c.Request.Context(), map[string]string{"role": "Invalid role value"})
		middleware.ValidationErrorResponse(c, "Invalid role", gin.H{"field": "role"})
		return
	}

	// Update user role via service based on requested role
	var err error
	switch models.UserRole(req.Role) {
	case models.RoleSimple:
		err = ac.userService.DemoteToSimple(c.Request.Context(), userID)
	case models.RoleTraveller:
		err = ac.userService.PromoteToTraveller(c.Request.Context(), userID)
	case models.RoleAdmin:
		err = ac.userService.MakeAdmin(c.Request.Context(), userID)
	}

	if err != nil {
		log.ServiceError(c.Request.Context(), "UserService", "UpdateUserRole", err, adminID)
		middleware.InternalErrorResponse(c, "Failed to update user role")
		return
	}

	log.SecurityEvent(c.Request.Context(), "USER_ROLE_CHANGED", adminID, "", "Changed role of user "+userID+" to "+req.Role)

	middleware.SuccessResponse(c, http.StatusOK, gin.H{
		"user_id": userID,
		"role":    req.Role,
		"message": "User role updated successfully",
	})
}

// DeactivateUser handles PATCH /api/v1/users/:id/deactivate - admin deactivate user
// @Summary Deactivate user account
// @Description Admin can soft-delete (deactivate) a user account. Admin cannot deactivate self. (admin only)
// @Tags admin
// @Security Bearer
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{} "User deactivated successfully"
// @Failure 401 {object} middleware.SwaggerErrorResponse "Not authenticated"
// @Failure 403 {object} middleware.SwaggerErrorResponse "Cannot deactivate own account"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /users/{id}/deactivate [patch]
func (ac *AdminController) DeactivateUser(c *gin.Context) {
	userID := c.Param("id")
	log := logger.GetLogger("AdminController")

	// Get current admin
	adminID, ok := utilities.GetUserIDFromContext(c)
	if !ok {
		log.AuthorizationError(c.Request.Context(), "", "User deactivation attempted without authentication")
		middleware.AuthErrorResponse(c, "Admin not authenticated")
		return
	}

	// Prevent self-deactivation
	if adminID == userID {
		log.AuthorizationError(c.Request.Context(), adminID, "Attempted to deactivate own account")
		middleware.ForbiddenErrorResponse(c, "You cannot deactivate yourself")
		return
	}

	// Deactivate user via service (soft delete)
	err := ac.userService.DeleteUser(c.Request.Context(), userID)
	if err != nil {
		log.ServiceError(c.Request.Context(), "UserService", "DeleteUser", err, adminID)
		middleware.InternalErrorResponse(c, "Failed to deactivate user")
		return
	}

	log.SecurityEvent(c.Request.Context(), "USER_DEACTIVATED", adminID, "", "Admin deactivated user "+userID)

	middleware.SuccessResponse(c, http.StatusOK, gin.H{
		"user_id": userID,
		"message": "User deactivated successfully",
	})
}

// ApprovePromotionRequest handles PATCH /api/v1/admin/promotions/:id/approve - approve promotion request
// @Summary Approve a promotion request
// @Description Admin can approve a promotion request. If no plan is specified, user role is upgraded to traveller and all user's draft nodes are auto-approved. If a plan is specified, user is promoted to traveller and the plan is published with nodes auto-approved.
// @Tags admin
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "Promotion Request ID"
// @Param request body map[string]string true "Admin notes (max 500 chars)"
// @Success 200 {object} map[string]interface{} "Promotion request approved"
// @Failure 404 {object} middleware.SwaggerErrorResponse "Promotion request not found"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /admin/promotions/{id}/approve [patch]
func (ac *AdminController) ApprovePromotionRequest(c *gin.Context) {
	requestID := c.Param("id")
	log := logger.GetLogger("AdminController")

	var req struct {
		AdminNotes string `json:"admin_notes"`
	}

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ValidationErrorResponse(c, "invalid request", nil)
		return
	}

	// Get the promotion request to get user ID and plan ID
	promotionReq, err := ac.promotionService.GetRequest(c.Request.Context(), requestID)
	if err != nil || promotionReq == nil {
		middleware.NotFoundErrorResponse(c, "Promotion request not found")
		return
	}

	userID := promotionReq.UserID

	// Approve the request
	err = ac.promotionService.ApproveRequest(c.Request.Context(), requestID, req.AdminNotes)
	if err != nil {
		if err == models.ErrNotFound {
			middleware.NotFoundErrorResponse(c, "Promotion request not found")
			return
		}
		if err == models.ErrValidation {
			middleware.ValidationErrorResponse(c, "Request cannot be approved (invalid status)", nil)
			return
		}
		middleware.InternalErrorResponse(c, "Failed to approve promotion request")
		return
	}

	// Always promote user to traveller role
	err = ac.userService.PromoteToTraveller(c.Request.Context(), userID)
	if err != nil {
		// Log error but don't fail - still send success response
		log.ServiceError(c.Request.Context(), "UserService", "PromoteToTraveller", err, userID)
	}

	// If plan is specified (plan promotion request), publish the plan
	if promotionReq.PlanID != nil && *promotionReq.PlanID != "" {
		planID := *promotionReq.PlanID

		// Publish the plan
		err = ac.planService.PublishPlan(c.Request.Context(), planID)
		if err != nil {
			// Log error but don't fail - user is already promoted
			log.ServiceError(c.Request.Context(), "PlanService", "PublishPlan", err, planID)
		}

		// Get all nodes in the plan and auto-approve user's draft nodes
		planNodes, err := ac.planService.GetPlanNodes(c.Request.Context(), planID)
		if err != nil {
			// Log error but don't fail - plan is already published
			c.Error(err)
		} else {
			// Auto-approve draft nodes created by this user in this plan
			for _, planNode := range planNodes {
				if planNode.NodeID != "" {
					node, err := ac.nodeService.GetNodeByID(c.Request.Context(), planNode.NodeID)
					if err == nil && node != nil && node.CreatedBy == userID && !node.IsApproved {
						if err := ac.nodeService.ApproveNode(c.Request.Context(), node.ID); err != nil {
							// Log error but don't fail - continue with other nodes
							c.Error(err)
						}
					}
				}
			}
		}
	} else {
		// Role upgrade request (no plan specified)
		// Get all user's draft nodes and auto-approve them
		draftNodes, err := ac.nodeService.ListDraftNodesByCreator(c.Request.Context(), userID, 0, 1000)
		if err != nil {
			// Log error but don't fail - promotion is already approved
			c.Error(err)
		} else {
			// Auto-approve each draft node
			for _, node := range draftNodes {
				if err := ac.nodeService.ApproveNode(c.Request.Context(), node.ID); err != nil {
					// Log error but don't fail - continue with other nodes
					c.Error(err)
				}
			}
		}
	}

	middleware.SuccessResponse(c, http.StatusOK, gin.H{
		"request_id": requestID,
		"status":     "approved",
	})
}

// ListPendingPromotions handles GET /api/v1/admin/promotions/pending - list pending promotion requests
// @Summary List pending promotion requests
// @Description Admin can view all pending promotion requests awaiting review (admin only)
// @Tags admin
// @Security Bearer
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Results per page" default(50)
// @Success 200 {object} map[string]interface{} "List of pending promotion requests"
// @Failure 401 {object} middleware.SwaggerErrorResponse "Not authenticated"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /admin/promotions/pending [get]
func (ac *AdminController) ListPendingPromotions(c *gin.Context) {
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

	// Get pending promotion requests
	requests, total, err := ac.promotionService.ListPending(c.Request.Context(), offset, limit)
	if err != nil {
		middleware.InternalErrorResponse(c, "Failed to fetch pending promotion requests")
		return
	}

	// Enrich requests with user and plan details
	type PromotionRequestDetailResponse struct {
		ID         string             `json:"id"`
		UserID     string             `json:"user_id"`
		User       *models.User       `json:"user,omitempty"`
		PlanID     *string            `json:"plan_id"`
		Plan       *models.TravelPlan `json:"plan,omitempty"`
		Status     string             `json:"status"`
		AdminNotes string             `json:"admin_notes"`
		CreatedAt  string             `json:"created_at"`
		ReviewedAt *string            `json:"reviewed_at"`
	}

	responses := make([]PromotionRequestDetailResponse, len(requests))
	for i, req := range requests {
		reviewedAtStr := ""
		if req.ReviewedAt != nil {
			reviewedAtStr = req.ReviewedAt.String()
		}

		// Get user details
		user, _ := ac.userService.GetUserByID(c.Request.Context(), req.UserID)

		// Get plan details if plan_id is specified
		var plan *models.TravelPlan
		if req.PlanID != nil && *req.PlanID != "" {
			plan, _ = ac.planService.GetPlanByID(c.Request.Context(), *req.PlanID)
		}

		responses[i] = PromotionRequestDetailResponse{
			ID:         req.ID,
			UserID:     req.UserID,
			User:       user,
			PlanID:     req.PlanID,
			Plan:       plan,
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

// RejectPromotionRequest handles PATCH /api/v1/admin/promotions/:id/reject - reject promotion request
// @Summary Reject a promotion request
// @Description Admin can reject a promotion request with optional notes. (admin only)
// @Tags admin
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "Promotion Request ID"
// @Param request body map[string]string true "Admin notes (max 500 chars)"
// @Success 200 {object} map[string]interface{} "Promotion request rejected"
// @Failure 404 {object} middleware.SwaggerErrorResponse "Promotion request not found"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /admin/promotions/{id}/reject [patch]
func (ac *AdminController) RejectPromotionRequest(c *gin.Context) {
	requestID := c.Param("id")

	var req struct {
		AdminNotes string `json:"admin_notes"`
	}

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ValidationErrorResponse(c, "invalid request", nil)
		return
	}

	// Reject the request
	err := ac.promotionService.RejectRequest(c.Request.Context(), requestID, req.AdminNotes)
	if err != nil {
		if err == models.ErrNotFound {
			middleware.NotFoundErrorResponse(c, "Promotion request not found")
			return
		}
		if err == models.ErrValidation {
			middleware.ValidationErrorResponse(c, "Request cannot be rejected (invalid status)", nil)
			return
		}
		middleware.InternalErrorResponse(c, "Failed to reject promotion request")
		return
	}

	middleware.SuccessResponse(c, http.StatusOK, gin.H{
		"request_id": requestID,
		"status":     "rejected",
	})
}
