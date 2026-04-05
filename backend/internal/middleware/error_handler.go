package middleware

import (
	"tll-backend/internal/models"

	"github.com/gin-gonic/gin"
)

// HandleServiceError converts service layer errors to appropriate HTTP responses
func HandleServiceError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	switch err {
	case models.ErrNotFound:
		NotFoundErrorResponse(c, "Resource not found")
	case models.ErrDuplicateEmail:
		ValidationErrorResponse(c, "Email already registered", map[string]string{
			"field":  "email",
			"reason": "Email must be unique",
		})
	case models.ErrDuplicateUsername:
		ValidationErrorResponse(c, "Username already taken", map[string]string{
			"field":  "username",
			"reason": "Username must be unique",
		})
	case models.ErrValidation:
		ValidationErrorResponse(c, "Validation failed", nil)
	case models.ErrUnauthorized:
		AuthErrorResponse(c, "Invalid credentials")
	case models.ErrInvalidRole:
		ForbiddenErrorResponse(c, "Invalid user role")
	default:
		InternalErrorResponse(c, "An unexpected error occurred")
	}

	return true
}

// UserAuthData represents user authentication data from JWT claims
type UserAuthData struct {
	UserID   string
	Email    string
	Username string
	Role     string
}

// ExtractUserAuthData extracts user authentication data from context
func ExtractUserAuthData(c *gin.Context) *UserAuthData {
	claims := GetUserClaims(c)
	if claims == nil {
		return nil
	}

	return &UserAuthData{
		UserID:   claims.UserID,
		Email:    claims.Email,
		Username: claims.Username,
		Role:     claims.Role,
	}
}

// EnsureUserOwnsResource checks if the authenticated user is the resource owner
func EnsureUserOwnsResource(c *gin.Context, resourceUserID string) bool {
	authData := ExtractUserAuthData(c)
	if authData == nil {
		ForbiddenErrorResponse(c, "User not authenticated")
		return false
	}

	if authData.UserID != resourceUserID {
		ForbiddenErrorResponse(c, "You do not have permission to modify this resource")
		return false
	}

	return true
}

// ErrorWithStatus is a helper to attach HTTP status to error handling
type ErrorWithStatus struct {
	StatusCode int
	Code       string
	Message    string
	Details    interface{}
}

// HandleErrorWithStatus handles errors with explicit status codes
func HandleErrorWithStatus(c *gin.Context, errStatus ErrorWithStatus) {
	ErrorResponse(c, errStatus.StatusCode, errStatus.Code, errStatus.Message, errStatus.Details)
}
