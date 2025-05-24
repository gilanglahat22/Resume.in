package utils

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/google/uuid"
)

// GenerateUUID generates a random UUID string
func GenerateUUID() string {
	return uuid.New().String()
}

// GenerateRandomString generates a random string of specified length
func GenerateRandomString(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		// Fallback to UUID if random generation fails
		return GenerateUUID()
	}
	return base64.URLEncoding.EncodeToString(b)[:length]
} 