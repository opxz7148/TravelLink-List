package utilities

import (
	"fmt"
	"tll-backend/internal/middleware"
	"tll-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// GetUserIDFromContext retrieves the authenticated user's ID from the request context
// Returns (userID, ok) where ok is false if user is not authenticated
func GetUserIDFromContext(c *gin.Context) (string, bool) {
	claims := middleware.GetUserClaims(c)
	fmt.Println(claims)
	if claims == nil {
		return "", false
	}
	return claims.UserID, true
}

// GetUserEmailFromContext retrieves the authenticated user's email from the request context
// Returns (email, ok) where ok is false if user is not authenticated
func GetUserEmailFromContext(c *gin.Context) (string, bool) {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		return "", false
	}
	return claims.Email, true
}

// GetUserRoleFromContext retrieves the authenticated user's role from the request context
// Returns (role, ok) where ok is false if user is not authenticated
func GetUserRoleFromContext(c *gin.Context) (string, bool) {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		return "", false
	}
	return claims.Role, true
}

// GetUserClaimsFromContext retrieves the full TokenClaims from the request context
// Returns (claims, ok) where ok is false if user is not authenticated
func GetUserClaimsFromContext(c *gin.Context) (*services.TokenClaims, bool) {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		return nil, false
	}
	return claims, true
}

// IsUserAdmin checks if the authenticated user has admin role
func IsUserAdmin(c *gin.Context) bool {
	role, ok := GetUserRoleFromContext(c)
	return ok && role == "admin"
}

// IsUserTraveller checks if the authenticated user has traveller role
func IsUserTraveller(c *gin.Context) bool {
	role, ok := GetUserRoleFromContext(c)
	return ok && role == "traveller"
}

// IsUserSimple checks if the authenticated user has simple role
func IsUserSimple(c *gin.Context) bool {
	role, ok := GetUserRoleFromContext(c)
	return ok && role == "simple"
}
