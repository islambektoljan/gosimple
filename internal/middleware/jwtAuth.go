package middleware

import (
	"github.com/form3tech-oss/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Токен не предоставлен"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		secret := os.Getenv("JWT_SECRET")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			log.Printf("Invalid token: %v\n", err) // Логирование ошибки
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный токен"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные claims"})
			c.Abort()
			return
		}

		// Проверка истечения токена
		if exp, ok := claims["exp"].(float64); ok && exp < float64(time.Now().Unix()) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Токен истёк"})
			c.Abort()
			return
		}

		// Приведение userID
		idFloat, ok := claims["id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный формат ID пользователя"})
			c.Abort()
			return
		}
		userID := uint(idFloat)
		c.Set("userID", userID)

		c.Next()

	}
}
