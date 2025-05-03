package routes

import (
	"github.com/gin-gonic/gin"
	"gosimple/internal/controllers"
)

func AuthRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", controllers.RegisterUser)
		auth.POST("/login", controllers.LoginUser)
	}
}
