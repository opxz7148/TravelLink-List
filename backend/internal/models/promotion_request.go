package models

import (
	"time"
)

// PromotionRequestStatus represents possible status values for a promotion request
type PromotionRequestStatus string

const (
	PromotionStatusPending  PromotionRequestStatus = "pending"
	PromotionStatusApproved PromotionRequestStatus = "approved"
	PromotionStatusRejected PromotionRequestStatus = "rejected"
)

// String returns the string representation of the status
func (s PromotionRequestStatus) String() string {
	return string(s)
}

// ValidPromotionStatuses returns all valid promotion statuses
func ValidPromotionStatuses() []string {
	return []string{
		PromotionStatusPending.String(),
		PromotionStatusApproved.String(),
		PromotionStatusRejected.String(),
	}
}

// PromotionRequest represents a request to promote a user or plan
type PromotionRequest struct {
	// id (UUID, primary key)
	ID string `gorm:"primaryKey;type:TEXT" json:"id"`

	// user_id (UUID, foreign key to User requesting)
	UserID string `gorm:"type:TEXT;not null;index:idx_promotion_requests_user_id" json:"user_id"`

	// plan_id (UUID, foreign key to TravelPlan being promoted, nullable if upgrading user role)
	PlanID *string `gorm:"type:TEXT;index:idx_promotion_requests_plan_id" json:"plan_id"`

	// status (enum: pending | approved | rejected, default: pending)
	Status string `gorm:"type:TEXT;default:'pending';index:idx_promotion_requests_status" json:"status"`

	// admin_notes (string, max 500 chars, optional, filled on rejection or approval)
	AdminNotes string `gorm:"type:TEXT;size:500" json:"admin_notes"`

	// created_at (timestamp, UTC, immutable)
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"created_at"`

	// reviewed_at (timestamp, UTC, nullable, set when admin approves/rejects)
	ReviewedAt *time.Time `json:"reviewed_at"`

	// Relationships
	User *User        `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Plan *TravelPlan  `gorm:"foreignKey:PlanID;references:ID" json:"plan,omitempty"`
}

// TableName specifies the database table name for PromotionRequest
func (PromotionRequest) TableName() string {
	return "promotion_requests"
}

// Validate checks if the promotion request is valid
func (p *PromotionRequest) Validate() bool {
	// User ID required
	if len(p.UserID) == 0 {
		return false
	}

	// Status must be valid
	validStatuses := map[string]bool{
		PromotionStatusPending.String():  true,
		PromotionStatusApproved.String(): true,
		PromotionStatusRejected.String(): true,
	}
	if !validStatuses[p.Status] {
		return false
	}

	// Admin notes must be within 500 chars
	if len(p.AdminNotes) > 500 {
		return false
	}

	return true
}

// CanTransitionTo checks if the current status can transition to the new status
func (p *PromotionRequest) CanTransitionTo(newStatus string) bool {
	// Only pending can transition to approved or rejected
	if p.Status != PromotionStatusPending.String() {
		return false
	}

	return newStatus == PromotionStatusApproved.String() || newStatus == PromotionStatusRejected.String()
}
