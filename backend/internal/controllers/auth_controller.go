package controllers

import (
	"regexp"
	"strings"

	"tll-backend/internal/logger"
	"tll-backend/internal/middleware"
	"tll-backend/internal/models"
	"tll-backend/internal/services"
	"tll-backend/internal/utilities"

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

// validatePassword checks if password meets security requirements:
// - At least 8 characters
// - At least one uppercase letter (A-Z)
// - At least one digit (0-9)
// - At least one special character (!@#$%^&*)
func validatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	
	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*]`).MatchString(password)
	
	return hasUppercase && hasDigit && hasSpecial
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
// @Summary Register a new user
// @Description Create a new user account with email, username, password, and display name
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration request"
// @Success 201 {object} middleware.ResponseEnvelope "User successfully registered"
// @Failure 400 {object} middleware.ResponseEnvelope "Validation error"
// @Failure 500 {object} middleware.ResponseEnvelope "Internal server error"
// @Router /auth/register [post]
func (ac *AuthController) Register(c *gin.Context) {
	var req RegisterRequest

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ValidationErrorResponse(c, "Invalid request format", nil)
		return
	}

	// Trim whitespace
	req.Email = strings.TrimSpace(req.Email)
	req.Username = strings.TrimSpace(req.Username)
	req.DisplayName = strings.TrimSpace(req.DisplayName)

	// Validate password meets security requirements
	if !validatePassword(req.Password) {
		middleware.ValidationErrorResponse(c, "Password does not meet security requirements (min 8 chars, must include uppercase, digit, and special character)", gin.H{"field": "password"})
		return
	}

	// Call service to register user
	log := logger.GetLogger("AuthController")
	user, tokenResp, err := ac.userService.Register(c.Request.Context(), req.Email, req.Username, req.Password, req.DisplayName)
	if err != nil {
		if err == models.ErrDuplicateEmail {
			log.ValidationError(c.Request.Context(), map[string]string{"email": "Email already registered"})
			middleware.ValidationErrorResponse(c, "Email already registered", gin.H{"field": "email"})
			return
		}
		if err == models.ErrDuplicateUsername {
			log.ValidationError(c.Request.Context(), map[string]string{"username": "Username already taken"})
			middleware.ValidationErrorResponse(c, "Username already taken", gin.H{"field": "username"})
			return
		}
		if err == models.ErrValidation {
			log.ValidationError(c.Request.Context(), map[string]string{"password": "Does not meet security requirements"})
			middleware.ValidationErrorResponse(c, "Password does not meet security requirements", gin.H{"field": "password"})
			return
		}
		log.ServiceError(c.Request.Context(), "UserService", "Register", err, "")
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

	log.UserRegistration(c.Request.Context(), user.ID, user.Email, user.Username)
	middleware.SuccessResponse(c, 201, resp)
}

// Login handles POST /api/v1/auth/login
// @Summary Authenticate user and get access token
// @Description Authenticate with email and password to receive JWT access token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login request with email and password"
// @Success 200 {object} middleware.ResponseEnvelope "User successfully authenticated"
// @Failure 400 {object} middleware.ResponseEnvelope "Validation error"
// @Failure 401 {object} middleware.ResponseEnvelope "Invalid credentials or inactive user"
// @Failure 500 {object} middleware.ResponseEnvelope "Internal server error"
// @Router /auth/login [post]
func (ac *AuthController) Login(c *gin.Context) {
	var req LoginRequest

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ValidationErrorResponse(c, "Invalid request format", nil)
		return
	}

	// Trim whitespace
	req.Email = strings.TrimSpace(req.Email)

	// Call service to authenticate user
	log := logger.GetLogger("AuthController")
	user, tokenResp, err := ac.userService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if err == models.ErrInvalidCredentials {
			log.AuthenticationError(c.Request.Context(), req.Email, "Invalid credentials provided")
			middleware.AuthErrorResponseWithCode(c, "INVALID_CREDENTIALS", "Email or password is incorrect")
			return
		}
		if err == models.ErrUserInactive {
			log.AuthenticationError(c.Request.Context(), req.Email, "Account inactive")
			middleware.AuthErrorResponseWithCode(c, "ACCOUNT_INACTIVE", "Your account has been deactivated by an administrator")
			return
		}
		if err == models.ErrNotFound {
			log.AuthenticationError(c.Request.Context(), req.Email, "User not found")
			middleware.AuthErrorResponseWithCode(c, "INVALID_CREDENTIALS", "Email or password is incorrect")
			return
		}
		log.ServiceError(c.Request.Context(), "UserService", "Login", err, "")
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

	log.UserAuthentication(c.Request.Context(), user.ID, user.Email)
	middleware.SuccessResponse(c, 200, resp)
}

// Logout handles POST /api/v1/auth/logout
// @Summary Logout user
// @Description Logout authenticated user. In a stateless JWT system, the client should discard the token.
// @Tags auth
// @Security Bearer
// @Success 204 "User successfully logged out"
// @Failure 401 {object} middleware.ResponseEnvelope "User not authenticated"
// @Router /auth/logout [post]
func (ac *AuthController) Logout(c *gin.Context) {
	// Get authenticated user (just to verify auth)
	userID, ok := utilities.GetUserIDFromContext(c)
	if !ok {
		log := logger.GetLogger("AuthController")
		log.AuthorizationError(c.Request.Context(), "", "Logout attempted without authentication")
		middleware.AuthErrorResponse(c, "User not authenticated")
		return
	}

	// Note: In a stateless JWT system, logout is client-side (discard token)
	// If implementing token blacklisting, add logic here
	// For now, just return 204 No Content to signal success
	log := logger.GetLogger("AuthController")
	log.UserLogout(c.Request.Context(), userID)
	c.Status(204) // No Content response
}
