package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"
	
	"github.com/gin-gonic/gin"
)

// CreatePenulis handles creating a new penulis
func CreatePenulis(c *gin.Context) {
	var input models.Penulis

	// Bind JSON input ke struct Penulis
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid: " + err.Error()})
		return
	}

	// Validasi: Cek apakah nama penulis sudah ada
	var existingPenulis models.Penulis
	if err := config.DB.Where("nama = ?", input.Nama).First(&existingPenulis).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama penulis sudah digunakan"})
		return
	}

	// Simpan penulis ke database
	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat penulis"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Penulis berhasil dibuat", "data": input})
}


// GetPenulisByID handles getting a penulis by ID
func GetPenulisByID(c *gin.Context) {
	id := c.Param("id")
	var penulis models.Penulis

	// Find the penulis by ID
	if err := config.DB.First(&penulis, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Penulis tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": penulis})
}

// GetAllPenulis handles getting all penulis
func GetAllPenulis(c *gin.Context) {
	var penulis []models.Penulis

	// Get all penulis from the database
	if err := config.DB.Find(&penulis).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data penulis"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": penulis})
}

// UpdatePenulis handles updating an existing penulis
func UpdatePenulis(c *gin.Context) {
	// Ambil ID penulis dari parameter URL
	penulisId := c.Param("id")

	// Cek apakah penulis dengan ID tersebut ada
	var penulis models.Penulis
	if err := config.DB.First(&penulis, penulisId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Penulis tidak ditemukan"})
		return
	}

	// Bind JSON input ke struct
	var input models.Penulis
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	// Validasi nama tidak boleh kosong
	if input.Nama == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama penulis tidak boleh kosong"})
		return
	}

	// Hanya lakukan pengecekan jika nama penulis berubah
	if input.Nama != penulis.Nama {
		var existingPenulis models.Penulis
		if err := config.DB.Where("nama = ?", input.Nama).First(&existingPenulis).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nama penulis sudah digunakan oleh penulis lain"})
			return
		}
	}

	// Update data penulis
	penulis.Nama = input.Nama
	if err := config.DB.Save(&penulis).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data penulis"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Penulis berhasil diperbarui", "data": penulis})
}



// DeletePenulis handles deleting a penulis by ID
func DeletePenulis(c *gin.Context) {
	id := c.Param("id")
	var penulis models.Penulis

	// Find the penulis by ID
	if err := config.DB.First(&penulis, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Penulis tidak ditemukan"})
		return
	}

	// Delete the penulis
	if err := config.DB.Delete(&penulis).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus penulis"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Penulis berhasil dihapus"})
}
