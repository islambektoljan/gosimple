package middleware

import (
	"github.com/gin-gonic/gin"
	"gosimple/internal/db"
	"gosimple/internal/models"
	"net/http"
)

func BookOwnerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		bookID := c.Param("id")

		var book models.Book
		if err := db.DB.Where("id = ?", bookID).First(&book).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Книга не найдена"})
			c.Abort()
			return
		}

		userID := c.MustGet("userID").(uint)

		if book.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"message": "Вы не можете редактировать или удалять эту книгу"})
			c.Abort()
			return
		}

		c.Next()
	}
}
