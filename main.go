package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gosimple/internal/db"
	"gosimple/internal/middleware"
	"gosimple/internal/routes"
)

func main() {
	// Инициализация базы данных
	db.InitDB()

	// Создание роутера
	router := gin.Default()

	// Middleware CORS
	router.Use(middleware.CORSMiddleware())

	// Маршруты
	routes.BookRoutes(router)
	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	// Обработчик ошибок
	router.Use(middleware.ErrorMiddleware())
	router.Use(middleware.LoggerMiddleware())

	// Запуск сервера
	port := os.Getenv("PORT")
	dbPort := os.Getenv("DB_PORT")
	log.Println("Значение переменной PORT:", port)
	log.Println("Значение переменной PORT:", dbPort)
	if port == "" || port == "--" {
		port = "8080"
	}

	log.Printf("🚀 Сервер запущен на http://localhost:%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
