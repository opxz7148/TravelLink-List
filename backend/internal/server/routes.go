package server

import (
	"net/http"
	"os"
	"time"

	"tll-backend/internal/controllers"
	"tll-backend/internal/repositories"
	"tll-backend/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	// Initialize repositories
	userRepo := repositories.NewRelationalUserRepository(s.db)
	planRepo := repositories.NewRelationalPlanRepository(s.db)
	nodeRepo := repositories.NewRelationalNodeRepository(s.db)

	// Initialize services
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key" // Fallback for testing
	}
	jwtService := services.NewRelationalJWTService(jwtSecret, 1*time.Hour)
	userService := services.NewRelationalUserService(userRepo, jwtService)
	planService := services.NewRelationalPlanService(planRepo, nodeRepo)
	nodeService := services.NewRelationalNodeService(nodeRepo)

	// Initialize controllers (for testing - routes will be registered later)
	_ = controllers.NewAuthController(userService, jwtService)
	_ = controllers.NewPlanController(planService)
	_ = controllers.NewNodeController(nodeService)
	_ = controllers.NewCommentController(planService)
	_ = controllers.NewRatingController(planService)
	_ = controllers.NewAdminController(planService, nodeService, userService)

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

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
