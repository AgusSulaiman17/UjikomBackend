package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func Penerbit(router *gin.Engine) {
	Penerbit := router.Group("/penerbit")
	{
		Penerbit.POST("/", controllers.CreatePenerbit)
		Penerbit.GET("/", controllers.GetAllPenerbit)
		Penerbit.GET("/:id", controllers.GetPenerbitByID)
		Penerbit.PUT("/:id", controllers.UpdatePenerbit)
		Penerbit.DELETE("/:id", controllers.DeletePenerbit)
	}
}