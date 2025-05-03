package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Internal server error: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Внутренняя ошибка сервера",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
