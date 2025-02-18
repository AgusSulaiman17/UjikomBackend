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
		Peminjaman.GET("/user/:id_user", controllers.GetPeminjamanByUserID)
		Peminjaman.PUT("/:id", controllers.UpdatePeminjaman)
		Peminjaman.DELETE("/:id", controllers.DeletePeminjaman)

		// Booking Routes
		Peminjaman.POST("/booking", controllers.CreateBooking)  
		Peminjaman.PUT("/approve/:id", controllers.ApproveBooking)  
		Peminjaman.PUT("/return/:id", controllers.ReturnBook) 
		Peminjaman.DELETE("/booking/:id", controllers.DeleteBooking)  
		Peminjaman.GET("/booking", controllers.GetAllBookings)
		Peminjaman.GET("/booking/user/:id_user", controllers.GetBookingByUserID) 
		Peminjaman.GET("/booking/:id", controllers.GetBookingByID)
	}
}
