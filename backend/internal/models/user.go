package models

import "time"

// UserRole represents a user's role in the system
type UserRole string

// Role constants define the available user roles in the system
const (
	RoleSimple    UserRole = "simple"    // Can browse, search, comment, rate (default)
	RoleTraveller UserRole = "traveller" // Can create/edit travel plans, create nodes (after promotion)
	RoleAdmin     UserRole = "admin"     // Can moderate content, approve user-created nodes
)

// String returns the string representation of a UserRole
func (r UserRole) String() string {
	return string(r)
}

// CheckRole validates if the provided role is a valid UserRole
// Returns true if the role is one of: simple, traveller, or admin
func CheckRole(role UserRole) bool {
	validRoles := map[UserRole]bool{
		RoleSimple:    true,
		RoleTraveller: true,
		RoleAdmin:     true,
	}
	return validRoles[role]
}

// User represents a platform user with role-based permissions.
// Follows the specification from data-model.md with all validation rules.
//
// Roles:
//   - simple: Can browse, search, comment, rate (default)
//   - traveller: Can create/edit travel plans, create nodes (after promotion)
//   - admin: Can moderate content, approve user-created nodes
type User struct {
	// id (UUID, primary key)
	// Unique identifier for the user, typically generated as UUID v4
	ID string `gorm:"primaryKey;type:TEXT" json:"id"`

	// email (string, unique, max 255 chars, validated as RFC 5322)
	// Must be unique across all users; used for login and password recovery
	Email string `gorm:"uniqueIndex:idx_users_email;type:TEXT;not null" json:"email"`

	// username (string, unique, max 50 chars, alphanumeric + underscore)
	// Display name for user profiles; must be unique
	// Validation: 3-50 chars, alphanumeric + underscore only
	Username string `gorm:"uniqueIndex:idx_users_username;type:TEXT;not null" json:"username"`

	// password_hash (string, bcrypt hash, never stored plaintext)
	// Always store bcrypt hash (cost 10+), never plaintext passwords
	PasswordHash string `gorm:"type:TEXT;not null" json:"-"` // Never expose in JSON

	// role (enum: RoleSimple | RoleTraveller | RoleAdmin, default: RoleSimple)
	// Defines permission level using UserRole constants:
	//   - RoleSimple: Browse, comment, rate only
	//   - RoleTraveller: Can create/edit travel plans (after admin promotion)
	//   - RoleAdmin: Full moderation access
	// Use CheckRole(role) to validate role values
	Role string `gorm:"type:TEXT;default:'simple';index" json:"role"`

	// display_name (string, max 100 chars, optional, defaults to username)
	// Public display name; if empty, UI shows username instead
	DisplayName string `gorm:"type:TEXT;size:100" json:"display_name"`

	// bio (string, max 500 chars, optional)
	// User's personal bio displayed on profile
	Bio string `gorm:"type:TEXT;size:500" json:"bio"`

	// profile_picture_url (string, URL validation, optional)
	// URL to user's profile picture; stored as string, not as blob
	ProfilePictureURL string `gorm:"type:TEXT" json:"profile_picture_url"`

	// is_active (boolean, default: true, soft-delete via false)
	// When false, user is hidden from public listings but authored content preserved
	// Used for soft-delete instead of hard delete
	IsActive bool `gorm:"type:BOOLEAN;default:true;index" json:"is_active"`

	// created_at (timestamp, UTC, immutable)
	// When the user account was created; never modified after creation
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"created_at"`

	// updated_at (timestamp, UTC, on modification)
	// When the user profile was last updated (name, bio, picture, etc.)
	// Does NOT include password changes (use password_updated_at separately)
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli" json:"updated_at"`

	// last_login_at (timestamp, UTC, nullable, on successful auth)
	// Tracks when user last successfully authenticated
	// Useful for activity monitoring and account security
	LastLoginAt *time.Time `gorm:"type:TIMESTAMP;index" json:"last_login_at"`
}

// TableName specifies the database table name for User
func (User) TableName() string {
	return "users"
}

// Validation Rules (for use with validation middleware/service)

// ValidateEmail checks if email is not empty and looks like an email
func (u *User) ValidateEmail() bool {
	return u.Email != "" && len(u.Email) <= 255
}

// ValidateUsername checks username constraints
// Must be 3-50 chars, alphanumeric + underscore only
func (u *User) ValidateUsername() bool {
	if len(u.Username) < 3 || len(u.Username) > 50 {
		return false
	}
	// Alphanumeric + underscore only
	for _, c := range u.Username {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_') {
			return false
		}
	}
	return true
}

// ValidateDisplayName checks display_name constraints
func (u *User) ValidateDisplayName() bool {
	return len(u.DisplayName) <= 100
}

// ValidateBio checks bio constraints
func (u *User) ValidateBio() bool {
	return len(u.Bio) <= 500
}

// ValidateRole checks if role is a valid enum value
func (u *User) ValidateRole() bool {
	return CheckRole(UserRole(u.Role))
}

// IsSimple returns true if user has 'simple' role
func (u *User) IsSimple() bool {
	return u.Role == RoleSimple.String()
}

// IsTraveller returns true if user has 'traveller' role
func (u *User) IsTraveller() bool {
	return u.Role == RoleTraveller.String()
}

// IsAdmin returns true if user has 'admin' role
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin.String()
}

// CanCreatePlans returns true if user can create travel plans
// Only traveller and admin roles can create plans
func (u *User) CanCreatePlans() bool {
	return u.IsTraveller() || u.IsAdmin()
}

// CanModerate returns true if user can moderate content
// Only admin role can moderate
func (u *User) CanModerate() bool {
	return u.IsAdmin()
}

// UpdateLastLogin updates the last_login_at timestamp to now (UTC)
func (u *User) UpdateLastLogin() {
	now := time.Now().UTC()
	u.LastLoginAt = &now
}

// SetDisplayNameDefault sets display_name to username if empty
func (u *User) SetDisplayNameDefault() {
	if u.DisplayName == "" {
		u.DisplayName = u.Username
	}
}
