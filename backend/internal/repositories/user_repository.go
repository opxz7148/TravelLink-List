package repositories

import (
	"context"
	"tll-backend/internal/models"
)

// UserRepository defines the interface for user data access operations
// All methods follow the repository pattern: data access abstraction layer
// This allows implementation switching (e.g., SQLite to PostgreSQL) without affecting services
type UserRepository interface {
	// CreateAndSave inserts a new user into the database and persists it
	// Returns error if email or username already exists (unique constraint violation)
	CreateAndSave(ctx context.Context, user *models.User) error

	// GetByID retrieves a user by their unique ID
	// Returns models.ErrNotFound if user doesn't exist
	GetByID(ctx context.Context, id string) (*models.User, error)

	// GetByEmail retrieves a user by their email address
	// Email lookups are case-sensitive
	// Returns models.ErrNotFound if user doesn't exist
	GetByEmail(ctx context.Context, email string) (*models.User, error)

	// GetByUsername retrieves a user by their username
	// Username lookups are case-sensitive
	// Returns models.ErrNotFound if user doesn't exist
	GetByUsername(ctx context.Context, username string) (*models.User, error)

	// Update updates an existing user record
	// Only updates fields that are non-zero; uses GORM's selective updates
	// Returns error if user doesn't exist
	Update(ctx context.Context, user *models.User) error

	// UpdatePassword updates a user's password hash
	// Useful for password change operations without modifying other fields
	// Returns error if user doesn't exist
	UpdatePassword(ctx context.Context, userID string, hashedPassword string) error

	// UpdateLastLogin updates the user's last_login_at timestamp to now
	// Used after successful authentication
	// Returns error if user doesn't exist
	UpdateLastLogin(ctx context.Context, userID string) error

	// Delete performs a soft delete by setting is_active to false
	// User record remains in database but is hidden from public listings
	Delete(ctx context.Context, userID string) error

	// HardDelete permanently removes a user record from the database
	// Use with caution: violates foreign key constraints if user has related data
	// Should only be used for admin/cleanup operations
	HardDelete(ctx context.Context, userID string) error

	// Restore reactivates a soft-deleted user (sets is_active to true)
	// No-op if user is already active
	Restore(ctx context.Context, userID string) error

	// ExistsByEmail checks if a user with the given email exists
	// Useful for frontend validation and duplicate prevention
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// ExistsByUsername checks if a user with the given username exists
	// Useful for frontend validation and duplicate prevention
	ExistsByUsername(ctx context.Context, username string) (bool, error)

	// GetActiveUsers retrieves all active users with pagination
	// offset: number of records to skip
	// limit: maximum number of records to return (0 = no limit)
	// Returns list of active users sorted by created_at descending
	GetActiveUsers(ctx context.Context, offset, limit int) ([]*models.User, error)

	// GetUsersByRole retrieves all users with a specific role
	// role: UserRole constant (RoleSimple, RoleTraveller, RoleAdmin)
	// Returns list of users with the specified role
	GetUsersByRole(ctx context.Context, role models.UserRole) ([]*models.User, error)

	// CountActiveUsers returns the total number of active users
	CountActiveUsers(ctx context.Context) (int64, error)

	// CountUsersByRole returns the count of users with a specific role
	CountUsersByRole(ctx context.Context, role models.UserRole) (int64, error)

	// Search searches users by email or username (partial match, case-insensitive)
	// query: search term (minimum 2 characters)
	// limit: maximum number of results to return
	// Returns matching users sorted by relevance
	Search(ctx context.Context, query string, limit int) ([]*models.User, error)

	// PromoteToTraveller upgrades a user's role from 'simple' to 'traveller'
	// No-op if user already has traveller or admin role
	PromoteToTraveller(ctx context.Context, userID string) error

	// DemoteToSimple downgrades a user's role to 'simple'
	// Use with caution: may affect user's travel plans and permissions
	DemoteToSimple(ctx context.Context, userID string) error

	// MakeAdmin upgrades a user's role to 'admin'
	// Only use for trusted users
	MakeAdmin(ctx context.Context, userID string) error
}
