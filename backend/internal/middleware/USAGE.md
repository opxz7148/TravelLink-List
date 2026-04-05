package middleware

// USAGE GUIDE FOR MIDDLEWARE
//
// The middleware package provides authentication, authorization, and response handling
// utilities for Gin-based HTTP handlers.
//
// 1. AUTHENTICATION MIDDLEWARE (RequireAuth)
// ============================================
// RequireAuth validates JWT tokens from Authorization header
//
// Usage in route registration:
//   r.POST("/api/v1/posts", middleware.RequireAuth(jwtService), handlers.CreatePost)
//
// In handler, extract user info:
//   func CreatePost(c *gin.Context) {
//       claims := middleware.GetUserClaims(c)
//       userID := claims.UserID
//       // ... create post
//   }
//
// 2. ROLE-BASED ACCESS CONTROL (RequireRole)
// ============================================
// RequireRole verifies user has required roles. Must follow RequireAuth.
//
// Usage:
//   r.DELETE("/api/v1/users/:id",
//       middleware.RequireAuth(jwtService),
//       middleware.RequireRole(models.RoleAdmin),
//       handlers.DeleteUser)
//
// Multiple allowed roles:
//   middleware.RequireRole(models.RoleAdmin, models.RoleTraveller)
//
// 3. RESPONSE HELPERS
// ===================
// Standardize all responses using middleware response functions
//
// Success response:
//   middleware.SuccessResponse(c, http.StatusOK, userData)
//   // Returns response envelope with success=true, data={}, error=null
//
// Error responses:
//   middleware.ValidationErrorResponse(c, "Email is required", nil)
//   middleware.AuthErrorResponse(c, "Invalid credentials")
//   middleware.ForbiddenErrorResponse(c, "Insufficient permissions")
//   middleware.NotFoundErrorResponse(c, "User not found")
//   middleware.InternalErrorResponse(c, "Database error")
//
// 4. ERROR HANDLING
// =================
// Convert service errors to HTTP responses:
//
//   user, tokenResp, err := userService.Register(ctx, email, username, password, displayName)
//   if middleware.HandleServiceError(c, err) {
//       return
//   }
//   middleware.SuccessResponse(c, http.StatusCreated, gin.H{
//       "user":         user,
//       "access_token": tokenResp.AccessToken,
//       "token_type":   tokenResp.TokenType,
//       "expires_in":   tokenResp.ExpiresIn,
//   })
//
// 5. CONTEXT HELPER
// =================
// Convenient access to user data in handlers:
//
//   helper := middleware.NewContextHelper()
//   userID := helper.GetUserID(c)
//   isAdmin := helper.IsAdmin(c)
//   role := helper.GetUserRole(c)
//   if helper.IsTraveller(c) {
//       // traveller-specific logic
//   }
//
// 6. RESOURCE OWNERSHIP VALIDATION
// =================================
// Verify user owns resource before modifying:
//
//   if !middleware.EnsureUserOwnsResource(c, postAuthorID) {
//       return
//   }
//   // Proceed with modification
//
// 7. FULL EXAMPLE HANDLER
// =======================
// Example: Update user profile (authenticated, owns resource)
//
//   r.PUT("/api/v1/users/:id",
//       middleware.RequireAuth(jwtService),
//       handlers.UpdateUserProfile)
//
//   func UpdateUserProfile(c *gin.Context) {
//       userID := c.Param("id")
//       
//       // Verify ownership
//       if !middleware.EnsureUserOwnsResource(c, userID) {
//           return
//       }
//
//       // Validate input
//       var req UpdateProfileRequest
//       if err := c.ShouldBindJSON(&req); err != nil {
//           middleware.ValidationErrorResponse(c, "Invalid request", nil)
//           return
//       }
//
//       // Call service
//       user, err := userService.UpdateProfile(c.Request.Context(), 
//           userID, req.DisplayName, req.Bio, req.ProfilePicture)
//       if middleware.HandleServiceError(c, err) {
//           return
//       }
//
//       // Return success response
//       middleware.SuccessResponse(c, http.StatusOK, user)
//   }
//
// 8. ADMIN-ONLY ENDPOINT
// ======================
// Example: Delete any user (admin only)
//
//   r.DELETE("/api/v1/users/:id",
//       middleware.RequireAuth(jwtService),
//       middleware.RequireRole(models.RoleAdmin),
//       handlers.DeleteUser)
//
//   func DeleteUser(c *gin.Context) {
//       userID := c.Param("id")
//       err := userService.DeleteUser(c.Request.Context(), userID)
//       if middleware.HandleServiceError(c, err) {
//           return
//       }
//       middleware.SuccessResponse(c, http.StatusOK, gin.H{
//           "message": "User deleted successfully",
//       })
//   }
