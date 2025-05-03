package controllers

import (
	"github.com/form3tech-oss/jwt-go"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gosimple/internal/db"
	"gosimple/internal/models"
)

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func generateToken(id uint, role string) (string, error) {
	claims := jwt.MapClaims{
		"id":   id,
		"role": role,
		"exp":  time.Now().Add(30 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	return token.SignedString([]byte(secret))
}

func RegisterUser(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Неверные данные"})
		return
	}

	var existingUser models.User
	// Ищем пользователя, включая soft-deleted
	if err := db.DB.Unscoped().Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		if existingUser.DeletedAt.Valid {
			// Восстанавливаем пользователя
			existingUser.Name = input.Name
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка сервера"})
				return
			}
			existingUser.Password = string(hashedPassword)
			existingUser.DeletedAt.Time = time.Time{}
			existingUser.DeletedAt.Valid = false

			if err := db.DB.Unscoped().Save(&existingUser).Error; err != nil {
				log.Println("Ошибка при восстановлении пользователя:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка восстановления пользователя"})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Пользователь уже существует"})
			return
		}
	} else {
		// Создаём нового пользователя
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка сервера"})
			return
		}

		existingUser = models.User{
			Name:     input.Name,
			Email:    input.Email,
			Password: string(hashedPassword),
			Role:     "user",
		}

		if err := db.DB.Create(&existingUser).Error; err != nil {
			log.Println("Ошибка при создании пользователя:", err)
			c.JSON(http.StatusBadRequest, gin.H{"message": "Ошибка регистрации"})
			return
		}
	}

	token, err := generateToken(existingUser.ID, existingUser.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка генерации токена"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"_id":   existingUser.ID,
		"name":  existingUser.Name,
		"email": existingUser.Email,
		"token": token,
	})
}

func LoginUser(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Неверные данные"})
		return
	}

	var user models.User
	if err := db.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Неверные учетные данные"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Неверные учетные данные"})
		return
	}

	token, err := generateToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка генерации токена"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"_id":   user.ID,
		"name":  user.Name,
		"email": user.Email,
		"token": token,
	})
}
