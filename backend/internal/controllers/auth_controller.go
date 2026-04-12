package controllers

import (
	"strings"

	"tll-backend/internal/middleware"
	"tll-backend/internal/models"
	"tll-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// AuthController handles user registration and authentication
type AuthController struct {
	userService services.UserService
	jwtService  services.JWTService
}

// NewAuthController creates a new auth controller
func NewAuthController(userService services.UserService, jwtService services.JWTService) *AuthController {
	return &AuthController{
		userService: userService,
		jwtService:  jwtService,
	}
}

// RegisterRequest represents user registration request payload
type RegisterRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Username    string `json:"username" binding:"required,min=3,max=50"`
	Password    string `json:"password" binding:"required,min=8"`
	DisplayName string `json:"display_name" binding:"max=100"`
}

// LoginRequest represents user login request payload
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterResponse represents user registration response
type RegisterResponse struct {
	User        *UserResponse `json:"user"`
	AccessToken string        `json:"access_token"`
	TokenType   string        `json:"token_type"`
	ExpiresIn   int           `json:"expires_in"`
}

// LoginResponse represents user login response
type LoginResponse struct {
	User        *UserResponse `json:"user"`
	AccessToken string        `json:"access_token"`
	TokenType   string        `json:"token_type"`
	ExpiresIn   int           `json:"expires_in"`
}

// UserResponse represents user data in API response
type UserResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Role        string `json:"role"`
	CreatedAt   string `json:"created_at"`
	LastLoginAt string `json:"last_login_at,omitempty"`
}

// Register handles POST /api/v1/auth/register
func (ac *AuthController) Register(c *gin.Context) {
	var req RegisterRequest

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ValidationErrorResponse(c, "invalid request", nil)
		return
	}

	// Trim whitespace
	req.Email = strings.TrimSpace(req.Email)
	req.Username = strings.TrimSpace(req.Username)
	req.DisplayName = strings.TrimSpace(req.DisplayName)

	// Call service to register user
	user, tokenResp, err := ac.userService.Register(c.Request.Context(), req.Email, req.Username, req.Password, req.DisplayName)
	if err != nil {
		if err == models.ErrDuplicateEmail {
			middleware.ValidationErrorResponse(c, "Email already registered", gin.H{"field": "email"})
			return
		}
		if err == models.ErrDuplicateUsername {
			middleware.ValidationErrorResponse(c, "Username already taken", gin.H{"field": "username"})
			return
		}
		if err == models.ErrValidation {
			middleware.ValidationErrorResponse(c, "Password does not meet security requirements", gin.H{"field": "password"})
			return
		}
		middleware.InternalErrorResponse(c, "Failed to register user")
		return
	}

	// Build response
	resp := RegisterResponse{
		User: &UserResponse{
			ID:          user.ID,
			Email:       user.Email,
			Username:    user.Username,
			DisplayName: user.DisplayName,
			Role:        user.Role,
			CreatedAt:   user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		},
		AccessToken: tokenResp.AccessToken,
		TokenType:   tokenResp.TokenType,
		ExpiresIn:   int(tokenResp.ExpiresIn),
	}

	middleware.SuccessResponse(c, 201, resp)
}

// Login handles POST /api/v1/auth/login
func (ac *AuthController) Login(c *gin.Context) {
	var req LoginRequest

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ValidationErrorResponse(c, "invalid request", nil)
		return
	}

	// Trim whitespace
	req.Email = strings.TrimSpace(req.Email)

	// Call service to authenticate user
	user, tokenResp, err := ac.userService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if err == models.ErrInvalidCredentials {
			middleware.AuthErrorResponse(c, "Email or password is incorrect")
			return
		}
		if err == models.ErrUserInactive {
			middleware.AuthErrorResponse(c, "User account is inactive")
			return
		}
		middleware.InternalErrorResponse(c, "Failed to login")
		return
	}

	// Build response
	lastLoginAt := ""
	if user.LastLoginAt != nil {
		lastLoginAt = user.LastLoginAt.Format("2006-01-02T15:04:05Z")
	}

	resp := LoginResponse{
		User: &UserResponse{
			ID:          user.ID,
			Email:       user.Email,
			Username:    user.Username,
			DisplayName: user.DisplayName,
			Role:        user.Role,
			CreatedAt:   user.CreatedAt.Format("2006-01-02T15:04:05Z"),
			LastLoginAt: lastLoginAt,
		},
		AccessToken: tokenResp.AccessToken,
		TokenType:   tokenResp.TokenType,
		ExpiresIn:   int(tokenResp.ExpiresIn),
	}

	middleware.SuccessResponse(c, 200, resp)
}
