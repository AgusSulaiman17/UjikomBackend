package routes

import (
	"github.com/gin-gonic/gin"
	"backend/controllers"
)

func Peminjaman(router *gin.Engine) {
	Peminjaman := router.Group("/peminjaman")
	{
		Peminjaman.POST("/", controllers.CreatePeminjaman)
		Peminjaman.GET("/", controllers.GetAllPeminjaman)
		Peminjaman.GET("/:id", controllers.GetPeminjaman)
		Peminjaman.PUT("/:id", controllers.UpdatePeminjaman)
		Peminjaman.DELETE("/:id", controllers.DeletePeminjaman)
	}
}