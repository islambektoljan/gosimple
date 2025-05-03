package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gosimple/internal/db"
	"gosimple/internal/models"
	"log"
	"net/http"
	"os"
	"strings"
)

func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Токен не предоставлен"})
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
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Неверный токен"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Ошибка в токене"})
			c.Abort()
			return
		}

		userIDFloat, ok := claims["id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Неверный ID в токене"})
			c.Abort()
			return
		}
		userID := uint(userIDFloat)

		var user models.User
		if err := db.DB.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Пользователь не найден"})
			c.Abort()
			return
		}

		if user.Role != requiredRole {
			log.Printf("User with ID %d tried to access admin route\n", userID) // Логирование попытки доступа
			c.JSON(http.StatusForbidden, gin.H{"message": "Недостаточно прав"})
			c.Abort()
			return
		}

		c.Set("userID", user.ID)
		c.Set("userRole", user.Role)
		c.Next()
	}
}
