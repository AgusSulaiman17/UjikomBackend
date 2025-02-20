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

	// Hitung jumlah peminjaman dan booking gabungan
	var totalCount int64
	if err := config.DB.Model(&models.Peminjaman{}).
		Where("id_user = ? AND (status = 'disetujui' OR status = 'pending')", idUser).
		Count(&totalCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung jumlah peminjaman dan booking"})
		return
	}
	if totalCount >= 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Anda sudah memiliki total 5 peminjaman dan booking"})
		return
	}

	// Cek apakah user sudah meminjam buku yang sama dan belum mengembalikannya
	var existingPeminjaman models.Peminjaman
	if err := config.DB.Where("id_user = ? AND id_buku = ? AND status = 'disetujui' AND status_kembali = false", idUser, idBuku).
		First(&existingPeminjaman).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Anda sudah meminjam buku ini, kembalikan terlebih dahulu sebelum meminjam lagi"})
		return
	}

	// Ambil data buku dari database
	var buku models.Buku
	if err := config.DB.First(&buku, idBuku).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Buku tidak ditemukan"})
		return
	}

	// Periksa stok buku
	if buku.Jumlah <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stok buku habis, peminjaman tidak dapat dilakukan"})
		return
	}

	// Kurangi stok buku
	buku.Jumlah -= 1
	if err := config.DB.Save(&buku).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate jumlah buku"})
		return
	}

	// Set data peminjaman
	peminjaman.DurasiHari = 1
	peminjaman.TanggalPinjam = time.Now()
	peminjaman.TanggalKembali = peminjaman.TanggalPinjam.Add(24 * time.Hour)
	peminjaman.Status = "disetujui"
	peminjaman.StatusKembali = false

	// Simpan peminjaman ke database
	if err := config.DB.Create(&peminjaman).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat peminjaman"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Peminjaman berhasil dibuat dan disetujui", "data": peminjaman})
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

// CreateBooking - Membuat booking baru
func CreateBooking(c *gin.Context) {
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

	// Hitung jumlah peminjaman dan booking gabungan
	var totalCount int64
	if err := config.DB.Model(&models.Peminjaman{}).
		Where("id_user = ? AND (status = 'disetujui' OR status = 'pending')", idUser).
		Count(&totalCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung jumlah peminjaman dan booking"})
		return
	}
	if totalCount >= 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Anda sudah memiliki total 5 peminjaman dan booking"})
		return
	}

	// Cek apakah user sudah membooking buku yang sama
	var existingBooking models.Peminjaman
	if err := config.DB.Where("id_user = ? AND id_buku = ? AND status = 'pending'", idUser, idBuku).
		First(&existingBooking).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Anda sudah membooking buku ini, tunggu hingga diproses"})
		return
	}

	// Set data booking
	peminjaman.DurasiHari = 1
	peminjaman.TanggalPinjam = time.Now()
	peminjaman.Status = "pending"
	peminjaman.StatusKembali = false

	// Simpan booking ke database
	if err := config.DB.Create(&peminjaman).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat booking"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking berhasil dibuat", "data": peminjaman})
}

func ApproveBooking(c *gin.Context) {
    id := c.Param("id")
    var peminjaman models.Peminjaman

    // Cari peminjaman berdasarkan ID
    if err := config.DB.First(&peminjaman, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Booking tidak ditemukan"})
        return
    }

    // Ambil data buku berdasarkan IDBuku
    var buku models.Buku
    if err := config.DB.First(&buku, peminjaman.IDBuku).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Buku tidak ditemukan"})
        return
    }

    // Periksa apakah stok buku masih tersedia
    if buku.Jumlah <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Stok buku habis, tidak dapat menyetujui booking"})
        return
    }

    // Kurangi jumlah buku
    buku.Jumlah -= 1
    if err := config.DB.Save(&buku).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate jumlah buku"})
        return
    }

    // Set status peminjaman menjadi disetujui
    peminjaman.Status = "disetujui"
    peminjaman.TanggalKembali = peminjaman.TanggalPinjam.Add(time.Duration(peminjaman.DurasiHari) * 24 * time.Hour)

    // Simpan perubahan peminjaman ke database
    if err := config.DB.Save(&peminjaman).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyetujui booking"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Booking disetujui, jumlah buku berkurang", "data": peminjaman})
}



func ReturnBook(c *gin.Context) {
    id := c.Param("id")
    var peminjaman models.Peminjaman

    // Cari peminjaman berdasarkan ID
    if err := config.DB.First(&peminjaman, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Peminjaman tidak ditemukan"})
        return
    }

    if peminjaman.Status != "disetujui" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Buku belum disetujui untuk dipinjam"})
        return
    }

    // Ambil data buku berdasarkan IDBuku
    var buku models.Buku
    if err := config.DB.First(&buku, peminjaman.IDBuku).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Buku tidak ditemukan"})
        return
    }

    // Tambah jumlah buku setelah dikembalikan
    buku.Jumlah += 1
    if err := config.DB.Save(&buku).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate jumlah buku"})
        return
    }

    // Set status peminjaman menjadi dikembalikan
    peminjaman.Status = "dikembalikan"
    peminjaman.StatusKembali = true
    peminjaman.TanggalKembali = time.Now()

    // Periksa apakah keterlambatan terjadi
    tanggalJatuhTempo := peminjaman.TanggalPinjam.Add(time.Duration(peminjaman.DurasiHari) * 24 * time.Hour)
    terlambatDurasi := peminjaman.TanggalKembali.Sub(tanggalJatuhTempo)
    hariTerlambat := int(terlambatDurasi.Hours() / 24)

    if hariTerlambat > 0 {
        peminjaman.Denda = float64(hariTerlambat * 10000)
    } else {
        peminjaman.Denda = 0
    }

    // Simpan perubahan peminjaman ke database
    if err := config.DB.Save(&peminjaman).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengembalikan buku"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Buku berhasil dikembalikan, stok buku bertambah", "data": peminjaman})
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

    if err := config.DB.Where("status = ?", "pending").Preload("User").Preload("Buku"). Preload("Buku.Penerbit").
	Preload("Buku.Penulis").
	Preload("Buku.Kategori").Find(&peminjaman).Error; err != nil {
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
		Preload("Buku.Penerbit").
        Preload("Buku.Penulis").
        Preload("Buku.Kategori").
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
		Preload("Buku.Penerbit").
        Preload("Buku.Penulis").
        Preload("Buku.Kategori").
        Find(&peminjaman).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data peminjaman"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Daftar booking user berhasil diambil", "data": peminjaman})
}

