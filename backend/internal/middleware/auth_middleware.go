package middleware

import (
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

// ExtractClaims is middleware that extracts and validates user token if present
// If Authorization header exists and token is valid, stores TokenClaims in context
// If missing or validation fails, stores nil in context (no error response)
// This allows handlers to check if user is authenticated without requiring it
func ExtractClaims(jwtService services.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// No auth header - store nil
			c.Set(UserContextKey, nil)
			c.Next()
			return
		}

		// Extract token from "Bearer <token>" format
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != BearerScheme {
			// Invalid format - store nil
			c.Set(UserContextKey, nil)
			c.Next()
			return
		}

		token := parts[1]

		// Validate token
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			// Invalid token - store nil
			c.Set(UserContextKey, nil)
			c.Next()
			return
		}

		// Store validated claims in context
		c.Set(UserContextKey, claims)
		c.Next()
	}
}

// RequireAuth is middleware that enforces authentication
// It checks that validated claims exist in the request context (from ExtractClaims)
// If claims are nil (extraction/validation failed), returns 401 Unauthorized
// Must be used after ExtractClaims middleware in the middleware chain
func RequireAuth(jwtService services.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if claims were extracted and validated by ExtractClaims middleware
		claims := GetUserClaims(c)
		if claims == nil {
			AuthErrorResponse(c, "missing or invalid authorization")
			c.Abort()
			return
		}
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

// RequireRole is middleware that verifies the user has one of the required roles from the database
// This checks the CURRENT role from the database, not the cached JWT role
// Must be used after RequireAuth middleware
// If user doesn't have required role, returns 403 Forbidden
func RequireRole(userService services.UserService, allowedRoles ...models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := GetUserClaims(c)
		if claims == nil {
			ForbiddenErrorResponse(c, "user information not found")
			c.Abort()
			return
		}

		// Fetch current user from database to get latest role
		user, err := userService.GetUserByID(c.Request.Context(), claims.UserID)
		if err != nil || user == nil {
			ForbiddenErrorResponse(c, "user not found")
			c.Abort()
			return
		}

		// Check if user's current role is in allowed roles
		userRole := models.UserRole(user.Role)
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

// RequireAdmin is a middleware that ensures the user is an admin (checked from database)
// This verifies the CURRENT admin status, not the cached JWT role
// Must be used after RequireAuth middleware
func RequireAdmin(userService services.UserService) gin.HandlerFunc {
	return RequireRole(userService, models.RoleAdmin)
}

// RequireTraveller is a middleware that ensures the user is at least a traveller (checked from database)
// This verifies the CURRENT traveller status, not the cached JWT role
// Must be used after RequireAuth middleware
func RequireTraveller(userService services.UserService) gin.HandlerFunc {
	return RequireRole(userService, models.RoleTraveller)
}

// RequireTravellerOrAdmin is a middleware that ensures the user is a traveller or admin (checked from database)
// This verifies the CURRENT role, not the cached JWT role
// Must be used after RequireAuth middleware
func RequireTravellerOrAdmin(userService services.UserService) gin.HandlerFunc {
	return RequireRole(userService, models.RoleTraveller, models.RoleAdmin)
}

// RequireNonAdmin is a middleware that ensures the user is NOT an admin (checked from database)
// Used for operations like plan creation that are restricted to non-admins
// This verifies the CURRENT admin status, not the cached JWT role
// Must be used after RequireAuth middleware
func RequireNonAdmin(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := GetUserClaims(c)
		if claims == nil {
			ForbiddenErrorResponse(c, "user information not found")
			c.Abort()
			return
		}

		// Fetch current user from database to get latest role
		user, err := userService.GetUserByID(c.Request.Context(), claims.UserID)
		if err != nil || user == nil {
			ForbiddenErrorResponse(c, "user not found")
			c.Abort()
			return
		}

		// Check if user is admin
		if models.UserRole(user.Role) == models.RoleAdmin {
			ForbiddenErrorResponse(c, "admins cannot perform this action")
			c.Abort()
			return
		}

		c.Next()
	}
}
