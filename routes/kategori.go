package routes

import (
	"backend/controllers"
	"github.com/gin-gonic/gin"
)

func Kategori(router *gin.Engine) {
	kategoriRoutes := router.Group("/kategori")
	{
		kategoriRoutes.GET("/", controllers.GetAllKategori)
		kategoriRoutes.GET("/:id", controllers.GetKategoriByID)
		kategoriRoutes.POST("/", controllers.CreateKategori)
		kategoriRoutes.PUT("/:id", controllers.UpdateKategori)
		kategoriRoutes.DELETE("/:id", controllers.DeleteKategori)
	}
}
