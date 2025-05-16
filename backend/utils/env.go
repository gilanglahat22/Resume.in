package utils

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// LoadEnv loads environment variables from a .env file
// It first looks in the current directory and then in the parent directory
func LoadEnv() error {
	// Try to load from the current directory first
	err := loadEnvFile(".env")
	if err != nil {
		// If not found in current directory, try the parent directory
		err = loadEnvFile(filepath.Join("..", ".env"))
		if err != nil {
			// Not a critical error, we'll just use environment variables that are already set
			Info("No .env file found, using existing environment variables")
			return nil
		}
	}
	return nil
}

// loadEnvFile reads a file and sets environment variables from it
func loadEnvFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Skip empty lines and comments
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		// Split by first equals sign
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // Skip invalid lines
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		if len(value) > 1 && (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
			(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
			value = value[1 : len(value)-1]
		}

		// Only set if not already set
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	Info("Loaded environment variables from %s", filename)
	return nil
} 