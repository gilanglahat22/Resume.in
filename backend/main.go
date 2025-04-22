package main

import (
	"fmt"
	"os"
	"time"

	"resume.in/backend/config"
	"resume.in/backend/controllers"
	"resume.in/backend/models"
	"resume.in/backend/routes"
	"resume.in/backend/utils"
)

func main() {
	// Initialize loggers
	utils.InitLoggers()

	// Load configuration
	cfg := config.LoadConfigFromEnv()

	// Connect to database with retry
	var resumeRepo models.ResumeRepository
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

		// Setup PostgreSQL repository
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
		break
	}

	// If we couldn't connect to the database, fall back to in-memory repository
	if resumeRepo == nil {
		utils.Info("Falling back to in-memory repository")
		
		// Initialize repository
		memoryRepo := models.NewInMemoryResumeRepository()
		memoryRepo.InitDemoData()
		resumeRepo = memoryRepo
	}

	// Initialize controller
	resumeController := controllers.NewResumeController(resumeRepo)

	// Setup router
	router := routes.SetupRouter(resumeController)

	// Start server
	serverAddr := fmt.Sprintf(":%d", cfg.ServerPort)
	utils.Info("Starting server on http://localhost%s", serverAddr)
	
	if err := router.Run(serverAddr); err != nil {
		utils.Error("Failed to start server: %v", err)
	}
}
