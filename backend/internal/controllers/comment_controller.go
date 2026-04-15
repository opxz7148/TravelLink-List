package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"tll-backend/internal/middleware"
	"tll-backend/internal/models"
	"tll-backend/internal/services"
	"tll-backend/internal/utilities"

	"github.com/gin-gonic/gin"
)

// CommentController handles comment operations
type CommentController struct {
	planService    services.PlanService
	commentService services.CommentService
}

// NewCommentController creates a new comment controller
func NewCommentController(planService services.PlanService, commentService services.CommentService) *CommentController {
	return &CommentController{
		planService:    planService,
		commentService: commentService,
	}
}

// CreateCommentRequest represents comment creation request
type CreateCommentRequest struct {
	Text string `json:"text" binding:"required,min=1,max=1000"`
}

// UpdateCommentRequest represents comment update request
type UpdateCommentRequest struct {
	Text string `json:"text" binding:"required,min=1,max=1000"`
}

// CommentResponse represents comment data in API response
type CommentResponse struct {
	ID               string      `json:"id"`
	PlanID           string      `json:"plan_id"`
	AuthorID         string      `json:"author_id"`
	Text             string      `json:"text"`
	IsDeletedByAdmin bool        `json:"is_deleted_by_admin"`
	CreatedAt        string      `json:"created_at"`
	UpdatedAt        *string     `json:"updated_at,omitempty"`
	Author           *AuthorInfo `json:"author,omitempty"`
}

// AuthorInfo represents basic author information for display
type AuthorInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// CreateComment handles POST /api/v1/plans/:id/comments - create comment
// @Summary Create a comment on a travel plan
// @Description Add a comment to a travel plan. Authenticated users can comment.
// @Tags comments
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "Plan ID"
// @Param request body CreateCommentRequest true "Comment creation request"
// @Success 201 {object} map[string]interface{} "Comment created with ID"
// @Failure 400 {object} middleware.SwaggerErrorResponse "Validation error"
// @Failure 401 {object} middleware.SwaggerErrorResponse "Not authenticated"
// @Failure 404 {object} middleware.SwaggerErrorResponse "Plan not found"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /plans/{id}/comments [post]
func (cc *CommentController) CreateComment(c *gin.Context) {
	planID := c.Param("id")
	var req CreateCommentRequest

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

	// Verify plan exists
	plan, _ := cc.planService.GetPlanByID(c.Request.Context(), planID)
	if plan == nil {
		middleware.NotFoundErrorResponse(c, "Plan not found")
		return
	}

	// Create comment model
	comment := &models.Comment{
		PlanID:   planID,
		AuthorID: userID,
		Text:     strings.TrimSpace(req.Text),
	}

	// Validate comment
	if err := comment.Validate(); err != nil {
		middleware.ValidationErrorResponse(c, "Invalid comment data", nil)
		return
	}

	// Create comment via service
	commentID, err := cc.commentService.CreateComment(c.Request.Context(), planID, userID, req.Text)
	if err != nil {
		middleware.InternalErrorResponse(c, "Failed to create comment")
		return
	}

	middleware.SuccessResponse(c, http.StatusCreated, gin.H{
		"id":      commentID,
		"message": "Comment created successfully",
	})
}

// GetComments handles GET /api/v1/plans/:id/comments - list comments
// @Summary Get comments for a travel plan
// @Description Retrieve paginated comments for a travel plan (public endpoint)
// @Tags comments
// @Produce json
// @Param id path string true "Plan ID"
// @Param offset query int false "Offset for pagination" default(0)
// @Param limit query int false "Limit for pagination" default(50)
// @Success 200 {object} map[string]interface{} "Comments list with pagination metadata"
// @Failure 404 {object} middleware.SwaggerErrorResponse "Plan not found"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /plans/{id}/comments [get]
func (cc *CommentController) GetComments(c *gin.Context) {
	planID := c.Param("id")

	// Verify plan exists
	plan, _ := cc.planService.GetPlanByID(c.Request.Context(), planID)
	if plan == nil {
		middleware.NotFoundErrorResponse(c, "Plan not found")
		return
	}

	// Get pagination parameters
	offset := 0
	limit := 50
	if o := c.Query("offset"); o != "" {
		_, _ = fmt.Sscanf(o, "%d", &offset)
	}
	if l := c.Query("limit"); l != "" {
		_, _ = fmt.Sscanf(l, "%d", &limit)
	}

	// Fetch comments via service
	comments, total, err := cc.commentService.ListCommentsByPlan(c.Request.Context(), planID, offset, limit)
	if err != nil {
		middleware.InternalErrorResponse(c, "Failed to fetch comments")
		return
	}

	// Convert to response format
	commentResponses := make([]CommentResponse, len(comments))
	for i, comment := range comments {
		var author *AuthorInfo
		if comment.Author != nil {
			author = &AuthorInfo{
				ID:       comment.Author.ID,
				Username: comment.Author.Username,
			}
		}
		commentResponses[i] = CommentResponse{
			ID:               comment.ID,
			PlanID:           comment.PlanID,
			AuthorID:         comment.AuthorID,
			Text:             comment.Text,
			IsDeletedByAdmin: comment.IsDeletedByAdmin,
			CreatedAt:        comment.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			Author:           author,
		}
		if comment.UpdatedAt != nil {
			updatedAt := comment.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
			commentResponses[i].UpdatedAt = &updatedAt
		}
	}

	middleware.SuccessResponse(c, http.StatusOK, gin.H{
		"comments": commentResponses,
		"total":    total,
		"offset":   offset,
		"limit":    limit,
	})
}

// UpdateComment handles PUT /api/v1/comments/:commentId - update comment
// @Summary Update a comment
// @Description Update a comment (author or admin only). Can only update within 30 days of creation.
// @Tags comments
// @Security Bearer
// @Accept json
// @Produce json
// @Param commentId path string true "Comment ID"
// @Param request body UpdateCommentRequest true "Comment update request"
// @Success 200 {object} map[string]string "Comment updated successfully"
// @Failure 400 {object} middleware.SwaggerErrorResponse "Validation error"
// @Failure 401 {object} middleware.SwaggerErrorResponse "Not authenticated"
// @Failure 403 {object} middleware.SwaggerErrorResponse "Permission denied or comment too old"
// @Failure 404 {object} middleware.SwaggerErrorResponse "Comment not found"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /comments/{commentId} [put]
func (cc *CommentController) UpdateComment(c *gin.Context) {
	commentID := c.Param("commentId")
	var req UpdateCommentRequest

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

	userRole, _ := utilities.GetUserRoleFromContext(c)

	// Get comment to verify ownership
	comment, err := cc.commentService.GetCommentByID(c.Request.Context(), commentID)
	if err != nil || comment == nil {
		middleware.NotFoundErrorResponse(c, "Comment not found")
		return
	}

	// Convert userRole string to models.UserRole
	role := models.UserRole(userRole)

	// Update comment via service (authorization handled in service)
	if err := cc.commentService.UpdateComment(c.Request.Context(), commentID, userID, role, strings.TrimSpace(req.Text)); err != nil {
		if err == models.ErrUnauthorized {
			middleware.ForbiddenErrorResponse(c, "You do not have permission to update this comment")
			return
		}
		middleware.InternalErrorResponse(c, "Failed to update comment")
		return
	}

	middleware.SuccessResponse(c, http.StatusOK, gin.H{
		"message": "Comment updated successfully",
	})
}

// DeleteComment handles DELETE /api/v1/comments/:commentId - delete comment
// @Summary Delete a comment
// @Description Delete a comment (author or admin only). Soft-deleted, admin can permanently delete.
// @Tags comments
// @Security Bearer
// @Produce json
// @Param commentId path string true "Comment ID"
// @Success 200 {object} map[string]string "Comment deleted successfully"
// @Failure 401 {object} middleware.SwaggerErrorResponse "Not authenticated"
// @Failure 403 {object} middleware.SwaggerErrorResponse "Permission denied"
// @Failure 404 {object} middleware.SwaggerErrorResponse "Comment not found"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /comments/{commentId} [delete]
func (cc *CommentController) DeleteComment(c *gin.Context) {
	commentID := c.Param("commentId")

	// Get current user
	userID, ok := utilities.GetUserIDFromContext(c)
	if !ok {
		middleware.AuthErrorResponse(c, "User not authenticated")
		return
	}

	userRole, _ := utilities.GetUserRoleFromContext(c)

	// Convert userRole string to models.UserRole
	role := models.UserRole(userRole)

	// Delete comment via service (authorization handled in service)
	if err := cc.commentService.DeleteComment(c.Request.Context(), commentID, userID, role); err != nil {
		if err == models.ErrUnauthorized {
			middleware.ForbiddenErrorResponse(c, "You do not have permission to delete this comment")
			return
		}
		middleware.InternalErrorResponse(c, "Failed to delete comment")
		return
	}

	middleware.SuccessResponse(c, http.StatusOK, gin.H{
		"message": "Comment deleted successfully",
	})
}
