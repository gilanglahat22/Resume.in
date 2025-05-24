package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "resume.in/backend/docs" // Import generated docs
	"resume.in/backend/config"
	"resume.in/backend/controllers"
	"resume.in/backend/middleware"
)

// SetupRouter configures all API routes
func SetupRouter(
	cfg *config.Config,
	authController *controllers.AuthController,
	chatbotController *controllers.ChatbotController,
	resumeController *controllers.ResumeController,
) *gin.Engine {
	router := gin.Default()

	// Enable CORS middleware
	router.Use(middleware.CORSMiddleware(cfg.AllowOrigins))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// API routes
	api := router.Group("/api")
	{
		// Authentication endpoints (public)
		auth := api.Group("/auth")
		{
			// Registration endpoints
			auth.POST("/register", authController.Register)
			auth.GET("/google/register", authController.GoogleRegister)
			
			// Google OAuth login endpoints (unified callback handles both login and registration)
			auth.GET("/google/login", authController.GoogleLogin)
			auth.GET("/google/callback", authController.GoogleCallback)
			
			// Token management
			auth.POST("/refresh", authController.RefreshToken)
		}

		// Protected authentication endpoints
		authProtected := auth.Group("")
		authProtected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			authProtected.POST("/logout", authController.Logout)
			authProtected.GET("/profile", authController.GetProfile)
		}

		// Chatbot endpoints (protected)
		if chatbotController != nil {
			chat := api.Group("/chat")
			chat.Use(middleware.AuthMiddleware(cfg.JWTSecret))
			{
				chat.POST("/message", chatbotController.SendMessage)
				chat.GET("/history/:sessionId", chatbotController.GetChatHistory)
				chat.POST("/document", chatbotController.UploadDocument)
				chat.POST("/generate-resume", chatbotController.GenerateATSResume)
			}
		}

		// Resume endpoints (protected)
		resume := api.Group("/resumes")
		resume.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			resume.GET("", resumeController.GetResumes)
			resume.GET("/:id", resumeController.GetResume)
			resume.POST("", resumeController.CreateResume)
			resume.PUT("/:id", resumeController.UpdateResume)
			resume.DELETE("/:id", resumeController.DeleteResume)
		}
	}

	// Swagger documentation - handle with a single route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
} 