package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"resume.in/backend/utils"
)

// DBConfig holds database configuration parameters
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// GetDBConfig loads database configuration from environment variables
func GetDBConfig() *DBConfig {
	return &DBConfig{
		Host:     getEnv("POSTGRES_HOST", "localhost"),
		Port:     getEnv("POSTGRES_PORT", "5432"),
		User:     getEnv("POSTGRES_USER", "resumeuser"),
		Password: getEnv("POSTGRES_PASSWORD", "resumepassword"),
		DBName:   getEnv("POSTGRES_DB", "resumedb"),
		SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
	}
}

// ConnectDB establishes a connection to the database
func ConnectDB() (*sql.DB, error) {
	dbConfig := GetDBConfig()

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	utils.Info("Connected to PostgreSQL database")
	return db, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
} 