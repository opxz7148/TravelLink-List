package server

import (
	"net/http"
	"os"
	"time"

	"tll-backend/internal/controllers"
	"tll-backend/internal/middleware"
	"tll-backend/internal/repositories"
	"tll-backend/internal/services"

	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	// Get allowed origins from environment, with sensible defaults for development
	allowedOrigins := []string{
		"http://localhost:5173",
		"http://localhost:5174",
		"http://localhost:3000",
		"http://127.0.0.1:5173",
		"http://127.0.0.1:5174",
		"http://127.0.0.1:3000",
	}
	
	// Add production origin if specified
	if frontendURL := os.Getenv("FRONTEND_URL"); frontendURL != "" {
		allowedOrigins = append(allowedOrigins, frontendURL)
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-Requested-With"},
		AllowCredentials: true, // Enable cookies/auth
		MaxAge:           12 * 60 * 60, // 12 hours
	}))

	// Initialize repositories
	userRepo := repositories.NewRelationalUserRepository(s.db)
	planRepo := repositories.NewRelationalPlanRepository(s.db)
	nodeRepo := repositories.NewRelationalNodeRepository(s.db)
	commentRepo := repositories.NewRelationalCommentRepository(s.db)
	ratingRepo := repositories.NewRelationalRatingRepository(s.db)
	promotionRepo := repositories.NewRelationalPromotionRequestRepository(s.db)

	// Initialize services
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key" // Fallback for testing
	}
	jwtService := services.NewRelationalJWTService(jwtSecret, 1*time.Hour)
	userService := services.NewRelationalUserService(userRepo, jwtService)
	planService := services.NewRelationalPlanService(planRepo, nodeRepo)
	nodeService := services.NewRelationalNodeService(nodeRepo)
	commentService := services.NewRelationalCommentService(commentRepo, planRepo)
	ratingService := services.NewRelationalRatingService(ratingRepo, planRepo)
	promotionService := services.NewRelationalPromotionService(promotionRepo, userRepo, planRepo)

	// Initialize controllers
	authController := controllers.NewAuthController(userService, jwtService)
	userController := controllers.NewUserController(userService)
	planController := controllers.NewPlanController(planService)
	nodeController := controllers.NewNodeController(nodeService)
	commentController := controllers.NewCommentController(planService, commentService)
	ratingController := controllers.NewRatingController(planService, ratingService)
	adminController := controllers.NewAdminController(planService, nodeService, userService, promotionService)
	promotionController := controllers.NewPromotionController(promotionService, userService)

	// Health & info endpoints
	r.GET("/", s.HelloWorldHandler)
	r.GET("/health", s.healthHandler)

	// Swagger UI endpoint
	r.GET("/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))

	// ============================================================================
	// USER-RELATED ROUTES
	// ============================================================================

	// Authentication routes (public, no auth required)
	authGroup := r.Group("/api/v1/auth")
	{
		authGroup.POST("/register", authController.Register)
		authGroup.POST("/login", authController.Login)
		authGroup.POST("/logout", middleware.RequireAuth(jwtService), authController.Logout)
	}

	// User profile routes
	userGroup := r.Group("/api/v1/users")
	{
		// Public: Get any user's profile
		userGroup.GET("/:id", userController.GetProfile)

		// Protected: Update own/admin profile
		userGroup.PUT("/:id", middleware.RequireAuth(jwtService), userController.UpdateProfile)

		// Protected: Change password (only own)
		userGroup.POST("/:id/change-password", middleware.RequireAuth(jwtService), userController.ChangePassword)
	}

	// Admin user management routes
	adminUserGroup := r.Group("/api/v1/users", middleware.RequireAuth(jwtService), middleware.RequireAdmin())
	{
		// Admin: Change user role
		adminUserGroup.PATCH("/:id/role", adminController.UpdateUserRole)

		// Admin: Deactivate user
		adminUserGroup.PATCH("/:id/deactivate", adminController.DeactivateUser)
	}

	// ============================================================================
	// NODE ROUTES (Attractions & Transitions)
	// ============================================================================

	// Public node discovery
	nodeGroup := r.Group("/api/v1/nodes")
	{
		// List approved nodes (with optional type filter)
		nodeGroup.GET("", nodeController.ListNodes)

		// Get single node details
		nodeGroup.GET("/:id", nodeController.GetNodeDetail)
	}

	// Protected node creation (traveller+)
	nodeCreationGroup := r.Group("/api/v1/nodes", middleware.RequireAuth(jwtService))
	{
		// Create new attraction node (pending admin approval)
		nodeCreationGroup.POST("/attraction", nodeController.CreateAttractionNode)

		// Create new transition node (pending admin approval)
		nodeCreationGroup.POST("/transition", nodeController.CreateTransitionNode)
	}

	// ============================================================================
	// TRAVEL PLAN ROUTES
	// ============================================================================

	// Public plan browsing
	planBrowseGroup := r.Group("/api/v1/plans")
	{
		// List published plans with pagination
		planBrowseGroup.GET("", planController.BrowsePlans)

		// Search plans by query
		planBrowseGroup.GET("/search", planController.SearchPlans)

		// Get plan details with linked list of nodes
		planBrowseGroup.GET("/:id", planController.GetPlanDetails)

		// Get rating statistics (public)
		planBrowseGroup.GET("/:id/ratings", ratingController.GetPlanRatingStats)

		// List comments (public)
		planBrowseGroup.GET("/:id/comments", commentController.GetComments)
	}

	// Protected plan operations (authenticated users)
	planProtectedGroup := r.Group("/api/v1/plans", middleware.RequireAuth(jwtService))
	{
		// Create draft plan (traveller+)
		planProtectedGroup.POST("", planController.CreatePlan)

		// Publish plan (plan owner or admin)
		planProtectedGroup.PATCH("/:id/publish", planController.PublishPlan)

		// Add/reorder/remove nodes from plan (plan owner or admin)
		planProtectedGroup.PATCH("/:id/nodes", planController.UpdatePlanNodes)

		// Submit comment (any authenticated user)
		planProtectedGroup.POST("/:id/comments", commentController.CreateComment)

		// Submit/update rating (any authenticated user)
		planProtectedGroup.POST("/:id/ratings", ratingController.SubmitRating)
		planProtectedGroup.PUT("/:id/ratings", ratingController.UpdateRating)

		// Get user's own rating (authenticated)
		planProtectedGroup.GET("/:id/my-rating", ratingController.GetUserRating)
	}

	// Protected comment operations
	commentGroup := r.Group("/api/v1/comments", middleware.RequireAuth(jwtService))
	{
		// Update own comment (author or admin)
		commentGroup.PUT("/:commentId", commentController.UpdateComment)

		// Delete own comment (author or admin)
		commentGroup.DELETE("/:commentId", commentController.DeleteComment)
	}

	// ============================================================================
	// PROMOTION REQUEST ROUTES
	// ============================================================================

	// Protected promotion request routes
	promotionGroup := r.Group("/api/v1/promotions", middleware.RequireAuth(jwtService))
	{
		// Submit promotion request (for role upgrade or plan promotion)
		promotionGroup.POST("/request", promotionController.SubmitRequest)

		// Get user's promotion requests
		promotionGroup.GET("/my-requests", promotionController.ListMyRequests)

		// Get specific promotion request status
		promotionGroup.GET("/request/:id", promotionController.GetRequestStatus)
	}

	// ============================================================================
	// ADMIN MODERATION ROUTES
	// ============================================================================

	// Admin-only moderation routes
	adminModGroup := r.Group("/api/v1/admin", middleware.RequireAuth(jwtService), middleware.RequireAdmin())
	{
		// Plan moderation
		adminModGroup.PATCH("/plans/:id/suspend", adminController.SuspendPlan)
		adminModGroup.DELETE("/plans/:id", adminController.DeletePlan)

		// Node approval/disapproval
		adminModGroup.PATCH("/nodes/:id/approve", adminController.ApproveNode)
		adminModGroup.PATCH("/nodes/:id/disapprove", adminController.DisapproveNode)

		// Node deletion
		adminModGroup.DELETE("/nodes/:id", adminController.DeleteNode)

		// Promotion request management
		adminModGroup.PATCH("/promotions/:id/approve", adminController.ApprovePromotionRequest)
		adminModGroup.PATCH("/promotions/:id/reject", adminController.RejectPromotionRequest)
	}

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
