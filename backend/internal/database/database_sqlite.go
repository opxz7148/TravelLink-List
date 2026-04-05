package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"tll-backend/internal/database/migrations"
)

// SQLiteService represents a SQLite3-based database service.
// It implements the Service interface using SQLite3 as the underlying database.
type SQLiteService struct {
	db     *gorm.DB
	dbPath string
}

var (
	sqliteInstance *SQLiteService
)

// NewSQLite creates a new SQLite3 database service.
// It initializes a SQLite3 database connection with GORM.
// If a connection already exists (singleton pattern), it returns the existing instance.
//
// Environment Variables:
//   - DB_PATH: Path to the SQLite database file (default: "backend/travellink.db")
//
// Returns a Service interface implementation using SQLite3.
func NewSQLite() Service {
	// Reuse Connection (singleton pattern)
	if sqliteInstance != nil {
		return sqliteInstance
	}

	// Get database path from environment or use default
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "backend/travellink.db"
	}

	// Ensure the directory exists
	dbDir := filepath.Dir(dbPath)
	if dbDir != "." {
		if err := os.MkdirAll(dbDir, 0o755); err != nil {
			log.Fatalf("Failed to create database directory: %v", err)
		}
	}

	// Configure GORM logger
	gormLogger := logger.Default.LogMode(logger.Info)

	// Open SQLite database with GORM
	// The DSN format for SQLite with GORM is just the file path
	// Additional parameters can be added using ?param=value syntax
	dsn := fmt.Sprintf("file:%s?cache=shared&mode=rwc&_journal_mode=WAL", dbPath)

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		log.Fatalf("Failed to connect to SQLite database: %v", err)
	}

	// Get underlying SQL DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get SQL DB instance: %v", err)
	}

	// Enable foreign key constraints for SQLite3
	if _, err := sqlDB.Exec("PRAGMA foreign_keys = ON"); err != nil {
		log.Fatalf("Failed to enable foreign key constraints: %v", err)
	}

	// Configure connection pool for SQLite
	// SQLite has different concurrency characteristics than PostgreSQL
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetConnMaxLifetime(time.Hour)

	sqliteInstance = &SQLiteService{
		db:     db,
		dbPath: dbPath,
	}

	log.Printf("Connected to SQLite database at: %s", dbPath)

	// Run database migrations
	// Get current working directory to resolve migrations path
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Warning: failed to get current working directory for migrations: %v", err)
		cwd = "."
	}

	if err := migrations.InitMigrations(db, cwd); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	return sqliteInstance
}

// Health checks the health of the SQLite database connection.
// It returns a map with keys indicating various health statistics.
// For SQLite3, some metrics differ from traditional databases since it's file-based.
func (s *SQLiteService) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Get underlying SQL DB from GORM
	sqlDB, err := s.db.DB()
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Printf("db down: %v", err)
		return stats
	}

	// Ping the database
	err = sqlDB.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Printf("db down: %v", err)
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"
	stats["database"] = s.dbPath

	// Get database stats
	dbStats := sqlDB.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)

	// Check file size
	if fileInfo, err := os.Stat(s.dbPath); err == nil {
		stats["database_size_bytes"] = strconv.FormatInt(fileInfo.Size(), 10)
		stats["database_size_mb"] = fmt.Sprintf("%.2f", float64(fileInfo.Size())/1024/1024)
	}

	// SQLite-specific pragma information
	var journalMode string
	if err := s.db.Raw("PRAGMA journal_mode").Scan(&journalMode).Error; err == nil {
		stats["journal_mode"] = journalMode
	}

	var syncMode string
	if err := s.db.Raw("PRAGMA synchronous").Scan(&syncMode).Error; err == nil {
		stats["synchronous"] = syncMode
	}

	// Check for lock status by attempting a query
	var tableCount int64
	if err := s.db.WithContext(ctx).Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table'").Scan(&tableCount).Error; err != nil {
		stats["message"] = "Database is locked or unavailable"
	} else {
		stats["table_count"] = strconv.FormatInt(tableCount, 10)
	}

	return stats
}

// Close closes the SQLite database connection.
// It logs a message indicating the disconnection from the database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *SQLiteService) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	log.Printf("Disconnected from SQLite database at: %s", s.dbPath)
	return sqlDB.Close()
}

// GetDB returns the underlying GORM DB instance.
// This allows other parts of the application to perform database operations using GORM.
func (s *SQLiteService) GetDB() *gorm.DB {
	return s.db
}
