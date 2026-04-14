package middleware

import (
	"fmt"
	"strings"

	"tll-backend/internal/models"
	"tll-backend/internal/services"

	"github.com/gin-gonic/gin"
)

const (
	// UserContextKey is the key used to store user claims in the request context
	UserContextKey = "user_claims"
	// BearerScheme is the authentication scheme for JWT tokens
	BearerScheme = "Bearer"
)

// RequireAuth is middleware that validates JWT tokens and extracts user claims
// It expects an "Authorization: Bearer <token>" header
// If valid, it stores the TokenClaims in the request context
// If invalid or missing, it returns 401 Unauthorized
func RequireAuth(jwtService services.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			AuthErrorResponse(c, "missing authorization header")
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>" format
		parts := strings.SplitN(authHeader, " ", 2)
		fmt.Println(authHeader)
		if len(parts) != 2 || parts[0] != BearerScheme {
			AuthErrorResponse(c, "invalid authorization header format")
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			AuthErrorResponse(c, "invalid or expired token")
			c.Abort()
			return
		}

		// Store claims in context for later use
		c.Set(UserContextKey, claims)
		c.Next()
	}
}

// GetUserClaims retrieves the TokenClaims from the request context
// Should be called after RequireAuth middleware
// Returns nil if claims are not found in context
func GetUserClaims(c *gin.Context) *services.TokenClaims {
	claims, exists := c.Get(UserContextKey)
	if !exists {
		return nil
	}
	tokenClaims, ok := claims.(*services.TokenClaims)
	if !ok {
		return nil
	}
	return tokenClaims
}

// RequireRole is middleware that verifies the user has one of the required roles
// Must be used after RequireAuth middleware
// If user doesn't have required role, returns 403 Forbidden
func RequireRole(allowedRoles ...models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := GetUserClaims(c)
		if claims == nil {
			ForbiddenErrorResponse(c, "user information not found")
			c.Abort()
			return
		}

		// Check if user's role is in allowed roles
		userRole := models.UserRole(claims.Role)
		for _, allowed := range allowedRoles {
			if userRole == allowed {
				c.Next()
				return
			}
		}

		// User does not have required role
		ForbiddenErrorResponse(c, "insufficient permissions for this action")
		c.Abort()
	}
}

// RequireAdmin is a middleware that ensures the user is an admin
// Must be used after RequireAuth middleware
func RequireAdmin() gin.HandlerFunc {
	return RequireRole(models.RoleAdmin)
}

// RequireTraveller is a middleware that ensures the user is at least a traveller
// Must be used after RequireAuth middleware
func RequireTraveller() gin.HandlerFunc {
	return RequireRole(models.RoleTraveller)
}

// RequireTravellerOrAdmin is a middleware that ensures the user is a traveller or admin
// Must be used after RequireAuth middleware
func RequireTravellerOrAdmin() gin.HandlerFunc {
	return RequireRole(models.RoleTraveller, models.RoleAdmin)
}
