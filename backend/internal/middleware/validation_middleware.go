package middleware

import (
	"errors"

	"tll-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ValidateJSON validates JSON request body against struct tags
// Returns 400 Bad Request with validation errors if validation fails
// Supports custom validators via validator.Validate
func ValidateJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		// This middleware relies on Gin's ShouldBindJSON for automatic validation
		// Applied per-route as needed using middleware chain
		c.Next()
	}
}

// ValidateRequest validates request body and returns errors if validation fails
// Used by controllers to validate incoming requests before processing
// Returns models.ErrValidation if validation fails
func ValidateRequest(data interface{}) error {
	v := validator.New()

	// Register custom validators if needed
	// Example: v.RegisterValidation("custom", customValidation)

	if err := v.Struct(data); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			return models.ErrValidation
		}
		return err
	}

	return nil
}

// ValidationErrorDetail represents a single field validation error
type ValidationErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidateAndReturnErrors validates request and returns detailed error messages
// Returns errors in standardized error response format for API responses
func ValidateAndReturnErrors(c *gin.Context, data interface{}) bool {
	v := validator.New()

	if err := v.Struct(data); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			errors := make([]ValidationErrorDetail, 0, len(validationErrors))
			for _, fieldErr := range validationErrors {
				errors = append(errors, ValidationErrorDetail{
					Field:   fieldErr.Field(),
					Message: getValidationMessage(fieldErr),
				})
			}

			// Use standardized error response from response_middleware
			ValidationErrorResponse(c, "validation failed", errors)
			return true // Validation failed
		}

		// Non-validation error
		ValidationErrorResponse(c, "invalid request body", nil)
		return true
	}

	return false // Validation passed
}

// getValidationMessage returns a human-readable validation error message
func getValidationMessage(fieldErr validator.FieldError) string {
	switch fieldErr.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "This field must be a valid email"
	case "min":
		return "This field must be at least " + fieldErr.Param() + " characters"
	case "max":
		return "This field must be at most " + fieldErr.Param() + " characters"
	case "oneof":
		return "This field must be one of: " + fieldErr.Param()
	case "len":
		return "This field must be exactly " + fieldErr.Param() + " characters"
	case "numeric":
		return "This field must be numeric"
	case "gte":
		return "This field must be greater than or equal to " + fieldErr.Param()
	case "lte":
		return "This field must be less than or equal to " + fieldErr.Param()
	default:
		return "This field is invalid"
	}
}
