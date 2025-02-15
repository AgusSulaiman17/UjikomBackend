package routes

import (
	"backend/controllers"
	"github.com/gin-gonic/gin"
)

func Penulis(router *gin.Engine) {
	penulis := router.Group("/penulis")
	{
		penulis.POST("/", controllers.CreatePenulis)
		penulis.GET("/", controllers.GetAllPenulis)
		penulis.GET("/:id", controllers.GetPenulisByID)
		penulis.PUT("/:id", controllers.UpdatePenulis)
		penulis.DELETE("/:id", controllers.DeletePenulis)
	}
}
