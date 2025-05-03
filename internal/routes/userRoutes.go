package routes

import (
	"github.com/gin-gonic/gin"
	"gosimple/internal/controllers"
	"gosimple/internal/middleware"
)

func UserRoutes(r *gin.Engine) {
	user := r.Group("/profile")
	user.Use(middleware.JWTAuthMiddleware())

	user.GET("", controllers.GetProfile)
	user.PUT("/update", controllers.UpdateUser)
	user.DELETE("/delete", controllers.DeleteUser)
}
