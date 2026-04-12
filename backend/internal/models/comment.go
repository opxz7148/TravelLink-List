package models

import (
	"strings"
	"time"
)

// Comment represents user commentary on a TravelPlan.
// Follows the specification from data-model.md with all validation rules.
//
// Comments enable community engagement on travel plans.
// Each comment is authored by a user and associated with a travel plan.
// Multiple comments per user per plan are allowed, but comments can be edited post-creation.
type Comment struct {
	// id (UUID, primary key)
	// Unique identifier for the comment, typically generated as UUID v4
	ID string `gorm:"primaryKey;type:TEXT" json:"id"`

	// plan_id (UUID, foreign key to TravelPlan)
	// The travel plan this comment belongs to
	PlanID string `gorm:"type:TEXT;not null;index" json:"plan_id"`

	// author_id (UUID, foreign key to User)
	// The user who authored this comment
	AuthorID string `gorm:"type:TEXT;not null;index" json:"author_id"`

	// text (string, max 1000 chars, required, min 1 char)
	// The comment content
	// Must be trimmed and non-empty
	Text string `gorm:"type:TEXT;not null;size:1000" json:"text"`

	// is_deleted_by_admin (boolean, default: false, soft-delete)
	// Soft-delete flag; true means comment is hidden but preserved
	IsDeletedByAdmin bool `gorm:"default:false;not null" json:"is_deleted_by_admin"`

	// created_at (timestamp, UTC, immutable)
	// Timestamp when comment was created
	CreatedAt time.Time `gorm:"autoCreateTime:milli;type:TIMESTAMP;not null" json:"created_at"`

	// updated_at (timestamp, UTC, nullable, only for non-deleted comments)
	// Timestamp when comment was last updated
	// Null if not yet updated after creation
	UpdatedAt *time.Time `gorm:"autoUpdateTime:milli;type:TIMESTAMP" json:"updated_at"`
}

// TableName specifies the table name for the Comment model
func (Comment) TableName() string {
	return "comments"
}

// Validate validates comment fields according to data-model.md rules
// Returns ErrValidation if validation fails
func (c *Comment) Validate() error {
	// Text validation: Required, 1-1000 chars, trimmed
	text := strings.TrimSpace(c.Text)
	if text == "" {
		return ErrValidation // Empty or whitespace-only text
	}
	if len(text) > 1000 {
		return ErrValidation // Text exceeds max length
	}

	// Plan must exist (not empty)
	if strings.TrimSpace(c.PlanID) == "" {
		return ErrValidation // Plan ID required
	}

	// Author must exist (not empty)
	if strings.TrimSpace(c.AuthorID) == "" {
		return ErrValidation // Author ID required
	}

	return nil
}

// CanBeEditedByUser checks if a user can edit this comment
// Only author or admin can edit
func (c *Comment) CanBeEditedByUser(userID string, userRole UserRole) bool {
	// Admin can always edit
	if userRole == RoleAdmin {
		return true
	}

	// Author can edit if not soft-deleted
	if c.AuthorID == userID && !c.IsDeletedByAdmin {
		return true
	}

	return false
}

// CanBeDeletedByUser checks if a user can delete (soft-delete) this comment
// Only author or admin can delete
func (c *Comment) CanBeDeletedByUser(userID string, userRole UserRole) bool {
	// Admin can always delete
	if userRole == RoleAdmin {
		return true
	}

	// Author can delete if not already deleted
	if c.AuthorID == userID && !c.IsDeletedByAdmin {
		return true
	}

	return false
}

// CanBeViewedByUser checks if a user can view this comment
// Deleted comments visible only to admin and author
func (c *Comment) CanBeViewedByUser(userID string, userRole UserRole) bool {
	// If not deleted, everyone can view
	if !c.IsDeletedByAdmin {
		return true
	}

	// If deleted, only author and admin can view
	if userRole == RoleAdmin || c.AuthorID == userID {
		return true
	}

	return false
}
