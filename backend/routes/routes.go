package routes

import (
	"github.com/gin-gonic/gin"
	"resume.in/backend/controllers"
	"resume.in/backend/middleware"
)

// SetupRouter configures all routes for the application
func SetupRouter(resumeController *controllers.ResumeController, chatbotController *controllers.ChatbotController) *gin.Engine {
	router := gin.Default()
	
	// Apply CORS middleware
	router.Use(middleware.CORS())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "up",
		})
	})

	// API routes
	api := router.Group("/api")
	{
		// Resume endpoints
		resumeRoutes := api.Group("/resume")
		{
			resumeRoutes.GET("/", resumeController.GetAllResumes)
			resumeRoutes.GET("/:id", resumeController.GetResume)
			resumeRoutes.POST("/", resumeController.CreateResume)
			resumeRoutes.PUT("/:id", resumeController.UpdateResume)
			resumeRoutes.DELETE("/:id", resumeController.DeleteResume)
		}
		
		// Skills endpoint
		api.GET("/skills", resumeController.GetAllSkills)
		
		// Experience endpoint
		api.GET("/experience", resumeController.GetAllExperience)
		
		// Chatbot endpoints
		chatRoutes := api.Group("/chat")
		{
			chatRoutes.POST("/message", chatbotController.SendMessage)
			chatRoutes.GET("/history/:sessionId", chatbotController.GetChatHistory)
			chatRoutes.POST("/document", chatbotController.UploadDocument)
		}
	}

	return router
} 