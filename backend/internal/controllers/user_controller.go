package controllers

import (
	"net/http"
	"strings"

	"tll-backend/internal/logger"
	"tll-backend/internal/middleware"
	"tll-backend/internal/models"
	"tll-backend/internal/services"
	"tll-backend/internal/utilities"

	"github.com/gin-gonic/gin"
)

// UserController handles user profile operations
type UserController struct {
	userService services.UserService
}

// NewUserController creates a new user controller
func NewUserController(userService services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// UpdateProfileRequest represents profile update request
type UpdateProfileRequest struct {
	DisplayName string `json:"display_name" binding:"max=100"`
	Bio         string `json:"bio" binding:"max=500"`
}

// ChangePasswordRequest represents password change request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=1"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// UserProfileResponse represents user profile in API response
type UserProfileResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	DisplayName string `json:"display_name"`
	Bio       string `json:"bio"`
	Role      string `json:"role"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// GetProfile handles GET /api/v1/users/:id - get user profile
// @Summary Get user profile by ID
// @Description Retrieve user profile information (public endpoint)
// @Tags users
// @Produce json
// @Param id path string true "User ID or 'me' for current user"
// @Success 200 {object} map[string]UserProfileResponse "User profile retrieved successfully"
// @Failure 404 {object} middleware.SwaggerErrorResponse "User not found"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /users/{id} [get]
func (uc *UserController) GetProfile(c *gin.Context) {
	userID := c.Param("id")

	// Handle special case: "me" refers to current authenticated user
	if userID == "me" {
		claims := middleware.GetUserClaims(c)
		if claims == nil || claims.UserID == "" {
			middleware.AuthErrorResponse(c, "Authentication required to access 'me' endpoint")
			return
		}
		userID = claims.UserID
	}

	// Get user via service
	user, err := uc.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		middleware.InternalErrorResponse(c, "Failed to fetch user profile")
		return
	}

	if user == nil {
		middleware.NotFoundErrorResponse(c, "User not found")
		return
	}

	// Convert to response format
	profileResp := UserProfileResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		DisplayName: user.DisplayName,
		Bio:       user.Bio,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	middleware.SuccessResponse(c, http.StatusOK, gin.H{
		"user": profileResp,
	})
}

// UpdateProfile handles PUT /api/v1/users/:id - update user profile
// @Summary Update user profile
// @Description Update user profile information (display_name, bio). Only user or admin can update.
// @Tags users
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body UpdateProfileRequest true "Profile update request"
// @Success 200 {object} map[string]interface{} "Profile updated with user and message"
// @Failure 400 {object} middleware.SwaggerErrorResponse "Validation error"
// @Failure 401 {object} middleware.SwaggerErrorResponse "User not authenticated"
// @Failure 403 {object} middleware.SwaggerErrorResponse "Permission denied"
// @Failure 404 {object} middleware.SwaggerErrorResponse "User not found"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /users/{id} [put]
func (uc *UserController) UpdateProfile(c *gin.Context) {
	userID := c.Param("id")
	var req UpdateProfileRequest
	log := logger.GetLogger("UserController")

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		log.ValidationError(c.Request.Context(), map[string]string{"request": "Invalid request format"})
		middleware.ValidationErrorResponse(c, "invalid request", nil)
		return
	}

	// Get authenticated user
	authUserID, ok := utilities.GetUserIDFromContext(c)
	if !ok {
		log.AuthorizationError(c.Request.Context(), "", "Attempted profile update without authentication")
		middleware.AuthErrorResponse(c, "User not authenticated")
		return
	}

	// Only allow users to update their own profile, unless admin
	userRole, _ := utilities.GetUserRoleFromContext(c)
	isAdmin := userRole == "admin"
	if authUserID != userID && !isAdmin {
		log.AuthorizationError(c.Request.Context(), authUserID, "Attempted to update another user's profile without admin role")
		middleware.ForbiddenErrorResponse(c, "You can only update your own profile")
		return
	}

	// Trim whitespace
	displayName := strings.TrimSpace(req.DisplayName)
	bio := strings.TrimSpace(req.Bio)

	// Update profile via service
	user, err := uc.userService.UpdateProfile(c.Request.Context(), userID, displayName, bio, "")
	if err != nil {
		if err == models.ErrValidation {
			log.ValidationError(c.Request.Context(), map[string]string{"profile": "Invalid profile data"})
			middleware.ValidationErrorResponse(c, "Invalid profile data", nil)
			return
		}
		if err == models.ErrNotFound {
			log.Info(c.Request.Context(), "Profile update attempted on non-existent user", logger.LogAttributes{UserID: userID})
			middleware.NotFoundErrorResponse(c, "User not found")
			return
		}
		log.ServiceError(c.Request.Context(), "UserService", "UpdateProfile", err, authUserID)
		middleware.InternalErrorResponse(c, "Failed to update profile")
		return
	}

	// Convert to response format
	profileResp := UserProfileResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		DisplayName: user.DisplayName,
		Bio:       user.Bio,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	log.Info(c.Request.Context(), "Profile updated successfully", logger.LogAttributes{UserID: authUserID, Details: "Updated display_name and bio"})

	middleware.SuccessResponse(c, http.StatusOK, gin.H{
		"user":    profileResp,
		"message": "Profile updated successfully",
	})
}

// ChangePassword handles POST /api/v1/users/:id/change-password - change password
// @Summary Change user password
// @Description Change password for authenticated user. Users can only change their own password.
// @Tags users
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body ChangePasswordRequest true "Password change request"
// @Success 200 {object} map[string]string "Password changed successfully"
// @Failure 400 {object} middleware.SwaggerErrorResponse "Validation error (invalid old password or weak new password)"
// @Failure 401 {object} middleware.SwaggerErrorResponse "User not authenticated"
// @Failure 403 {object} middleware.SwaggerErrorResponse "Can only change own password"
// @Failure 404 {object} middleware.SwaggerErrorResponse "User not found"
// @Failure 500 {object} middleware.SwaggerErrorResponse "Internal server error"
// @Router /users/{id}/change-password [post]
func (uc *UserController) ChangePassword(c *gin.Context) {
	userID := c.Param("id")
	var req ChangePasswordRequest
	log := logger.GetLogger("UserController")

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		log.ValidationError(c.Request.Context(), map[string]string{"request": "Invalid request format"})
		middleware.ValidationErrorResponse(c, "invalid request", nil)
		return
	}

	// Get authenticated user
	authUserID, ok := utilities.GetUserIDFromContext(c)
	if !ok {
		log.AuthorizationError(c.Request.Context(), "", "Attempted password change without authentication")
		middleware.AuthErrorResponse(c, "User not authenticated")
		return
	}

	// Only allow users to change their own password
	if authUserID != userID {
		log.AuthorizationError(c.Request.Context(), authUserID, "Attempted to change another user's password")
		middleware.ForbiddenErrorResponse(c, "You can only change your own password")
		return
	}

	// Change password via service
	err := uc.userService.ChangePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword)
	if err != nil {
		if err == models.ErrInvalidCredentials {
			log.AuthenticationError(c.Request.Context(), "", "Invalid current password provided")
			middleware.ValidationErrorResponse(c, "Current password is incorrect", gin.H{"field": "old_password"})
			return
		}
		if err == models.ErrValidation {
			log.ValidationError(c.Request.Context(), map[string]string{"new_password": "Does not meet security requirements"})
			middleware.ValidationErrorResponse(c, "New password does not meet security requirements", gin.H{"field": "new_password"})
			return
		}
		if err == models.ErrNotFound {
			log.Info(c.Request.Context(), "Password change attempted on non-existent user", logger.LogAttributes{UserID: userID})
			middleware.NotFoundErrorResponse(c, "User not found")
			return
		}
		log.ServiceError(c.Request.Context(), "UserService", "ChangePassword", err, authUserID)
		middleware.InternalErrorResponse(c, "Failed to change password")
		return
	}

	log.SecurityEvent(c.Request.Context(), "PASSWORD_CHANGED", authUserID, "", "User changed their password")

	middleware.SuccessResponse(c, http.StatusOK, gin.H{
		"message": "Password changed successfully",
	})
}
