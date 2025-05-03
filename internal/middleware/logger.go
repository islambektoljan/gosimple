package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		log.Printf("%s - [%s] %s %s %d %s",
			c.ClientIP(),
			start.Format(time.RFC1123),
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			duration,
		)
	}
}
