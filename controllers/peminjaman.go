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
	peminjaman.IDUser = idUser

	idBukuStr := c.PostForm("id_buku")
	idBuku, err := utils.ParseUint(idBukuStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID Buku tidak valid"})
		return
	}
	peminjaman.IDBuku = idBuku

	// Ambil data buku dari database
	var buku models.Buku
	if err := config.DB.First(&buku, idBuku).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Buku tidak ditemukan"})
		return
	}

	// Periksa apakah stok buku masih tersedia
	if buku.Jumlah <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stok buku habis, peminjaman tidak dapat dilakukan"})
		return
	}

	// Kurangi jumlah buku
	buku.Jumlah -= 1
	if err := config.DB.Save(&buku).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate jumlah buku"})
		return
	}

	// Set data peminjaman
	peminjaman.DurasiHari = 5
	peminjaman.TanggalPinjam = time.Now()
	peminjaman.TanggalKembali = peminjaman.TanggalPinjam.Add(time.Duration(peminjaman.DurasiHari) * 24 * time.Hour)
	peminjaman.Status = "disetujui"
	peminjaman.StatusKembali = false

	// Simpan peminjaman ke database
	if err := config.DB.Create(&peminjaman).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat peminjaman"})
		return
	}

	// Preload data peminjaman setelah disimpan
	var peminjamanWithDetails models.Peminjaman
	if err := config.DB.Preload("User").
		Preload("Buku.Penerbit").
		Preload("Buku.Penulis").
		Preload("Buku.Kategori").
		First(&peminjamanWithDetails, peminjaman.IDPeminjaman).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data peminjaman setelah dibuat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Peminjaman berhasil dibuat dan disetujui", "data": peminjamanWithDetails})
}




// GetAllPeminjaman - Mengambil semua data peminjaman
func GetAllPeminjaman(c *gin.Context) {
	var peminjaman []models.Peminjaman

	// Cek apakah tabel peminjaman dan buku tersedia
	if config.DB.Migrator().HasTable(&models.Peminjaman{}) && config.DB.Migrator().HasTable(&models.Buku{}) {
		if err := config.DB.
			Preload("User").
			Preload("Buku").
			Preload("Buku.Penerbit").
			Preload("Buku.Penulis").
			Preload("Buku.Kategori").
			Find(&peminjaman).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data peminjaman"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": peminjaman})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tabel peminjaman atau buku tidak ditemukan"})
	}
}

// GetPeminjaman - Mengambil data peminjaman berdasarkan ID
func GetPeminjaman(c *gin.Context) {
	id := c.Param("id")

	var peminjaman models.Peminjaman
	if config.DB.Migrator().HasTable(&models.Peminjaman{}) && config.DB.Migrator().HasTable(&models.Buku{}) {
		if err := config.DB.
			Preload("User").
			Preload("Buku").
			Preload("Buku.Penerbit").
			Preload("Buku.Penulis").
			Preload("Buku.Kategori").
			First(&peminjaman, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Peminjaman tidak ditemukan"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": peminjaman})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tabel peminjaman atau buku tidak ditemukan"})
	}
}

func GetPeminjamanByUserID(c *gin.Context) {
    idUser := c.Param("id_user")
    var peminjaman []models.Peminjaman

    if err := config.DB.Where("id_user = ?", idUser).
        Preload("User").
        Preload("Buku").
        Preload("Buku.Penerbit").
        Preload("Buku.Penulis").
        Preload("Buku.Kategori").
        Find(&peminjaman).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data peminjaman"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Daftar peminjaman user berhasil diambil", "data": peminjaman})
}


func UpdatePeminjaman(c *gin.Context) {
    id := c.Param("id")

    var peminjaman models.Peminjaman
    if err := config.DB.First(&peminjaman, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Peminjaman tidak ditemukan"})
        return
    }

    // Jika buku dikembalikan
    if c.PostForm("status_kembali") == "true" {
        peminjaman.StatusKembali = true
        peminjaman.Status = "dikembalikan"
        peminjaman.TanggalKembali = time.Now()

        // Hitung denda jika terlambat dalam jam
        if peminjaman.TanggalKembali.After(peminjaman.TanggalPinjam.Add(time.Duration(peminjaman.DurasiHari) * 24 * time.Hour)) {
            jamTerlambat := int(time.Since(peminjaman.TanggalPinjam.Add(time.Duration(peminjaman.DurasiHari) * 24 * time.Hour)).Hours())
            peminjaman.Denda = float64(jamTerlambat * 50) // 50 per jam
        } else {
            peminjaman.Denda = 0
        }
    }

    peminjaman.DiperbaruiPada = time.Now()

    if err := config.DB.Save(&peminjaman).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui peminjaman"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Peminjaman berhasil diperbarui", "data": peminjaman})
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

func CreateBooking(c *gin.Context) {
	var peminjaman models.Peminjaman

	// Ambil id_user dari form
	idUserStr := c.PostForm("id_user")
	idUser, err := utils.ParseUint(idUserStr)
	
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID User tidak valid"})
		return
	}
	peminjaman.IDUser = idUser

	// Ambil id_buku dari form
	idBukuStr := c.PostForm("id_buku")
	idBuku, err := utils.ParseUint(idBukuStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID Buku tidak valid"})
		return
	}
	peminjaman.IDBuku = idBuku

	// Set DurasiHari ke 5 hari
	peminjaman.DurasiHari = 5
	peminjaman.TanggalPinjam = time.Now()

	// Set Status menjadi "pending"
	peminjaman.Status = "pending"
	peminjaman.StatusKembali = false

	// Simpan peminjaman
	if err := config.DB.Create(&peminjaman).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat booking"})
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

	c.JSON(http.StatusOK, gin.H{"message": "Booking berhasil dibuat", "data": peminjamanWithDetails})
}


func ApproveBooking(c *gin.Context) {
    id := c.Param("id")
    var peminjaman models.Peminjaman

    if err := config.DB.First(&peminjaman, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Booking tidak ditemukan"})
        return
    }

    peminjaman.Status = "disetujui"
    peminjaman.TanggalKembali = peminjaman.TanggalPinjam.Add(time.Duration(peminjaman.DurasiHari) * 24 * time.Hour)

    if err := config.DB.Save(&peminjaman).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyetujui booking"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Booking disetujui", "data": peminjaman})
}


func ReturnBook(c *gin.Context) {
    id := c.Param("id")
    var peminjaman models.Peminjaman

    if err := config.DB.First(&peminjaman, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Peminjaman tidak ditemukan"})
        return
    }

    if peminjaman.Status != "disetujui" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Buku belum disetujui untuk dipinjam"})
        return
    }

    peminjaman.Status = "dikembalikan"
    peminjaman.StatusKembali = true
    peminjaman.TanggalKembali = time.Now()

    // Periksa apakah keterlambatan terjadi
    if peminjaman.TanggalKembali.After(peminjaman.TanggalPinjam.Add(time.Duration(peminjaman.DurasiHari) * 24 * time.Hour)) {
        jamTerlambat := int(time.Since(peminjaman.TanggalPinjam.Add(time.Duration(peminjaman.DurasiHari) * 24 * time.Hour)).Hours())
        peminjaman.Denda = float64(jamTerlambat * 50) // 50 per jam
    } else {
        peminjaman.Denda = 0
    }

    if err := config.DB.Save(&peminjaman).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengembalikan buku"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Buku berhasil dikembalikan", "data": peminjaman})
}

func DeleteBooking(c *gin.Context) {
    id := c.Param("id")
    var peminjaman models.Peminjaman

    if err := config.DB.First(&peminjaman, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Booking tidak ditemukan"})
        return
    }

    if err := config.DB.Delete(&peminjaman).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus booking"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Booking berhasil dihapus"})
}

func GetAllBookings(c *gin.Context) {
    var peminjaman []models.Peminjaman

    if err := config.DB.Where("status = ?", "pending").Preload("User").Preload("Buku").Find(&peminjaman).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data peminjaman"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Daftar booking yang pending berhasil diambil", "data": peminjaman})
}

func GetBookingByID(c *gin.Context) {
    id := c.Param("id")
    var peminjaman models.Peminjaman

    if err := config.DB.Where("id_peminjaman = ? AND status = ?", id, "pending").
        Preload("User").
        Preload("Buku").
        First(&peminjaman).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Booking tidak ditemukan atau tidak dalam status pending"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Detail booking berhasil diambil", "data": peminjaman})
}
func GetBookingByUserID(c *gin.Context) {
    idUser := c.Param("id_user")
    var peminjaman []models.Peminjaman

    if err := config.DB.Where("id_user = ? AND status = ?", idUser, "pending").
        Preload("User").
        Preload("Buku").
        Find(&peminjaman).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data peminjaman"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Daftar booking user berhasil diambil", "data": peminjaman})
}

