package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"resume.in/backend/config"
	"resume.in/backend/controllers"
	_ "resume.in/backend/docs" // Import generated docs
	"resume.in/backend/models"
	"resume.in/backend/routes"
	"resume.in/backend/utils"
)

// @title Resume.in API
// @version 1.0
// @description API Server for Resume.in application with chatbot capabilities and OAuth authentication
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.resume.in/support
// @contact.email support@resume.in

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
// @schemes http https

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description JWT Authorization header using the Bearer scheme. Example: "Bearer {token}"

func main() {
	// Initialize loggers
	utils.InitLoggers()

	// Load environment variables from .env file
	if err := utils.LoadEnv(); err != nil {
		utils.Error("Failed to load environment variables: %v", err)
		// Continue anyway, as we may have environment variables set directly
	}

	// Load configuration
	cfg := config.LoadConfigFromEnv()

	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to database with retry
	var resumeRepo models.ResumeRepository
	var chatbotRepo models.ChatbotRepository
	var userRepo models.UserRepository
	var maxRetries = 5
	var retryDelay = 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		db, err := config.ConnectDB()
		if err != nil {
			utils.Error("Failed to connect to database: %v. Retrying in %v... (%d/%d)", 
				err, retryDelay, i+1, maxRetries)
			time.Sleep(retryDelay)
			continue
		}

		// Run database migrations
		dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
			cfg.PostgresUser,
			cfg.PostgresPassword,
			cfg.PostgresHost,
			cfg.PostgresPort,
			cfg.PostgresDB,
			cfg.PostgresSSLMode,
		)
		if err := utils.RunMigrationsWithFS(dbURL, migrations); err != nil {
			utils.Error("Failed to run migrations: %v", err)
			os.Exit(1)
		}

		// Setup PostgreSQL repository for resume
		postgresRepo, err := models.NewPostgresResumeRepository(db)
		if err != nil {
			utils.Error("Failed to initialize database repository: %v", err)
			os.Exit(1)
		}

		// Initialize sample data
		if err := postgresRepo.InitDemoData(); err != nil {
			utils.Error("Failed to initialize demo data: %v", err)
		}

		resumeRepo = postgresRepo

		// Setup PostgreSQL repository for chatbot
		// Use simplified implementation to avoid LangChain dependency issues
		repo, err := models.NewPostgresChatbotRepository(db)
		if err != nil {
			utils.Error("Failed to initialize chatbot repository: %v", err)
			// Continue with other features
		} else {
			chatbotRepo = repo
			utils.Info("Chatbot repository initialized")
		}

		// Setup PostgreSQL repository for users
		userRepo = models.NewPostgresUserRepository(db)
		utils.Info("User repository initialized")

		break
	}

	// If we couldn't connect to the database, fall back to in-memory repository
	if resumeRepo == nil {
		utils.Info("Falling back to in-memory repository")
		
		// Initialize repository
		memoryRepo := models.NewInMemoryResumeRepository()
		memoryRepo.InitDemoData()
		resumeRepo = memoryRepo
		
		// For users, we need a database connection
		utils.Error("Authentication features require a database connection")
		os.Exit(1)
	}

	// Initialize controllers
	authController := controllers.NewAuthController(cfg, userRepo)
	resumeController := controllers.NewResumeController(resumeRepo)
	
	// Initialize chatbot controller if repository is available
	var chatbotController *controllers.ChatbotController
	if chatbotRepo != nil {
		chatbotController = controllers.NewChatbotController(chatbotRepo, resumeRepo)
		utils.Info("Chatbot controller initialized")
	} else {
		utils.Warning("Chatbot features will be disabled")
	}

	// Setup router
	router := routes.SetupRouter(cfg, authController, chatbotController, resumeController)

	// Remove the Swagger setup from here as it's now in routes.go
	utils.Info("Swagger UI available at http://localhost:%d/swagger/index.html", cfg.ServerPort)

	// Start server
	serverAddr := fmt.Sprintf(":%d", cfg.ServerPort)
	utils.Info("Starting server on http://localhost%s", serverAddr)
	
	if err := router.Run(serverAddr); err != nil {
		utils.Error("Failed to start server: %v", err)
	}
}
