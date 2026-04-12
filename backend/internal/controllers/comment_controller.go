package controllers

import (
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
	planService services.PlanService
	// CommentService will be added when implemented
	// commentService services.CommentService
}

// NewCommentController creates a new comment controller
func NewCommentController(planService services.PlanService) *CommentController {
	return &CommentController{
		planService: planService,
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
	ID               string `json:"id"`
	PlanID           string `json:"plan_id"`
	AuthorID         string `json:"author_id"`
	Text             string `json:"text"`
	IsDeletedByAdmin bool   `json:"is_deleted_by_admin"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        *string `json:"updated_at,omitempty"`
}

// CreateComment handles POST /api/v1/plans/:id/comments - create comment
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

	// TODO: When CommentService is implemented, call it here
	// commentID, err := cc.commentService.CreateComment(c.Request.Context(), comment)

	middleware.SuccessResponse(c, http.StatusCreated, gin.H{
		"message": "Comment created successfully (service not yet implemented)",
	})
}

// GetComments handles GET /api/v1/plans/:id/comments - list comments
func (cc *CommentController) GetComments(c *gin.Context) {
	planID := c.Param("id")

	// Verify plan exists
	plan, _ := cc.planService.GetPlanByID(c.Request.Context(), planID)
	if plan == nil {
		middleware.NotFoundErrorResponse(c, "Plan not found")
		return
	}

	// TODO: When CommentService is implemented, fetch comments here
	// comments, err := cc.commentService.GetCommentsByPlan(c.Request.Context(), planID)

	middleware.SuccessResponse(c, http.StatusOK, gin.H{
		"comments": []CommentResponse{},
		"message": "Comments feature not yet implemented",
	})
}

// UpdateComment handles PUT /api/v1/comments/:commentId - update comment
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

	// TODO: When CommentService is implemented
	// Verify ownership
	// comment, _ := cc.commentService.GetCommentByID(c.Request.Context(), commentID)
	// if comment.AuthorID != userID {
	//     middleware.ForbiddenErrorResponse(c, "You do not have permission to update this comment")
	//     return
	// }
	// err := cc.commentService.UpdateComment(c.Request.Context(), commentID, req.Text)

	_ = userID // Suppress unused warning
	_ = commentID

	middleware.SuccessResponse(c, http.StatusOK, gin.H{
		"message": "Update comment feature not yet implemented",
	})
}

// DeleteComment handles DELETE /api/v1/comments/:commentId - delete comment
func (cc *CommentController) DeleteComment(c *gin.Context) {
	commentID := c.Param("commentId")

	// Get current user
	userID, ok := utilities.GetUserIDFromContext(c)
	if !ok {
		middleware.AuthErrorResponse(c, "User not authenticated")
		return
	}

	// TODO: When CommentService is implemented
	// Verify ownership or admin
	// comment, _ := cc.commentService.GetCommentByID(c.Request.Context(), commentID)
	// if comment.AuthorID != userID {
	//     middleware.ForbiddenErrorResponse(c, "You do not have permission to delete this comment")
	//     return
	// }
	// err := cc.commentService.DeleteComment(c.Request.Context(), commentID)

	_ = userID // Suppress unused warning
	_ = commentID

	middleware.SuccessResponse(c, http.StatusOK, gin.H{
		"message": "Delete comment feature not yet implemented",
	})
}
