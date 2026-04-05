package migrations

import (
	"fmt"
	"log"
	"path/filepath"

	"gorm.io/gorm"
)

// InitMigrations initializes and runs all pending database migrations
// This should be called once during application startup after the database connection is established
//
// Parameters:
//   - db: GORM database instance
//   - basePath: Base path to the migrations directory (usually project root)
//
// Example:
//
//	if err := migrations.InitMigrations(db, "."); err != nil {
//		log.Fatalf("Migration failed: %v", err)
//	}
func InitMigrations(db *gorm.DB, basePath string) error {
	// Construct path to migrations directory
	migrationsPath := filepath.Join(basePath, "internal", "database", "migrations")

	// Create migration runner
	runner := NewMigrationRunner(db, migrationsPath)

	// Run all pending migrations
	if err := runner.Run(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("✓ Database migrations completed")
	return nil
}

// GetStatus returns the status of all database migrations
// Useful for health checks and debugging
func GetStatus(db *gorm.DB, basePath string) (map[string]bool, error) {
	migrationsPath := filepath.Join(basePath, "internal", "database", "migrations")
	runner := NewMigrationRunner(db, migrationsPath)

	// Ensure migrations table exists
	if err := runner.ensureMigrationsTable(); err != nil {
		return nil, err
	}

	// Load migrations
	if err := runner.loadMigrations(); err != nil {
		return nil, err
	}

	return runner.GetMigrationStatus()
}
