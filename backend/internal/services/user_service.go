package services

import (
	"context"
	"regexp"
	"strings"

	"tll-backend/internal/models"
	"tll-backend/internal/repositories"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserService defines the interface for user-related business logic operations
type UserService interface {
	// Register creates a new user account with validation and password hashing
	// Returns user and JWT token response on success
	Register(ctx context.Context, email, username, password, displayName string) (*models.User, *TokenResponse, error)

	// Login authenticates a user by email and password
	// Returns user and JWT token response on success
	Login(ctx context.Context, email, password string) (*models.User, *TokenResponse, error)

	// GetUserByID retrieves a user by their ID
	GetUserByID(ctx context.Context, userID string) (*models.User, error)

	// GetUserByEmail retrieves a user by their email address
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)

	// GetUserByUsername retrieves a user by their username
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)

	// UpdateProfile updates user profile information
	UpdateProfile(ctx context.Context, userID, displayName, bio, profilePictureURL string) (*models.User, error)

	// ChangePassword changes a user's password
	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error

	// PromoteToTraveller upgrades a user from 'simple' role to 'traveller' role
	PromoteToTraveller(ctx context.Context, userID string) error

	// DemoteToSimple downgrades a user to 'simple' role
	DemoteToSimple(ctx context.Context, userID string) error

	// MakeAdmin upgrades a user to 'admin' role
	MakeAdmin(ctx context.Context, userID string) error

	// DeleteUser performs a soft delete
	DeleteUser(ctx context.Context, userID string) error
}

// RelationalUserService implements UserService with relational database backend
type RelationalUserService struct {
	userRepo   repositories.UserRepository
	jwtService JWTService
}

// NewRelationalUserService creates a new user service
func NewRelationalUserService(userRepo repositories.UserRepository, jwtService JWTService) UserService {
	return &RelationalUserService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

const (
	minPasswordLength = 8
	bcryptCost        = 12
)

// ValidateEmail checks if email is valid RFC 5322 format
func ValidateEmail(email string) bool {
	if len(email) > 255 || len(email) == 0 {
		return false
	}
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailPattern, email)
	return matched
}

// ValidatePassword checks password meets security requirements
func ValidatePassword(password string) bool {
	if len(password) < minPasswordLength {
		return false
	}

	hasUpper := false
	hasDigit := false
	hasSpecial := false

	for _, c := range password {
		switch {
		case c >= 'A' && c <= 'Z':
			hasUpper = true
		case c >= '0' && c <= '9':
			hasDigit = true
		case strings.ContainsAny(string(c), "!@#$%^&*()-_=+[]{}|;:,.<>?"):
			hasSpecial = true
		}
	}

	return hasUpper && hasDigit && hasSpecial
}

// HashPassword creates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ComparePasswordHash verifies a password against its bcrypt hash
func ComparePasswordHash(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// Register creates a new user account
func (s *RelationalUserService) Register(ctx context.Context, email, username, password, displayName string) (*models.User, *TokenResponse, error) {
	// Validate inputs
	if !ValidateEmail(email) {
		return nil, nil, models.ErrValidation
	}

	email = strings.ToLower(strings.TrimSpace(email))
	username = strings.TrimSpace(username)
	displayName = strings.TrimSpace(displayName)

	tempUser := &models.User{Username: username}
	if !tempUser.ValidateUsername() {
		return nil, nil, models.ErrValidation
	}

	if !ValidatePassword(password) {
		return nil, nil, models.ErrValidation
	}

	// Check for duplicates
	exists, err := s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, nil, err
	}
	if exists {
		return nil, nil, models.ErrDuplicateEmail
	}

	exists, err = s.userRepo.ExistsByUsername(ctx, username)
	if err != nil {
		return nil, nil, err
	}
	if exists {
		return nil, nil, models.ErrDuplicateUsername
	}

	// Hash password
	passwordHash, err := HashPassword(password)
	if err != nil {
		return nil, nil, err
	}

	// Create user
	user := &models.User{
		ID:           uuid.New().String(),
		Email:        email,
		Username:     username,
		PasswordHash: passwordHash,
		DisplayName:  displayName,
		Role:         models.RoleSimple.String(),
		IsActive:     true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, nil, err
	}

	// Generate JWT token
	_, tokenResponse, err := s.jwtService.GenerateToken(user)
	if err != nil {
		return nil, nil, err
	}

	return user, tokenResponse, nil
}

// Login authenticates a user
func (s *RelationalUserService) Login(ctx context.Context, email, password string) (*models.User, *TokenResponse, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, nil, err
	}

	// Verify password
	if !ComparePasswordHash(user.PasswordHash, password) {
		return nil, nil, models.ErrValidation
	}

	// Update last login timestamp
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		return nil, nil, err
	}

	// Generate JWT token
	_, tokenResponse, err := s.jwtService.GenerateToken(user)
	if err != nil {
		return nil, nil, err
	}

	return user, tokenResponse, nil
}

// GetUserByID retrieves a user by ID
func (s *RelationalUserService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

// GetUserByEmail retrieves a user by email
func (s *RelationalUserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	return s.userRepo.GetByEmail(ctx, email)
}

// GetUserByUsername retrieves a user by username
func (s *RelationalUserService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	return s.userRepo.GetByUsername(ctx, strings.TrimSpace(username))
}

// UpdateProfile updates user profile information
func (s *RelationalUserService) UpdateProfile(ctx context.Context, userID, displayName, bio, profilePictureURL string) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.DisplayName = strings.TrimSpace(displayName)
	user.Bio = strings.TrimSpace(bio)
	user.ProfilePictureURL = strings.TrimSpace(profilePictureURL)

	if !user.ValidateDisplayName() || !user.ValidateBio() {
		return nil, models.ErrValidation
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword changes a user's password
func (s *RelationalUserService) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify old password
	if !ComparePasswordHash(user.PasswordHash, oldPassword) {
		return models.ErrValidation
	}

	// Validate new password
	if !ValidatePassword(newPassword) {
		return models.ErrValidation
	}

	// Hash new password
	newHash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(ctx, userID, newHash)
}

// PromoteToTraveller upgrades a user to traveller role
func (s *RelationalUserService) PromoteToTraveller(ctx context.Context, userID string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Only promote if currently simple role
	if !user.IsSimple() {
		return nil // No-op if already traveller or admin
	}

	return s.userRepo.PromoteToTraveller(ctx, userID)
}

// DemoteToSimple downgrades a user to simple role
func (s *RelationalUserService) DemoteToSimple(ctx context.Context, userID string) error {
	return s.userRepo.DemoteToSimple(ctx, userID)
}

// MakeAdmin upgrades a user to admin role
func (s *RelationalUserService) MakeAdmin(ctx context.Context, userID string) error {
	return s.userRepo.MakeAdmin(ctx, userID)
}

// DeleteUser performs a soft delete
func (s *RelationalUserService) DeleteUser(ctx context.Context, userID string) error {
	return s.userRepo.Delete(ctx, userID)
}
