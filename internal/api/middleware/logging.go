package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Log request details
		log.Printf(
			"[%s] %s %s %v %d %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			latency,
			c.Writer.Status(),
			c.Errors.String(),
		)
	}
}
