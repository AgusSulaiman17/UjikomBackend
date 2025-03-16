package controllers

import (
	"backend/config"
	"backend/models"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Create Buku
func CreateBuku(c *gin.Context) {
	// Form data
	var input models.Buku
	input.Judul = c.PostForm("judul")
	input.IDPenerbit = uint(parseUint(c.PostForm("id_penerbit")))
	input.IDPenulis = uint(parseUint(c.PostForm("id_penulis")))
	input.IDKategori = uint(parseUint(c.PostForm("id_kategori")))
	input.Deskripsi = c.PostForm("deskripsi")
	input.Jumlah = int(parseUint(c.PostForm("jumlah")))
	input.Status = true // default status tersedia

	// Ambil ISBN
	input.ISBN = c.PostForm("isbn")
	if input.ISBN == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ISBN wajib diisi"})
		return
	}

	// Cek apakah ISBN sudah ada
	var existingBuku models.Buku
	if err := config.DB.Where("isbn = ?", input.ISBN).First(&existingBuku).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "ISBN sudah digunakan untuk buku lain"})
		return
	}


	// Proses file gambar
	file, err := c.FormFile("gambar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gambar wajib diunggah"})
		return
	}

	// Simpan file ke folder uploads
	uploadPath := filepath.Join("uploads", file.Filename)
	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan gambar"})
		return
	}

	input.Gambar = uploadPath // Simpan path gambar ke database

	// Simpan data buku ke database
	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan buku"})
		return
	}

	// Preload relasi setelah create
	if err := config.DB.Preload("Penerbit").Preload("Penulis").Preload("Kategori").First(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat relasi buku"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Buku berhasil ditambahkan", "data": input})
}

// Helper untuk parse string ke uint
func parseUint(value string) uint64 {
	parsed, _ := strconv.ParseUint(value, 10, 64)
	return parsed
}

// Get All Buku
func GetAllBuku(c *gin.Context) {
	var buku []models.Buku
	// Preload relasi
	if err := config.DB.Preload("Penerbit").Preload("Penulis").Preload("Kategori").Find(&buku).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data buku"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": buku})
}

// Get Buku by ID
func GetBukuByID(c *gin.Context) {
	id := c.Param("id")
	var buku models.Buku

	// Preload relasi
	if err := config.DB.Preload("Penerbit").Preload("Penulis").Preload("Kategori").First(&buku, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Buku tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": buku})
}

// Update Buku
func UpdateBuku(c *gin.Context) {
	// Ambil ID dari parameter URL
	id := c.Param("id")

	// Cari buku berdasarkan ID
	var buku models.Buku
	if err := config.DB.First(&buku, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Buku tidak ditemukan"})
		return
	}

	// Form data
	judul := c.PostForm("judul")
	idPenerbit := uint(parseUint(c.PostForm("id_penerbit")))
	idPenulis := uint(parseUint(c.PostForm("id_penulis")))
	idKategori := uint(parseUint(c.PostForm("id_kategori")))
	deskripsi := c.PostForm("deskripsi")
	jumlah := int(parseUint(c.PostForm("jumlah")))
	status := c.PostForm("status") == "true" // Mengubah status ke boolean

	// Ambil ISBN jika ada perubahan
	isbn := c.PostForm("isbn")
	if isbn != "" {
		// Cek apakah ISBN sudah ada di buku lain
		var existingBuku models.Buku
		if err := config.DB.Where("isbn = ?", isbn).First(&existingBuku).Error; err == nil && existingBuku.IDBuku != buku.IDBuku {
			c.JSON(http.StatusConflict, gin.H{"error": "ISBN sudah digunakan untuk buku lain"})
			return
		}
		buku.ISBN = isbn
	}

	// Update field jika ada perubahan
	if judul != "" {
		buku.Judul = judul
	}

	if idPenerbit != 0 {
		buku.IDPenerbit = idPenerbit
	}
	if idPenulis != 0 {
		buku.IDPenulis = idPenulis
	}
	if idKategori != 0 {
		buku.IDKategori = idKategori
	}
	if deskripsi != "" {
		buku.Deskripsi = deskripsi
	}
	if jumlah != 0 {
		buku.Jumlah = jumlah
	}
	buku.Status = status

	// Periksa jika ada file gambar baru
	file, err := c.FormFile("gambar")
	if err == nil { // Jika file ada
		// Simpan file gambar baru ke folder uploads
		uploadPath := filepath.Join("uploads", file.Filename)
		if err := c.SaveUploadedFile(file, uploadPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan gambar"})
			return
		}

		// Hapus gambar lama jika ada
		if buku.Gambar != "" {
			if err := os.Remove(buku.Gambar); err != nil {
				fmt.Printf("Gagal menghapus file lama: %v\n", err)
			}
		}

		// Update path gambar baru
		buku.Gambar = uploadPath
	}

	// Update waktu diperbarui
	buku.DiperbaruiPada = time.Now()

	// Simpan perubahan ke database
	if err := config.DB.Save(&buku).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui buku"})
		return
	}

	// Preload relasi setelah update
	if err := config.DB.Preload("Penerbit").Preload("Penulis").Preload("Kategori").First(&buku).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat relasi buku"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Buku berhasil diperbarui", "data": buku})
}

// Delete Buku
func DeleteBuku(c *gin.Context) {
	id := c.Param("id")
	var buku models.Buku

	// Cek apakah buku ada
	if err := config.DB.First(&buku, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Buku tidak ditemukan"})
		return
	}

	// Hapus buku
	if err := config.DB.Delete(&buku).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus buku"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Buku berhasil dihapus"})
}
