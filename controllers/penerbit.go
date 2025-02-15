package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreatePenerbit(c *gin.Context) {
	var penerbit models.Penerbit

	// Bind JSON input ke struct
	if err := c.ShouldBindJSON(&penerbit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	// Validasi nama tidak boleh kosong
	if penerbit.Nama == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama penerbit tidak boleh kosong"})
		return
	}

	// Simpan ke database
	if err := config.DB.Create(&penerbit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat penerbit"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Penerbit berhasil dibuat",
		"data":    penerbit,
	})
}

// GetPenerbitByID handles getting a penerbit by ID
func GetPenerbitByID(c *gin.Context) {
	id := c.Param("id")
	var penerbit models.Penerbit

	// Find the penerbit by ID
	if err := config.DB.First(&penerbit, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Penerbit tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": penerbit})
}

// GetAllPenerbit handles getting all penerbit
func GetAllPenerbit(c *gin.Context) {
	var penerbits []models.Penerbit

	// Get all penerbit from the database
	if err := config.DB.Find(&penerbits).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data penerbit"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": penerbits})
}

func UpdatePenerbit(c *gin.Context) {
	// Ambil ID penerbit dari parameter URL
	id := c.Param("id")

	// Cek apakah penerbit dengan ID tersebut ada
	var penerbit models.Penerbit
	if err := config.DB.First(&penerbit, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Penerbit tidak ditemukan"})
		return
	}

	// Bind JSON input ke struct
	var input models.Penerbit
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	// Validasi nama tidak boleh kosong
	if input.Nama == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama penerbit tidak boleh kosong"})
		return
	}

	// Update data penerbit
	penerbit.Nama = input.Nama
	if err := config.DB.Save(&penerbit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui penerbit"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Penerbit berhasil diperbarui",
		"data":    penerbit,
	})
}

// DeletePenerbit handles deleting a penerbit by ID
func DeletePenerbit(c *gin.Context) {
	id := c.Param("id")
	var penerbit models.Penerbit

	// Find the penerbit by ID
	if err := config.DB.First(&penerbit, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Penerbit tidak ditemukan"})
		return
	}

	// Delete the penerbit
	if err := config.DB.Delete(&penerbit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus penerbit"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Penerbit berhasil dihapus"})
}