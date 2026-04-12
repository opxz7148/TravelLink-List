package models

import (
	"fmt"
	"time"
)

// Rating represents user rating of a TravelPlan.
// Follows the specification from data-model.md with all validation rules.
//
// Ratings enable community engagement on travel plans.
// Each rating is a 1-5 star rating from a user on a travel plan.
// One rating per user per plan (update replaces previous).
type Rating struct {
	// id (UUID, primary key)
	// Unique identifier for the rating, typically generated as UUID v4
	ID string `gorm:"primaryKey;type:TEXT" json:"id"`

	// plan_id (UUID, foreign key to TravelPlan)
	// The travel plan being rated
	PlanID string `gorm:"type:TEXT;not null;index" json:"plan_id"`

	// user_id (UUID, foreign key to User)
	// The user who submitted this rating
	UserID string `gorm:"type:TEXT;not null;index" json:"user_id"`

	// stars (integer, 1-5, required)
	// The star rating value (1, 2, 3, 4, or 5)
	Stars int `gorm:"type:INTEGER;not null;check:stars >= 1 AND stars <= 5" json:"stars"`

	// created_at (timestamp, UTC, immutable)
	// Timestamp when rating was created
	CreatedAt time.Time `gorm:"autoCreateTime:milli;type:TIMESTAMP;not null" json:"created_at"`

	// updated_at (timestamp, UTC, shows when rating last updated)
	// Timestamp when rating was last updated (when user changes their rating)
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli;type:TIMESTAMP;not null" json:"updated_at"`

	// User is the rating author (optional, for display)
	User *User `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`

	// TravelPlan is the rated plan (optional, for display)
	TravelPlan *TravelPlan `gorm:"foreignKey:PlanID;references:ID" json:"plan,omitempty"`
}

// TableName specifies the table name for the Rating model
func (Rating) TableName() string {
	return "ratings"
}

// Validate validates rating fields according to data-model.md rules
// Returns ErrValidation if any rule is violated
func (r *Rating) Validate() error {
	// Stars validation: Required, must be 1-5
	if r.Stars < 1 || r.Stars > 5 {
		return fmt.Errorf("%w: stars must be between 1 and 5, got %d", ErrValidation, r.Stars)
	}

	// Plan ID validation: Required
	if r.PlanID == "" {
		return fmt.Errorf("%w: plan_id cannot be empty", ErrValidation)
	}

	// User ID validation: Required
	if r.UserID == "" {
		return fmt.Errorf("%w: user_id cannot be empty", ErrValidation)
	}

	return nil
}
