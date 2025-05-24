package routes

import (
	"github.com/gin-gonic/gin"
	"resume.in/backend/controllers"
)

// SetupAuthRoutes sets up all auth related routes
func SetupAuthRoutes(router *gin.RouterGroup, authController *controllers.AuthController) {
	auth := router.Group("/auth")
	{
		// Registration
		auth.POST("/register", authController.Register)
		auth.GET("/google/register", authController.GoogleRegister)
		auth.GET("/google/register/callback", authController.GoogleRegisterCallback)

		// Google OAuth
		auth.GET("/google/login", authController.GoogleLogin)
		auth.GET("/google/callback", authController.GoogleCallback)

		// Token management
		auth.POST("/refresh", authController.RefreshToken)
		auth.POST("/logout", authController.Logout)

		// Protected routes
		auth.GET("/profile", authController.GetProfile)
	}
} 