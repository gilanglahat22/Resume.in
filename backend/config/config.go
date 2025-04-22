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
}

// NewConfig creates and returns a new Config with default values
func NewConfig() *Config {
	return &Config{
		ServerPort:       8080,
		Environment:      "development",
		AllowOrigins:     "*",
		LogLevel:         "debug",
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
	
	return config
} 