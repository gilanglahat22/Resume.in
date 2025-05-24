package config

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
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
func ConnectDB() (*sqlx.DB, error) {
	dbConfig := GetDBConfig()

	// Try to resolve the database host first to check network connectivity
	// This can help diagnose DNS issues in Docker networks
	_, err := net.LookupHost(dbConfig.Host)
	if err != nil {
		utils.Error("DNS lookup for database host '%s' failed: %v", dbConfig.Host, err)
		return nil, fmt.Errorf("DNS lookup for database host '%s' failed: %w", dbConfig.Host, err)
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s connect_timeout=10",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.SSLMode,
	)

	utils.Info("Attempting to connect to database at %s:%s", dbConfig.Host, dbConfig.Port)
	
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Set reasonable connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	
	// Ping with timeout
	err = ping(db, 10*time.Second)
	if err != nil {
		return nil, err
	}

	utils.Info("Successfully connected to PostgreSQL database at %s:%s", dbConfig.Host, dbConfig.Port)
	return db, nil
}

// ping attempts to ping the database with a timeout
func ping(db *sqlx.DB, timeout time.Duration) error {
	errc := make(chan error, 1)
	go func() {
		errc <- db.Ping()
	}()

	select {
	case err := <-errc:
		return err
	case <-time.After(timeout):
		return fmt.Errorf("database ping timed out after %v", timeout)
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
} 