package middleware

import (
	"fmt"
	"net/http"
	"time"

	"tll-backend/internal/logger"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware logs all HTTP requests and responses
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.GetLogger("HTTP")

		// Extract user information if available
	userID := ""
	userEmail := ""
	if claims := GetUserClaims(c); claims != nil {
		userID = claims.UserID
		}

		// Record request start time
		startTime := time.Now()

		// Log incoming request
		log.Info(c.Request.Context(), fmt.Sprintf("Incoming %s request", c.Request.Method),
			logger.LogAttributes{
				UserID:   userID,
				Email:    userEmail,
				Endpoint: c.Request.URL.Path,
				Method:   c.Request.Method,
			})

		// Process request
		c.Next()

		// Calculate request duration
		duration := time.Since(startTime)

		// Log response
		statusCode := c.Writer.Status()
		if statusCode >= 400 {
			// Log errors and denied access
			reason := ""
			switch statusCode {
			case http.StatusUnauthorized:
				reason = "Unauthorized"
			case http.StatusForbidden:
				reason = "Forbidden"
			case http.StatusNotFound:
				reason = "Not Found"
			case http.StatusBadRequest:
				reason = "Bad Request"
			case http.StatusConflict:
				reason = "Conflict"
			case http.StatusInternalServerError:
				reason = "Internal Server Error"
			default:
				reason = http.StatusText(statusCode)
			}

			if statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden {
				log.AuthorizationDenied(c.Request.Context(), userID, userEmail, c.Request.URL.Path, c.Request.Method, reason)
			} else {
				log.RequestLogging(c.Request.Context(), c.Request.Method, c.Request.URL.Path, userID, statusCode)
			}
		} else {
			log.RequestLogging(c.Request.Context(), c.Request.Method, c.Request.URL.Path, userID, statusCode)
		}

		// Log slow requests (over 500ms)
		if duration > 500*time.Millisecond {
			log.Warn(c.Request.Context(), fmt.Sprintf("Slow Request: %s %s took %dms", c.Request.Method, c.Request.URL.Path, duration.Milliseconds()),
				logger.LogAttributes{
					UserID:     userID,
					Email:      userEmail,
					Endpoint:   c.Request.URL.Path,
					Method:     c.Request.Method,
					Status:     statusCode,
					Details:    fmt.Sprintf("Duration: %d ms", duration.Milliseconds()),
				})
		}
	}
}

// ErrorHandlingMiddleware catches panics and server errors
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log := logger.GetLogger("PanicRecovery")

				userID := ""
				userEmail := ""
				if claims := GetUserClaims(c); claims != nil {
					userID = claims.UserID
					userEmail = claims.Email
				}

				// Log panic as error
				log.Error(c.Request.Context(), fmt.Sprintf("Panic recovered: %v", err), nil,
					logger.LogAttributes{
						UserID:   userID,
						Email:    userEmail,
						Endpoint: c.Request.URL.Path,
						Method:   c.Request.Method,
						Details:  fmt.Sprintf("%v", err),
					})

				// Return 500 error
				c.JSON(500, gin.H{"error": "Internal server error"})
			}
		}()

		c.Next()
	}
}
