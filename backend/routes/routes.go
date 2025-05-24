package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"resume.in/backend/config"
	"resume.in/backend/controllers"
	"resume.in/backend/docs"
	"resume.in/backend/middleware"
)

// SetupRouter configures all API routes
func SetupRouter(
	cfg *config.Config,
	authController *controllers.AuthController,
	resumeController *controllers.ResumeController,
	chatbotController *controllers.ChatbotController,
) *gin.Engine {
	router := gin.Default()

	// Enable CORS middleware
	router.Use(middleware.CORSMiddleware())

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
			auth.GET("/google/login", authController.GoogleLogin)
			auth.GET("/google/callback", authController.GoogleCallback)
			auth.POST("/refresh", authController.RefreshToken)
		}

		// Protected authentication endpoints
		authProtected := auth.Group("")
		authProtected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			authProtected.POST("/logout", authController.Logout)
			authProtected.GET("/profile", authController.GetProfile)
		}

		// Resume endpoints (protected)
		resumesProtected := api.Group("/resumes")
		resumesProtected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			resumesProtected.GET("", resumeController.GetResumes)
			resumesProtected.GET("/:id", resumeController.GetResume)
			resumesProtected.POST("", resumeController.CreateResume)
			resumesProtected.PUT("/:id", resumeController.UpdateResume)
			resumesProtected.DELETE("/:id", resumeController.DeleteResume)
		}

		// Chatbot endpoints (protected)
		if chatbotController != nil {
			chatProtected := api.Group("/chat")
			chatProtected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
			{
				chatProtected.POST("/message", chatbotController.SendMessage)
				chatProtected.GET("/history/:sessionId", chatbotController.GetChatHistory)
				chatProtected.POST("/document", chatbotController.UploadDocument)
				chatProtected.POST("/generate-resume", chatbotController.GenerateATSResume)
			}
		}
	}

	// Swagger documentation - handle with a single route
	router.GET("/swagger/*any", func(c *gin.Context) {
		// Get the path
		path := c.Param("any")
		
		// Check if this is a request for doc.json
		if path == "/doc.json" || path == "doc.json" {
			c.Header("Content-Type", "application/json")
			c.String(http.StatusOK, docs.DocJSON)
			return
		}
		
		// For all other paths, use the standard Swagger UI handler
		ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
	})

	return router
} 