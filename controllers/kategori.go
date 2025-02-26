package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Get all categories
func GetAllKategori(c *gin.Context) {
	var kategori []models.Kategori
	if err := config.DB.Find(&kategori).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": kategori})
}

// Get a single category by ID
func GetKategoriByID(c *gin.Context) {
	id := c.Param("id")
	var kategori models.Kategori

	if err := config.DB.First(&kategori, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": kategori})
}

// CreateKategori handles creating a new category
func CreateKategori(c *gin.Context) {
	var input struct {
		Kategori string `json:"kategori" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi: Cek apakah kategori sudah ada
	var existingKategori models.Kategori
	if err := config.DB.Where("kategori = ?", input.Kategori).First(&existingKategori).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kategori sudah digunakan"})
		return
	}

	// Simpan kategori ke database
	kategori := models.Kategori{Kategori: input.Kategori}
	if err := config.DB.Create(&kategori).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat kategori"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Kategori berhasil dibuat", "data": kategori})
}

// UpdateKategori handles updating an existing category
func UpdateKategori(c *gin.Context) {
	id := c.Param("id")
	var kategori models.Kategori

	// Cek apakah kategori dengan ID tersebut ada
	if err := config.DB.First(&kategori, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
		return
	}

	// Bind JSON input ke struct
	var input struct {
		Kategori string `json:"kategori" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi: Pastikan kategori tidak kosong
	if input.Kategori == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama kategori tidak boleh kosong"})
		return
	}

	// Hanya lakukan pengecekan jika kategori berubah
	if input.Kategori != kategori.Kategori {
		var existingKategori models.Kategori
		if err := config.DB.Where("kategori = ?", input.Kategori).First(&existingKategori).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Kategori sudah digunakan"})
			return
		}
	}

	// Update kategori
	kategori.Kategori = input.Kategori
	if err := config.DB.Save(&kategori).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui kategori"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kategori berhasil diperbarui", "data": kategori})
}

// Delete a category
func DeleteKategori(c *gin.Context) {
	id := c.Param("id")
	var kategori models.Kategori

	if err := config.DB.First(&kategori, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	if err := config.DB.Delete(&kategori).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
