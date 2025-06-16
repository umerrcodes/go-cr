package handler

import (
	apiHandler "dummy-backend/internal/handler"
	"dummy-backend/internal/repository"
	"dummy-backend/internal/service"
	"dummy-backend/pkg/config"
	"dummy-backend/pkg/database"
	"dummy-backend/pkg/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var router *gin.Engine

func init() {
	// Load environment variables
	godotenv.Load()

	// Load configuration
	cfg := config.LoadConfig()

	// Set gin mode to release for production
	gin.SetMode(gin.ReleaseMode)

	// Initialize database
	db := database.NewPostgresDB(cfg.DatabaseDSN)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	taskService := service.NewTaskService(taskRepo)

	// Initialize handlers
	authHandler := apiHandler.NewAuthHandler(authService)
	taskHandler := apiHandler.NewTaskHandler(taskService)

	// Initialize router
	router = gin.New()
	router.Use(gin.Recovery())

	// Apply middleware
	router.Use(middleware.CORSMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		if err := database.HealthCheck(db); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// API routes
	api := router.Group("/api")
	{
		// Auth routes (no authentication required)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Task routes (authentication required)
		tasks := api.Group("/tasks")
		tasks.Use(middleware.AuthMiddleware(authService))
		{
			tasks.POST("", taskHandler.CreateTask)
			tasks.GET("", taskHandler.GetAllTasks)
			tasks.GET("/:id", taskHandler.GetTaskByID)
			tasks.PUT("/:id", taskHandler.UpdateTask)
			tasks.DELETE("/:id", taskHandler.DeleteTask)
		}
	}
}

// Handler is the exported function that Vercel will call
func Handler(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}
