package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ResponseEnvelope is the standard response wrapper for all API responses
type ResponseEnvelope struct {
	Success    bool        `json:"success"`
	APIVersion string      `json:"api_version"`
	Data       interface{} `json:"data"`
	Error      *ErrorInfo  `json:"error"`
	Timestamp  string      `json:"timestamp"`
}

// ErrorInfo contains error details in the response
type ErrorInfo struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// SwaggerErrorResponse is a response structure for Swagger documentation purposes
// This type is used to document error responses in Swagger/OpenAPI specs
type SwaggerErrorResponse struct {
	Success bool      `json:"success" example:"false"`
	Error   ErrorInfo `json:"error"`
}

// SuccessResponse returns a standard success response envelope
func SuccessResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, ResponseEnvelope{
		Success:    true,
		APIVersion: "1.0",
		Data:       data,
		Error:      nil,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	})
}

// ErrorResponse returns a standard error response envelope
func ErrorResponse(c *gin.Context, statusCode int, code, message string, details interface{}) {
	c.JSON(statusCode, ResponseEnvelope{
		Success:    false,
		APIVersion: "1.0",
		Data:       nil,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// AuthErrorResponse returns a standard authentication error response
func AuthErrorResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnauthorized, "AUTH_ERROR", message, nil)
}

// AuthErrorResponseWithCode returns an authentication error with specific error code
func AuthErrorResponseWithCode(c *gin.Context, code, message string) {
	ErrorResponse(c, http.StatusUnauthorized, code, message, nil)
}

// ValidationErrorResponse returns a standard validation error response
func ValidationErrorResponse(c *gin.Context, message string, details interface{}) {
	ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", message, details)
}

// ForbiddenErrorResponse returns a standard forbidden error response
func ForbiddenErrorResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusForbidden, "FORBIDDEN_ERROR", message, nil)
}

// NotFoundErrorResponse returns a standard not found error response
func NotFoundErrorResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", message, nil)
}

// InternalErrorResponse returns a standard internal server error response
func InternalErrorResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", message, nil)
}
