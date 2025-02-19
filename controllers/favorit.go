package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)
type CreateFavoritRequest struct {
    IDUser uint `json:"id_user" binding:"required"`
    IDBuku uint `json:"id_buku" binding:"required"`
}

func CreateFavorit(c *gin.Context) {
    var request CreateFavoritRequest

    // Bind JSON ke struct request agar tidak ikut memproses relasi User dan Buku
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid: " + err.Error()})
        return
    }

    // Periksa apakah User dan Buku ada di database
    var user models.User
    if err := config.DB.First(&user, request.IDUser).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
        return
    }

    var buku models.Buku
    if err := config.DB.First(&buku, request.IDBuku).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Buku tidak ditemukan"})
        return
    }

    // Periksa apakah favorit sudah ada
    var existingFavorit models.Favorit
    if err := config.DB.Where("id_user = ? AND id_buku = ?", request.IDUser, request.IDBuku).First(&existingFavorit).Error; err == nil {
        c.JSON(http.StatusConflict, gin.H{"error": "Buku sudah ada di daftar favorit"})
        return
    }

    // Simpan favorit baru
    favorit := models.Favorit{
        IDUser: request.IDUser,
        IDBuku: request.IDBuku,
    }

    if err := config.DB.Create(&favorit).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan ke favorit"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Buku ditambahkan ke favorit", "data": favorit})
}



// Ambil semua favorit pengguna
func GetFavoritByUser(c *gin.Context) {
	userID := c.Param("user_id")
	var favorit []models.Favorit

	// Ambil daftar favorit berdasarkan UserID dan preload relasi User dan Buku
	if err := config.DB.Preload("User").
	Preload("Buku.Penerbit").
	Preload("Buku.Penulis").
	Preload("Buku.Kategori").Preload("Buku").Where("id_user = ?", userID).Find(&favorit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil daftar favorit"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": favorit})
}

// Hapus buku dari favorit
func DeleteFavorit(c *gin.Context) {
	userID := c.Param("user_id")
	bukuID := c.Param("buku_id")

	// Cari favorit berdasarkan UserID dan BukuID
	var favorit models.Favorit
	if err := config.DB.Where("id_user = ? AND id_buku = ?", userID, bukuID).First(&favorit).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Favorit tidak ditemukan"})
		return
	}

	// Hapus favorit dari database
	if err := config.DB.Delete(&favorit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus dari favorit"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Buku dihapus dari favorit"})
}
