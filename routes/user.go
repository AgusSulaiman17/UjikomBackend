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
		userGroup.GET("/unapp/", controllers.GetUnapprovedUsers)
		userGroup.POST("/app/", controllers.CreateApprovedUser)
		userGroup.PUT("/approve/:id", controllers.ApproveUser)
		userGroup.POST("/", controllers.CreateUser)
		userGroup.PUT("/:id", controllers.UpdateUser)
		userGroup.DELETE("/:id", controllers.DeleteUser)
	}
}
