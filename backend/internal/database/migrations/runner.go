package migrations

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Migration represents a single database migration
type Migration struct {
	Version   string
	Name      string
	SQL       string
	AppliedAt *time.Time
}

// MigrationRunner handles executing database migrations
type MigrationRunner struct {
	db             *gorm.DB
	migrationsPath string
	migrations     []Migration
}

// NewMigrationRunner creates a new migration runner instance
func NewMigrationRunner(db *gorm.DB, migrationsPath string) *MigrationRunner {
	return &MigrationRunner{
		db:             db,
		migrationsPath: migrationsPath,
		migrations:     []Migration{},
	}
}

// Run executes all pending migrations
func (mr *MigrationRunner) Run() error {
	log.Println("Starting migration runner...")

	// Create schema_migrations table if it doesn't exist
	if err := mr.ensureMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Load migration files from disk
	if err := mr.loadMigrations(); err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// Get list of applied migrations
	appliedVersions, err := mr.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Run pending migrations
	for _, migration := range mr.migrations {
		// Check if migration has already been applied
		if contains(appliedVersions, migration.Version) {
			log.Printf("⏭ Migration %s already applied, skipping\n", migration.Version)
			continue
		}

		// Execute migration
		log.Printf("⬆ Running migration %s: %s\n", migration.Version, migration.Name)
		if err := mr.applyMigration(&migration); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", migration.Version, err)
		}
	}

	log.Println("✓ All migrations completed successfully")
	return nil
}

// ensureMigrationsTable creates the schema_migrations table if it doesn't exist
func (mr *MigrationRunner) ensureMigrationsTable() error {
	const schema = `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`
	return mr.db.Exec(schema).Error
}

// loadMigrations loads all SQL migration files from the migrations directory
func (mr *MigrationRunner) loadMigrations() error {
	// Check if migrations directory exists
	if _, err := os.Stat(mr.migrationsPath); os.IsNotExist(err) {
		log.Printf("Migrations directory does not exist: %s. No migrations to run.\n", mr.migrationsPath)
		return nil
	}

	// Read all files in migrations directory
	entries, err := os.ReadDir(mr.migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Filter SQL files and sort by version
	var sqlFiles []os.DirEntry
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			sqlFiles = append(sqlFiles, entry)
		}
	}

	// Sort by filename (which ensures version order: 001, 002, 003, etc.)
	sort.Slice(sqlFiles, func(i, j int) bool {
		return sqlFiles[i].Name() < sqlFiles[j].Name()
	})

	// Read and parse each migration file
	for _, file := range sqlFiles {
		filePath := filepath.Join(mr.migrationsPath, file.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file.Name(), err)
		}

		// Extract version and name from filename
		parts := strings.SplitN(file.Name(), "_", 2)
		if len(parts) < 2 {
			log.Printf("Warning: skipping migration file with invalid name format: %s\n", file.Name())
			continue
		}

		version := parts[0]
		name := strings.TrimSuffix(parts[1], ".sql")

		migration := Migration{
			Version: version,
			Name:    name,
			SQL:     string(content),
		}

		mr.migrations = append(mr.migrations, migration)
		log.Printf("Loaded migration: %s - %s\n", version, name)
	}

	return nil
}

// applyMigration executes a single migration and records it in schema_migrations
// Uses a database transaction to ensure atomicity: both SQL and record must succeed or both fail
func (mr *MigrationRunner) applyMigration(migration *Migration) error {
	// Use GORM's Transaction mechanism for atomicity
	err := mr.db.Transaction(func(tx *gorm.DB) error {
		// Execute the migration SQL within transaction
		if err := tx.Exec(migration.SQL).Error; err != nil {
			return fmt.Errorf("SQL error: %w", err)
		}

		// Record migration as applied within same transaction
		now := time.Now().UTC()
		migration.AppliedAt = &now

		if err := tx.Exec(
			"INSERT INTO schema_migrations (version, name, applied_at) VALUES (?, ?, ?)",
			migration.Version,
			migration.Name,
			now,
		).Error; err != nil {
			return fmt.Errorf("failed to record migration: %w", err)
		}

		return nil
	})

	return err
}

// getAppliedMigrations returns a list of migration versions that have been applied
func (mr *MigrationRunner) getAppliedMigrations() ([]string, error) {
	var versions []string
	if err := mr.db.Raw("SELECT version FROM schema_migrations ORDER BY applied_at ASC").Scan(&versions).Error; err != nil {
		return nil, err
	}
	return versions, nil
}

// contains checks if a slice contains a string value
func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

// GetMigrationStatus returns the status of all migrations
func (mr *MigrationRunner) GetMigrationStatus() (map[string]bool, error) {
	status := make(map[string]bool)

	// Mark all loaded migrations as not applied
	for _, m := range mr.migrations {
		status[m.Version] = false
	}

	// Check which have been applied
	applied, err := mr.getAppliedMigrations()
	if err != nil {
		return nil, err
	}

	for _, version := range applied {
		status[version] = true
	}

	return status, nil
}
