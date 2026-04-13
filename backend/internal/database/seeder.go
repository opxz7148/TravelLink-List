package database

import (
	"context"
	"log"
	"time"

	"tll-backend/internal/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedDatabase initializes the database with sample data if it's empty
func SeedDatabase(db *gorm.DB) error {
	ctx := context.Background()

	// Count existing users
	var userCount int64
	if err := db.Model(&models.User{}).Count(&userCount).Error; err != nil {
		log.Printf("Warning: Could not count users for seeding check: %v", err)
		return nil // Don't fail on seeding
	}

	// Only seed if database is empty
	// Ignore system user that got created during migration, so check for more than 1 user
	if userCount > 1 {
		log.Println("Database already has users, skipping seeding")
		return nil
	}

	log.Println("Database is empty, seeding with sample data...")

	// Sample users to seed
	sampleUsers := []struct {
		email       string
		username    string
		password    string
		displayName string
		role        string
		bio         string
	}{
		{
			email:       "admin@travellink.local",
			username:    "admin",
			password:    "AdminPass123!",
			displayName: "Administrator",
			role:        models.RoleAdmin.String(),
			bio:         "System administrator",
		},
		{
			email:       "traveller@travellink.local",
			username:    "traveller",
			password:    "TravellerPass123!",
			displayName: "Travel Enthusiast",
			role:        models.RoleTraveller.String(),
			bio:         "Loves exploring new places",
		},
		{
			email:       "user@travellink.local",
			username:    "user",
			password:    "UserPass123!",
			displayName: "Simple User",
			role:        models.RoleSimple.String(),
			bio:         "Just browsing travel plans",
		},
	}

	for _, userData := range sampleUsers {
		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password for %s: %v", userData.username, err)
			continue
		}

		// Create user model
		user := &models.User{
			ID:           uuid.New().String(),
			Email:        userData.email,
			Username:     userData.username,
			PasswordHash: string(hashedPassword),
			DisplayName:  userData.displayName,
			Bio:          userData.bio,
			Role:         userData.role,
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		// Create user directly in database via GORM
		if err := db.WithContext(ctx).Create(user).Error; err != nil {
			log.Printf("Warning: Could not seed user %s: %v", userData.username, err)
			continue
		}

		log.Printf("✓ Seeded user: %s (%s) - Password: %s", userData.username, userData.email, userData.password)
	}

	log.Println("Database seeding completed successfully!")
	return nil
}

