package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func Favorit(router *gin.Engine) {
	favorit := router.Group("/favorit")
	{
		favorit.POST("/", controllers.CreateFavorit)           // Tambahkan buku ke favorit
		favorit.GET("/:user_id", controllers.GetFavoritByUser) // Ambil daftar favorit berdasarkan user
		favorit.DELETE("/:user_id/:buku_id", controllers.DeleteFavorit) // Hapus buku dari favorit
	}
}
