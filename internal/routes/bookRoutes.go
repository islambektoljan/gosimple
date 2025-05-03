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

	books.Use(middleware.JWTAuthMiddleware())
	books.POST("", controllers.CreateBook)

	// Маршруты для администраторов
	//books.Use(middleware.RoleMiddleware("admin"))
	//{
	//	// Администраторы могут редактировать и удалять книгу
	//	books.PUT("/:id", controllers.UpdateBook)
	//	books.DELETE("/:id", controllers.DeleteBook)
	//}

	books.Use(middleware.BookOwnerMiddleware())
	{
		books.PUT("/:id", controllers.UpdateBook)
		books.DELETE("/:id", controllers.DeleteBook)
	}
}
