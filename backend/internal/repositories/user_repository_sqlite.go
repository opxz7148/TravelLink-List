package repositories

import (
	"context"
	"errors"

	"tll-backend/internal/database"
	"tll-backend/internal/models"

	"gorm.io/gorm"
)

// RelationalUserRepository implements UserRepository using GORM for relational database access
type RelationalUserRepository struct {
	dbService database.Service
}

// NewRelationalUserRepository creates a new relational database user repository
// Takes a database.Service which provides access to the underlying GORM DB instance
func NewRelationalUserRepository(dbService database.Service) UserRepository {
	return &RelationalUserRepository{dbService: dbService}
}

// getDB is a helper method to get the underlying GORM DB instance
func (r *RelationalUserRepository) getDB() *gorm.DB {
	return r.dbService.GetDB()
}

// Create inserts a new user into the database
func (r *RelationalUserRepository) Create(ctx context.Context, user *models.User) error {
	if err := r.getDB().WithContext(ctx).Create(user).Error; err != nil {
		// Handle constraint violations
		if err.Error() == "UNIQUE constraint failed: users.email" {
			return models.ErrDuplicateEmail
		}
		if err.Error() == "UNIQUE constraint failed: users.username" {
			return models.ErrDuplicateUsername
		}
		return err
	}
	return nil
}

// GetByID retrieves a user by their unique ID
func (r *RelationalUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	if err := r.getDB().WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by their email address
func (r *RelationalUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.getDB().WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetByUsername retrieves a user by their username
func (r *RelationalUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	if err := r.getDB().WithContext(ctx).First(&user, "username = ?", username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

// Update updates an existing user record
func (r *RelationalUserRepository) Update(ctx context.Context, user *models.User) error {
	if err := r.getDB().WithContext(ctx).Model(user).Updates(user).Error; err != nil {
		return err
	}
	return nil
}

// UpdatePassword updates a user's password hash
func (r *RelationalUserRepository) UpdatePassword(ctx context.Context, userID string, hashedPassword string) error {
	result := r.getDB().WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Update("password_hash", hashedPassword)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return models.ErrNotFound
	}
	return nil
}

// UpdateLastLogin updates the user's last_login_at timestamp to now
func (r *RelationalUserRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	result := r.getDB().WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Update("last_login_at", gorm.Expr("CURRENT_TIMESTAMP"))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return models.ErrNotFound
	}
	return nil
}

// Delete performs a soft delete by setting is_active to false
func (r *RelationalUserRepository) Delete(ctx context.Context, userID string) error {
	result := r.getDB().WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Update("is_active", false)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return models.ErrNotFound
	}
	return nil
}

// HardDelete permanently removes a user record from the database
func (r *RelationalUserRepository) HardDelete(ctx context.Context, userID string) error {
	result := r.getDB().WithContext(ctx).Where("id = ?", userID).Delete(&models.User{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return models.ErrNotFound
	}
	return nil
}

// Restore reactivates a soft-deleted user (sets is_active to true)
func (r *RelationalUserRepository) Restore(ctx context.Context, userID string) error {
	result := r.getDB().WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Update("is_active", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return models.ErrNotFound
	}
	return nil
}

// ExistsByEmail checks if a user with the given email exists
func (r *RelationalUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.getDB().WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByUsername checks if a user with the given username exists
func (r *RelationalUserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	if err := r.getDB().WithContext(ctx).Model(&models.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetActiveUsers retrieves all active users with pagination
func (r *RelationalUserRepository) GetActiveUsers(ctx context.Context, offset, limit int) ([]*models.User, error) {
	var users []*models.User
	query := r.getDB().WithContext(ctx).Where("is_active = ?", true).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetUsersByRole retrieves all users with a specific role
func (r *RelationalUserRepository) GetUsersByRole(ctx context.Context, role models.UserRole) ([]*models.User, error) {
	var users []*models.User
	if err := r.getDB().WithContext(ctx).
		Where("role = ?", role.String()).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// CountActiveUsers returns the total number of active users
func (r *RelationalUserRepository) CountActiveUsers(ctx context.Context) (int64, error) {
	var count int64
	if err := r.getDB().WithContext(ctx).Model(&models.User{}).Where("is_active = ?", true).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountUsersByRole returns the count of users with a specific role
func (r *RelationalUserRepository) CountUsersByRole(ctx context.Context, role models.UserRole) (int64, error) {
	var count int64
	if err := r.getDB().WithContext(ctx).Model(&models.User{}).Where("role = ?", role.String()).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Search searches users by email or username (partial match, case-insensitive)
func (r *RelationalUserRepository) Search(ctx context.Context, query string, limit int) ([]*models.User, error) {
	var users []*models.User
	searchPattern := "%" + query + "%"
	err := r.getDB().WithContext(ctx).
		Where("email LIKE ? OR username LIKE ?", searchPattern, searchPattern).
		Order("created_at DESC").
		Limit(limit).
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// PromoteToTraveller upgrades a user's role from 'simple' to 'traveller'
func (r *RelationalUserRepository) PromoteToTraveller(ctx context.Context, userID string) error {
	result := r.getDB().WithContext(ctx).Model(&models.User{}).
		Where("id = ? AND role = ?", userID, models.RoleSimple.String()).
		Update("role", models.RoleTraveller.String())
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		// User either doesn't exist or isn't simple role, but we still consider it successful
		return nil
	}
	return nil
}

// DemoteToSimple downgrades a user's role to 'simple'
func (r *RelationalUserRepository) DemoteToSimple(ctx context.Context, userID string) error {
	result := r.getDB().WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Update("role", models.RoleSimple.String())
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return models.ErrNotFound
	}
	return nil
}

// MakeAdmin upgrades a user's role to 'admin'
func (r *RelationalUserRepository) MakeAdmin(ctx context.Context, userID string) error {
	result := r.getDB().WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Update("role", models.RoleAdmin.String())
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return models.ErrNotFound
	}
	return nil
}

