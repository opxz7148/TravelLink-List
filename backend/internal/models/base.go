package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Standard error definitions for repository operations
var (
	// ErrNotFound is returned when a requested record doesn't exist
	ErrNotFound = errors.New("record not found")

	// ErrDuplicateEmail is returned when trying to create a user with an email that already exists
	ErrDuplicateEmail = errors.New("email already exists")

	// ErrDuplicateUsername is returned when trying to create a user with a username that already exists
	ErrDuplicateUsername = errors.New("username already exists")

	// ErrInvalidRole is returned when an invalid role is used
	ErrInvalidRole = errors.New("invalid role")

	// ErrUnauthorized is returned when a user doesn't have permission for an operation
	ErrUnauthorized = errors.New("unauthorized")

	// ErrValidation is returned when a model fails validation
	ErrValidation = errors.New("validation failed")
)

// BaseModel contains common fields for all models
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
