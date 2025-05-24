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
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	
	// Frontend URL
	FrontendURL      string
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
	
	// Frontend URL
	if frontendURL := os.Getenv("FRONTEND_URL"); frontendURL != "" {
		config.FrontendURL = frontendURL
	}
	
	return config
} 