package models

import (
	"strings"
	"time"
)

// TravelPlanStatus represents the status of a travel plan
type TravelPlanStatus string

// Travel plan status constants
const (
	TravelPlanStatusDraft     TravelPlanStatus = "draft"
	TravelPlanStatusPublished TravelPlanStatus = "published"
	TravelPlanStatusSuspended TravelPlanStatus = "suspended"
	TravelPlanStatusDeleted   TravelPlanStatus = "deleted"
)

// String returns the string representation of a TravelPlanStatus
func (s TravelPlanStatus) String() string {
	return string(s)
}

// CheckTravelPlanStatus validates if the provided status is a valid TravelPlanStatus
func CheckTravelPlanStatus(status TravelPlanStatus) bool {
	validStatuses := map[TravelPlanStatus]bool{
		TravelPlanStatusDraft:     true,
		TravelPlanStatusPublished: true,
		TravelPlanStatusSuspended: true,
		TravelPlanStatusDeleted:   true,
	}
	return validStatuses[status]
}

// TravelPlan represents a complete travel itinerary as a linked list of nodes.
// Follows the specification from data-model.md with all validation rules.
//
// A travel plan organizes a sequence of attractions and transitions into a coherent trip.
// Status can be draft (private), published (public), or suspended (admin action).
type TravelPlan struct {
	// id (UUID, primary key)
	// Unique identifier for the travel plan, typically generated as UUID v4
	ID string `gorm:"primaryKey;type:TEXT" json:"id"`

	// title (string, max 150 chars, required)
	// Display name for the travel plan
	Title string `gorm:"type:TEXT;not null;size:150;index" json:"title"`

	// description (string, max 1000 chars, optional)
	// Detailed explanation of the travel plan
	Description string `gorm:"type:TEXT;size:1000" json:"description"`

	// destination (string, max 200 chars, required; searchable, not normalized location)
	// Primary destination or region being traveled
	// Used for search/filtering, not strict geocoding
	Destination string `gorm:"type:TEXT;not null;size:200;index" json:"destination"`

	// author_id (UUID, foreign key to User, denormalized for query efficiency)
	// The user who created and owns this travel plan
	AuthorID string `gorm:"type:TEXT;not null;index" json:"author_id"`

	// status (enum: draft | published | suspended, default: draft)
	// Visibility and lifecycle state of the plan
	// Only published plans visible to non-authors
	Status string `gorm:"type:TEXT;not null;default:'draft';index" json:"status"`

	// rating_count (integer, default: 0, denormalized for fast retrieval)
	// Number of ratings received (denormalized for efficient listing)
	RatingCount int `gorm:"type:INTEGER;default:0;not null" json:"rating_count"`

	// rating_sum (integer, default: 0, used to calculate average)
	// Sum of all star ratings (used to calculate average rating)
	// Average = rating_sum / rating_count (when count > 0)
	RatingSum int `gorm:"type:INTEGER;default:0;not null" json:"rating_sum"`

	// comment_count (integer, default: 0, denormalized for display)
	// Number of comments on this plan (denormalized for efficient listing)
	CommentCount int `gorm:"type:INTEGER;default:0;not null" json:"comment_count"`

	// is_deleted_by_admin (boolean, default: false, soft-delete flag)
	// When true, plan is hidden from listings but preserved for history
	IsDeletedByAdmin bool `gorm:"type:BOOLEAN;default:false;not null" json:"is_deleted_by_admin"`

	// created_at (timestamp, UTC, immutable)
	// When the travel plan was created
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"created_at"`

	// updated_at (timestamp, UTC, on modification or node reordering)
	// When the travel plan was last modified
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli" json:"updated_at"`
}

// TableName specifies the database table name for TravelPlan
func (TravelPlan) TableName() string {
	return "travel_plans"
}

// Validate performs domain validation on the TravelPlan
// Returns ErrValidation if any rule is violated
func (tp *TravelPlan) Validate() error {
	// Title validation: Required, 1-150 chars, non-empty after trim
	if strings.TrimSpace(tp.Title) == "" {
		return ErrValidation
	}
	if len(tp.Title) > 150 {
		return ErrValidation
	}

	// Destination validation: Required, 1-200 chars, enables search but not strict geocoding
	if strings.TrimSpace(tp.Destination) == "" {
		return ErrValidation
	}
	if len(tp.Destination) > 200 {
		return ErrValidation
	}

	// Description validation: Optional but max 1000 chars if provided
	if len(tp.Description) > 1000 {
		return ErrValidation
	}

	// Status validation: must be one of the valid statuses
	if !CheckTravelPlanStatus(TravelPlanStatus(tp.Status)) {
		return ErrValidation
	}

	// AuthorID validation: Must not be empty
	if strings.TrimSpace(tp.AuthorID) == "" {
		return ErrValidation
	}

	return nil
}

// HasStatus checks if the plan has a specific status
func (tp *TravelPlan) HasStatus(status TravelPlanStatus) bool {
	return tp.Status == status.String()
}

// IsPublished returns true if the plan status is published
func (tp *TravelPlan) IsPublished() bool {
	return tp.HasStatus(TravelPlanStatusPublished)
}

// IsDraft returns true if the plan status is draft
func (tp *TravelPlan) IsDraft() bool {
	return tp.HasStatus(TravelPlanStatusDraft)
}

// IsSuspended returns true if the plan status is suspended
func (tp *TravelPlan) IsSuspended() bool {
	return tp.HasStatus(TravelPlanStatusSuspended)
}

// userCanModify checks if a user has authorization to modify this plan
// Returns (isAuthor, isAdmin)
func (tp *TravelPlan) userCanModify(userID string, userRole UserRole) (bool, bool) {
	return tp.AuthorID == userID, userRole == RoleAdmin
}

// CanBeViewedBy checks if a user can view this plan
// Published plans visible to anyone; draft/suspended visible only to author or admin
func (tp *TravelPlan) CanBeViewedBy(userID string, userRole UserRole) bool {
	// Any non-deleted published plan can be viewed by anyone
	if tp.IsPublished() && !tp.IsDeletedByAdmin {
		return true
	}

	// Only author or admin can view draft/suspended/deleted plans
	isAuthor, isAdmin := tp.userCanModify(userID, userRole)
	return isAuthor || isAdmin
}

// CanBeEditedBy checks if a user can edit this plan
// Only author or admin can edit; cannot edit published plans except admin
func (tp *TravelPlan) CanBeEditedBy(userID string, userRole UserRole) bool {
	isAuthor, isAdmin := tp.userCanModify(userID, userRole)

	// Admin can always edit
	if isAdmin {
		return true
	}

	// Author can edit only if plan is draft
	if isAuthor && tp.IsDraft() {
		return true
	}

	return false
}
