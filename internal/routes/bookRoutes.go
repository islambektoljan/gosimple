package routes

import (
	"github.com/gin-gonic/gin"
	"gosimple/internal/controllers"
	"gosimple/internal/middleware"
)

func BookRoutes(router *gin.Engine) {
	books := router.Group("/books")

	books.GET("", controllers.GetBooks)
	books.GET("/:id", controllers.GetBook)

	// Маршруты для авторизованных пользователей
	books.Use(middleware.JWTAuthMiddleware())
	books.POST("", controllers.CreateBook)

	// Маршруты для администраторов
	//books.Use(middleware.RoleMiddleware("admin"))
	//{
	//	// Администраторы могут редактировать и удалять книгу
	//	books.PUT("/:id", controllers.UpdateBook)
	//	books.DELETE("/:id", controllers.DeleteBook)
	//}

	// Маршруты для владельцев книги
	books.Use(middleware.BookOwnerMiddleware())
	{
		// Владельцы книги могут редактировать и удалять только свои книги
		books.PUT("/:id", controllers.UpdateBook)
		books.DELETE("/:id", controllers.DeleteBook)
	}
}
