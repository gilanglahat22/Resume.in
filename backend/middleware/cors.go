package middleware

import (
	"strings"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware adds Cross-Origin Resource Sharing headers to responses
func CORSMiddleware(allowOrigins string) gin.HandlerFunc {
	origins := strings.Split(allowOrigins, ",")
	
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		// Check if the request origin is allowed
		for _, allowedOrigin := range origins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}
		
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
} 