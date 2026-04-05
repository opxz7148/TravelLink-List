package middleware

import (
	"github.com/gin-gonic/gin"
)

// ContextHelper provides convenient methods to access data stored in request context
type ContextHelper struct{}

// NewContextHelper creates a new context helper
func NewContextHelper() *ContextHelper {
	return &ContextHelper{}
}

// GetUserID retrieves the authenticated user's ID from the request context
func (h *ContextHelper) GetUserID(c *gin.Context) string {
	claims := GetUserClaims(c)
	if claims == nil {
		return ""
	}
	return claims.UserID
}

// GetUserEmail retrieves the authenticated user's email from the request context
func (h *ContextHelper) GetUserEmail(c *gin.Context) string {
	claims := GetUserClaims(c)
	if claims == nil {
		return ""
	}
	return claims.Email
}

// GetUserUsername retrieves the authenticated user's username from the request context
func (h *ContextHelper) GetUserUsername(c *gin.Context) string {
	claims := GetUserClaims(c)
	if claims == nil {
		return ""
	}
	return claims.Username
}

// GetUserRole retrieves the authenticated user's role from the request context
func (h *ContextHelper) GetUserRole(c *gin.Context) string {
	claims := GetUserClaims(c)
	if claims == nil {
		return ""
	}
	return claims.Role
}

// IsAdmin checks if the authenticated user is an admin
func (h *ContextHelper) IsAdmin(c *gin.Context) bool {
	claims := GetUserClaims(c)
	if claims == nil {
		return false
	}
	return claims.Role == "admin"
}

// IsTraveller checks if the authenticated user is a traveller
func (h *ContextHelper) IsTraveller(c *gin.Context) bool {
	claims := GetUserClaims(c)
	if claims == nil {
		return false
	}
	return claims.Role == "traveller"
}

// IsSimple checks if the authenticated user has simple role
func (h *ContextHelper) IsSimple(c *gin.Context) bool {
	claims := GetUserClaims(c)
	if claims == nil {
		return false
	}
	return claims.Role == "simple"
}

// HasRole checks if the authenticated user has one of the specified roles
func (h *ContextHelper) HasRole(c *gin.Context, roles ...string) bool {
	claims := GetUserClaims(c)
	if claims == nil {
		return false
	}
	for _, role := range roles {
		if claims.Role == role {
			return true
		}
	}
	return false
}
