package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"resume.in/backend/controllers"
	"resume.in/backend/docs"
	"resume.in/backend/middleware"
)

// SetupRouter configures all API routes
func SetupRouter(
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
		// Resume endpoints
		api.GET("/resumes", resumeController.GetResumes)
		api.GET("/resumes/:id", resumeController.GetResume)
		api.POST("/resumes", resumeController.CreateResume)
		api.PUT("/resumes/:id", resumeController.UpdateResume)
		api.DELETE("/resumes/:id", resumeController.DeleteResume)

		// Chatbot endpoints
		if chatbotController != nil {
			api.POST("/chat/message", chatbotController.SendMessage)
			api.GET("/chat/history/:sessionId", chatbotController.GetChatHistory)
			api.POST("/chat/document", chatbotController.UploadDocument)
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