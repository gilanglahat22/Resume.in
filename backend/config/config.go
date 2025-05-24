package config

import (
	"os"
	"strconv"
)

// Config stores all configuration settings
type Config struct {
	ServerPort       int
	Environment      string
	AllowOrigins     string
	LogLevel         string
	OpenRouterAPIKey string
	OpenRouterModel  string
	
	// JWT configuration
	JWTSecret        string
	
	// OAuth configuration
	GoogleClientID           string
	GoogleClientSecret       string
	GoogleRedirectURL        string
	GoogleRegisterRedirectURL string
	
	// Frontend URL
	FrontendURL      string

	// Database configuration
	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresSSLMode  string
}

// NewConfig creates and returns a new Config with default values
func NewConfig() *Config {
	return &Config{
		ServerPort:       8080,
		Environment:      "development",
		AllowOrigins:     "*",
		LogLevel:         "debug",
		OpenRouterAPIKey: "",
		OpenRouterModel:  "anthropic/claude-3-opus:beta", // Default to a powerful model
		JWTSecret:        "your-secret-key-change-in-production",
		FrontendURL:      "http://localhost:3000",
		GoogleRedirectURL: "http://localhost:8080/api/auth/google/callback",
		GoogleRegisterRedirectURL: "http://localhost:8080/api/auth/google/register/callback",
		
		// Default database configuration
		PostgresHost:     "localhost",
		PostgresPort:     5432,
		PostgresUser:     "resumeuser",
		PostgresPassword: "resumepassword",
		PostgresDB:       "resumedb",
		PostgresSSLMode:  "disable",
	}
}

// LoadConfigFromEnv loads configuration from environment variables
func LoadConfigFromEnv() *Config {
	config := NewConfig()
	
	// Server port
	if port := os.Getenv("SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.ServerPort = p
		}
	}
	
	// Environment
	if env := os.Getenv("ENVIRONMENT"); env != "" {
		config.Environment = env
	}
	
	// CORS allowed origins
	if origins := os.Getenv("ALLOW_ORIGINS"); origins != "" {
		config.AllowOrigins = origins
	}
	
	// Log level
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		config.LogLevel = logLevel
	}
	
	// Open Router API key
	if apiKey := os.Getenv("OPEN_ROUTER_API_KEY"); apiKey != "" {
		config.OpenRouterAPIKey = apiKey
	}
	
	// Open Router model
	if model := os.Getenv("OPEN_ROUTER_MODEL"); model != "" {
		config.OpenRouterModel = model
	}
	
	// JWT Secret
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		config.JWTSecret = jwtSecret
	}
	
	// Google OAuth
	if clientID := os.Getenv("GOOGLE_CLIENT_ID"); clientID != "" {
		config.GoogleClientID = clientID
	}
	
	if clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET"); clientSecret != "" {
		config.GoogleClientSecret = clientSecret
	}
	
	if redirectURL := os.Getenv("GOOGLE_REDIRECT_URL"); redirectURL != "" {
		config.GoogleRedirectURL = redirectURL
	}

	if registerRedirectURL := os.Getenv("GOOGLE_REGISTER_REDIRECT_URL"); registerRedirectURL != "" {
		config.GoogleRegisterRedirectURL = registerRedirectURL
	}
	
	// Frontend URL
	if frontendURL := os.Getenv("FRONTEND_URL"); frontendURL != "" {
		config.FrontendURL = frontendURL
	}

	// Database configuration
	if host := os.Getenv("POSTGRES_HOST"); host != "" {
		config.PostgresHost = host
	}
	
	if port := os.Getenv("POSTGRES_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.PostgresPort = p
		}
	}
	
	if user := os.Getenv("POSTGRES_USER"); user != "" {
		config.PostgresUser = user
	}
	
	if password := os.Getenv("POSTGRES_PASSWORD"); password != "" {
		config.PostgresPassword = password
	}
	
	if db := os.Getenv("POSTGRES_DB"); db != "" {
		config.PostgresDB = db
	}
	
	if sslMode := os.Getenv("POSTGRES_SSLMODE"); sslMode != "" {
		config.PostgresSSLMode = sslMode
	}
	
	return config
} 