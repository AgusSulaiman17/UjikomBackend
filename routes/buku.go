package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func BukuRoutes(router *gin.Engine) {
	buku := router.Group("/buku")
	{
		buku.POST("/", controllers.CreateBuku)     // Create buku
		buku.GET("/", controllers.GetAllBuku)     // Get all buku
		buku.GET("/:id", controllers.GetBukuByID) // Get buku by ID
		buku.PUT("/:id", controllers.UpdateBuku)  // Update buku
		buku.DELETE("/:id", controllers.DeleteBuku) // Delete buku
	}
}