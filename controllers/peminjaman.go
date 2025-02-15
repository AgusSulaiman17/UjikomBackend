package controllers

import (
	"net/http"
	"time"

	"backend/config"
	"backend/models"
	"backend/utils"
	"github.com/gin-gonic/gin"
)

// CreatePeminjaman - Membuat peminjaman baru
func CreatePeminjaman(c *gin.Context) {
	var peminjaman models.Peminjaman

	idUserStr := c.PostForm("id_user")
	idUser, err := utils.ParseUint(idUserStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID User tidak valid"})
		return
	}

	// Ambil data peminjaman lainnya dari form
	peminjaman.IDUser = idUser // Set ID User

	idBukuStr := c.PostForm("id_buku")
	idBuku, err := utils.ParseUint(idBukuStr) // Parse id_buku
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID Buku tidak valid"})
		return
	}
	peminjaman.IDBuku = idBuku // Set ID Buku

	// Set DurasiHari 5 hari
	peminjaman.DurasiHari = 5

	// Menghitung Tanggal Kembali berdasarkan DurasiHari
	peminjaman.TanggalPinjam = time.Now()
	peminjaman.TanggalKembali = peminjaman.TanggalPinjam.Add(time.Duration(peminjaman.DurasiHari) * 24 * time.Hour)

	// Set Status Kembali
	peminjaman.StatusKembali = false

	// Simpan peminjaman
	if err := config.DB.Create(&peminjaman).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat peminjaman"})
		return
	}

	// Preload relasi setelah peminjaman berhasil dibuat
	var peminjamanWithDetails models.Peminjaman
	if err := config.DB.Preload("User").
		Preload("Buku.Penerbit").
		Preload("Buku.Penulis").
		Preload("Buku.Kategori").
		First(&peminjamanWithDetails, peminjaman.IDPeminjaman).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data peminjaman setelah dibuat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Peminjaman berhasil dibuat", "data": peminjamanWithDetails})
}


// GetAllPeminjaman - Mengambil semua data peminjaman
func GetAllPeminjaman(c *gin.Context) {
	var peminjaman []models.Peminjaman

	// Ambil semua data peminjaman dan preload relasi User, Buku, Penerbit, Penulis, dan Kategori
	if err := config.DB.Preload("User").
		Preload("Buku.Penerbit").
		Preload("Buku.Penulis").
		Preload("Buku.Kategori").
		Find(&peminjaman).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data peminjaman"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": peminjaman})
}

// GetPeminjaman - Mengambil data peminjaman berdasarkan ID
func GetPeminjaman(c *gin.Context) {
	id := c.Param("id")

	var peminjaman models.Peminjaman
	if err := config.DB.Preload("User").
		Preload("Buku.Penerbit").
		Preload("Buku.Penulis").
		Preload("Buku.Kategori").
		First(&peminjaman, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Peminjaman tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": peminjaman})
}

// UpdatePeminjaman - Memperbarui status peminjaman dan menghitung denda
func UpdatePeminjaman(c *gin.Context) {
	id := c.Param("id")

	var peminjaman models.Peminjaman
	if err := config.DB.First(&peminjaman, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Peminjaman tidak ditemukan"})
		return
	}

	// Jika sudah dikembalikan
	if c.PostForm("status_kembali") == "true" {
		peminjaman.StatusKembali = true
		peminjaman.TanggalKembali = time.Now() // Tanggal pengembalian adalah saat ini
		// Hitung denda jika terlambat
		lateDuration := time.Since(peminjaman.TanggalKembali)
		if lateDuration.Hours() > float64(peminjaman.DurasiHari*24) {
			// Hitung denda: 100 per jam keterlambatan
			lateHours := lateDuration.Hours() - float64(peminjaman.DurasiHari*24)
			peminjaman.Denda = lateHours * 100
		}
	}

	// Update waktu diperbarui
	peminjaman.DiperbaruiPada = time.Now()

	// Simpan perubahan ke database
	if err := config.DB.Save(&peminjaman).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui peminjaman"})
		return
	}

	// Preload relasi setelah update
	var peminjamanWithDetails models.Peminjaman
	if err := config.DB.Preload("User").
		Preload("Buku.Penerbit").
		Preload("Buku.Penulis").
		Preload("Buku.Kategori").
		First(&peminjamanWithDetails, peminjaman.IDPeminjaman).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data peminjaman setelah update"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Peminjaman berhasil diperbarui", "data": peminjamanWithDetails})
}


// DeletePeminjaman - Menghapus peminjaman berdasarkan ID
func DeletePeminjaman(c *gin.Context) {
	id := c.Param("id")

	var peminjaman models.Peminjaman
	if err := config.DB.First(&peminjaman, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Peminjaman tidak ditemukan"})
		return
	}

	// Hapus data peminjaman
	if err := config.DB.Delete(&peminjaman).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus peminjaman"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Peminjaman berhasil dihapus"})
}
