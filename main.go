package main

import (
	"backend/config"
	"backend/routes"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Koneksi ke database
	config.ConnectDB()


	// Inisialisasi router Gin
	router := gin.Default()
	

	// Middleware CORS agar tidak terjadi masalah akses dari frontend
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Sesuaikan dengan frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           36 * time.Hour,
	}))

	// Tambahkan log untuk memastikan middleware berjalan
	router.Use(func(c *gin.Context) {
		fmt.Println("CORS Middleware applied")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// Log header yang dikirim
		fmt.Println("Response Headers:")
		fmt.Println("Access-Control-Allow-Origin:", c.Writer.Header().Get("Access-Control-Allow-Origin"))
		fmt.Println("Access-Control-Allow-Methods:", c.Writer.Header().Get("Access-Control-Allow-Methods"))
		fmt.Println("Access-Control-Allow-Headers:", c.Writer.Header().Get("Access-Control-Allow-Headers"))

		// Jika request OPTIONS, langsung response 204 (tanpa konten)
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Register routes
	routes.Login(router)
	routes.UserRoutes(router)
	routes.Penerbit(router)
	routes.Kategori(router)
	routes.Penulis(router)
	routes.BukuRoutes(router)
	routes.Peminjaman(router)
	router.Static("/uploads", "uploads")

	// Jalankan server di port 8080
	router.Run(":8080")
}
