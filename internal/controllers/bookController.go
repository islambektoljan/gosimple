package controllers

import (
	"github.com/gin-gonic/gin"
	"gosimple/internal/db"
	"gosimple/internal/models"
	"gosimple/internal/utils"
	"net/http"
)

func CreateBook(c *gin.Context) {
	var input models.Book
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Неверные данные"})
		return
	}

	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	input.UserID = userID

	// Создаём книгу
	if err := db.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при создании книги"})
		return
	}

	// Повторно загружаем книгу с данными пользователя
	if err := db.DB.Preload("User").First(&input, input.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при загрузке пользователя"})
		return
	}

	c.JSON(http.StatusCreated, input)
}

func GetBooks(c *gin.Context) {
	var books []models.Book
	userID, err := utils.GetUserID(c)

	query := db.DB.Preload("User").Model(&models.Book{})

	if err == nil {
		query = query.Where("is_private = ? OR user_id = ?", false, userID)
	} else {
		query = query.Where("is_private = ?", false)
	}

	if err := query.Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при получении книг"})
		return
	}

	c.JSON(http.StatusOK, books)
}

func GetBook(c *gin.Context) {
	var book models.Book
	id := c.Param("id")

	if err := db.DB.Preload("User").First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Книга не найдена"})
		return
	}

	userID, err := utils.GetUserID(c)

	if book.IsPrivate && (err != nil || userID != book.UserID) {
		c.JSON(http.StatusForbidden, gin.H{"message": "Книга приватная"})
		return
	}

	c.JSON(http.StatusOK, book)
}

func UpdateBook(c *gin.Context) {
	id := c.Param("id")

	var book models.Book
	if err := db.DB.Preload("User").First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Книга не найдена"})
		return
	}

	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	userRole, _ := c.Get("userRole")

	if userRole != "admin" && userID != book.UserID {
		c.JSON(http.StatusForbidden, gin.H{"message": "Нет доступа к изменению этой книги"})
		return
	}

	var input models.Book
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Неверные данные"})
		return
	}

	book.Title = input.Title
	book.Author = input.Author
	book.Category = input.Category
	book.ImageUrl = input.ImageUrl
	book.IsPrivate = input.IsPrivate

	if err := db.DB.Save(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при обновлении книги"})
		return
	}

	c.JSON(http.StatusOK, book)
}

func DeleteBook(c *gin.Context) {
	var book models.Book
	id := c.Param("id")

	if err := db.DB.Preload("User").First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Книга не найдена"})
		return
	}

	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	userRole, _ := c.Get("userRole")

	if userRole != "admin" && book.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"message": "У вас нет прав на удаление этой книги"})
		return
	}

	if err := db.DB.Delete(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при удалении книги"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Книга удалена"})
}
