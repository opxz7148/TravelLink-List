package services

import (
	"errors"
	"time"

	"tll-backend/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

// TokenClaims represents the JWT claims for our auth tokens
type TokenClaims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// TokenResponse represents the response structure for token operations
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"` // Seconds until token expiration
}

// JWTService defines the interface for JWT token operations
type JWTService interface {
	// GenerateToken creates a new JWT token for a user
	// Returns token string and token response with metadata
	GenerateToken(user *models.User) (string, *TokenResponse, error)

	// ValidateToken verifies and parses a JWT token
	// Returns the claims if valid, error otherwise
	ValidateToken(tokenString string) (*TokenClaims, error)

	// ParseToken extracts claims from a token without validating signature
	// Useful for getting claims from context during middleware
	ParseToken(tokenString string) (*TokenClaims, error)
}

// RelationalJWTService implements JWTService using HMAC signing
type RelationalJWTService struct {
	secretKey string
	expiresIn time.Duration
}

// NewRelationalJWTService creates a new JWT service
// secretKey: secret string for HMAC signing (should be from environment)
// expiresIn: token expiration duration (typically 1 hour)
func NewRelationalJWTService(secretKey string, expiresIn time.Duration) JWTService {
	return &RelationalJWTService{
		secretKey: secretKey,
		expiresIn: expiresIn,
	}
}

// GenerateToken creates a new JWT token for a user
func (s *RelationalJWTService) GenerateToken(user *models.User) (string, *TokenResponse, error) {
	now := time.Now()
	expiresAt := now.Add(s.expiresIn)

	claims := TokenClaims{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", nil, err
	}

	response := &TokenResponse{
		AccessToken: tokenString,
		TokenType:   "Bearer",
		ExpiresIn:   int64(s.expiresIn.Seconds()),
	}

	return tokenString, response, nil
}

// ValidateToken verifies a JWT token and returns the claims
func (s *RelationalJWTService) ValidateToken(tokenString string) (*TokenClaims, error) {
	claims := &TokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// ParseToken extracts claims from a token without validating signature
// Used during middleware processing where we trust the token structure
func (s *RelationalJWTService) ParseToken(tokenString string) (*TokenClaims, error) {
	claims := &TokenClaims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	return claims, nil
}
