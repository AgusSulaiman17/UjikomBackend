package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	userGroup := router.Group("/users")
	{
		userGroup.GET("/", controllers.GetAllUsers)
		userGroup.GET("/:id", controllers.GetUserByID)
		userGroup.POST("/", controllers.CreateUser)
		userGroup.PUT("/:id", controllers.UpdateUser)
		userGroup.DELETE("/:id", controllers.DeleteUser)
	}
}
