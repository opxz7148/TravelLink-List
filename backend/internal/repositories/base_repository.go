package repositories

import (
	"context"
	"errors"

	"tll-backend/internal/database"
	"tll-backend/internal/models"

	"gorm.io/gorm"
)

// BaseRepository provides common database operations for all relational repositories
// Embeds the database service and provides generic CRUD helpers using Go generics
type BaseRepository struct {
	dbService database.Service
}

// NewBaseRepository creates a new base repository
func NewBaseRepository(dbService database.Service) *BaseRepository {
	return &BaseRepository{
		dbService: dbService,
	}
}

// getDB returns the underlying GORM database instance
func (b *BaseRepository) getDB() *gorm.DB {
	return b.dbService.GetDB()
}

// isRecordNotFound checks if an error is a GORM "record not found" error
func (b *BaseRepository) isRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// convertNotFoundError converts GORM "record not found" error to models.ErrNotFound
func (b *BaseRepository) convertNotFoundError(err error) error {
	if b.isRecordNotFound(err) {
		return models.ErrNotFound
	}
	return err
}

// FindFirst retrieves a single record matching the query condition
// Returns models.ErrNotFound if no record matches
func (b *BaseRepository) FindFirst(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	// Combine query and args for GORM's variadic signature
	conditions := append([]interface{}{query}, args...)
	if err := b.getDB().WithContext(ctx).First(dest, conditions...).Error; err != nil {
		return b.convertNotFoundError(err)
	}
	return nil
}

// FindMany retrieves multiple records matching the query condition
// dest should be a pointer to a slice (e.g., *[]*models.User)
func (b *BaseRepository) FindMany(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	// Combine query and args for GORM's variadic signature
	conditions := append([]interface{}{query}, args...)
	if err := b.getDB().WithContext(ctx).Where(conditions[0], conditions[1:]...).Find(dest).Error; err != nil {
		return err
	}
	return nil
}

// UpdateFields updates one or more fields of a record identified by the where clause
// updates is a map of field names to new values
// Returns models.ErrNotFound if no record matches the query
func (b *BaseRepository) UpdateFields(ctx context.Context, modelType interface{}, updates map[string]interface{}, query string, args ...interface{}) error {
	conditions := append([]interface{}{query}, args...)
	result := b.getDB().WithContext(ctx).Model(modelType).
		Where(conditions[0], conditions[1:]...).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return models.ErrNotFound
	}

	return nil
}

// UpdateField updates a single field of a record identified by the where clause
// field is the field name, value is the new value
// Returns models.ErrNotFound if no record matches the query
func (b *BaseRepository) UpdateField(ctx context.Context, modelType interface{}, field string, value interface{}, query string, args ...interface{}) error {
	conditions := append([]interface{}{query}, args...)
	result := b.getDB().WithContext(ctx).Model(modelType).
		Where(conditions[0], conditions[1:]...).
		Update(field, value)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return models.ErrNotFound
	}

	return nil
}

// UpdateFieldWithExpr updates a single field using a GORM expression
// Useful for updates like: UPDATE table SET last_login_at = CURRENT_TIMESTAMP
// Returns models.ErrNotFound if no record matches the query
func (b *BaseRepository) UpdateFieldWithExpr(ctx context.Context, modelType interface{}, field string, expr interface{}, query string, args ...interface{}) error {
	conditions := append([]interface{}{query}, args...)
	result := b.getDB().WithContext(ctx).Model(modelType).
		Where(conditions[0], conditions[1:]...).
		Update(field, expr)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return models.ErrNotFound
	}

	return nil
}

// HardDelete permanently deletes records matching the query condition
// Returns models.ErrNotFound if no record matches the query
func (b *BaseRepository) HardDelete(ctx context.Context, modelType interface{}, query string, args ...interface{}) error {
	conditions := append([]interface{}{query}, args...)
	result := b.getDB().WithContext(ctx).Where(conditions[0], conditions[1:]...).Delete(modelType)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return models.ErrNotFound
	}

	return nil
}

// Count returns the count of records matching the query condition
func (b *BaseRepository) Count(ctx context.Context, modelType interface{}, query string, args ...interface{}) (int64, error) {
	var count int64
	conditions := append([]interface{}{query}, args...)
	if err := b.getDB().WithContext(ctx).Model(modelType).
		Where(conditions[0], conditions[1:]...).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Exists checks if a record exists matching the query condition
func (b *BaseRepository) Exists(ctx context.Context, modelType interface{}, query string, args ...interface{}) (bool, error) {
	count, err := b.Count(ctx, modelType, query, args...)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
