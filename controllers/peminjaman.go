package controllers

import (
	"net/http"
	"time"

	"backend/config"
	"backend/models"
	"backend/utils"
	"github.com/gin-gonic/gin"
    "fmt"
)

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

	// Ambil data user untuk notifikasi email
	var user models.User
	if err := config.DB.First(&user, idUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pengguna"})
		return
	}

	// Kurangi stok buku
	buku.Jumlah -= 1
	if err := config.DB.Save(&buku).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate jumlah buku"})
		return
	}

	// Set data peminjaman
	peminjaman.DurasiHari = 5
	peminjaman.TanggalPinjam = time.Now()
	peminjaman.TanggalKembali = peminjaman.TanggalPinjam.Add(5 * 24 * time.Hour) // 5 hari
	peminjaman.Status = "disetujui"
	peminjaman.StatusKembali = false

	// Simpan peminjaman ke database
	if err := config.DB.Create(&peminjaman).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat peminjaman"})
		return
	}

	// Kirim email notifikasi
	subject := "Peminjaman Buku Berhasil"
	body := fmt.Sprintf(
		"Halo %s,\n\nAnda telah berhasil meminjam buku dengan judul '%s'.\n\nTanggal Peminjaman: %s\nTanggal Kembali: %s\n\nHarap kembalikan buku sebelum tanggal kembali untuk menghindari denda.\n\nTerima kasih telah menggunakan layanan kami.\n\nSalam,\nPerpustakaan",
		user.Name, buku.Judul, peminjaman.TanggalPinjam.Format("02-01-2006"), peminjaman.TanggalKembali.Format("02-01-2006"),
	)

	err = utils.SendEmail(user.Email, subject, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Peminjaman berhasil dibuat, namun gagal mengirim email notifikasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Peminjaman berhasil dibuat dan email notifikasi telah dikirim", "data": peminjaman})
}


func GetAllPeminjaman(c *gin.Context) {
    var peminjaman []models.Peminjaman

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

// GetPeminjamanByUserID - Mengambil daftar peminjaman berdasarkan ID user
func GetPeminjamanByUserID(c *gin.Context) {
    idUser := c.Param("id_user")
    var peminjaman []models.Peminjaman

    if err := config.DB.Where("id_user = ? AND is_deleted_by_user = false", idUser).
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

// DeletePeminjamanByUser - Menandai peminjaman sebagai dihapus oleh user
func DeletePeminjamanByUser(c *gin.Context) {
    id := c.Param("id")

    var peminjaman models.Peminjaman
    if err := config.DB.First(&peminjaman, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Peminjaman tidak ditemukan"})
        return
    }

    // Tandai sebagai dihapus oleh user, bukan benar-benar menghapus
    peminjaman.IsDeletedByUser = true
    if err := config.DB.Save(&peminjaman).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus peminjaman"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Peminjaman berhasil dihapus oleh user"})
}

func CreateBooking(c *gin.Context) {
	var peminjaman models.Peminjaman
	var buku models.Buku

	// Ambil ID User dari request
	idUserStr := c.PostForm("id_user")
	idUser, err := utils.ParseUint(idUserStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID User tidak valid"})
		return
	}
	peminjaman.IDUser = idUser

	// Ambil ID Buku dari request
	idBukuStr := c.PostForm("id_buku")
	idBuku, err := utils.ParseUint(idBukuStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID Buku tidak valid"})
		return
	}
	peminjaman.IDBuku = idBuku

	// Cek apakah buku tersedia
	if err := config.DB.First(&buku, idBuku).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Buku tidak ditemukan"})
		return
	}
	if buku.Jumlah <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Buku sedang tidak tersedia"})
		return
	}

	// Hitung total peminjaman dan booking user
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

	// Ambil data user untuk notifikasi email
	var user models.User
	if err := config.DB.First(&user, idUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pengguna"})
		return
	}

	// Kurangi jumlah buku
	buku.Jumlah -= 1
	if err := config.DB.Save(&buku).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate jumlah buku"})
		return
	}

	// Set data booking
	peminjaman.DurasiHari = 5
	peminjaman.TanggalPinjam = time.Now()
	peminjaman.Status = "pending"
	peminjaman.StatusKembali = false

	// Simpan booking ke database
	if err := config.DB.Create(&peminjaman).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat booking"})
		return
	}

	// Kirim email notifikasi
	subject := "Booking Buku Berhasil Dibuat"
	body := fmt.Sprintf(
		"Halo %s,\n\nAnda telah berhasil membooking buku dengan judul '%s'.\n\nStatus booking Anda saat ini: %s\nSilakan menunggu konfirmasi dari admin.\n\nSalam,\nPerpustakaan",
		user.Name, buku.Judul, peminjaman.Status,
	)

	err = utils.SendEmail(user.Email, subject, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Booking berhasil dibuat, namun gagal mengirim email notifikasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking berhasil dibuat dan email notifikasi telah dikirim", "data": peminjaman})
}

func ApproveBooking(c *gin.Context) {
    id := c.Param("id")
    var peminjaman models.Peminjaman

    // Cari peminjaman berdasarkan ID
    if err := config.DB.First(&peminjaman, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Booking tidak ditemukan"})
        return
    }

    // Cek apakah peminjaman sudah disetujui sebelumnya
    if peminjaman.Status == "disetujui" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Booking sudah disetujui sebelumnya"})
        return
    }

    // Ambil data buku berdasarkan IDBuku
    var buku models.Buku
    if err := config.DB.First(&buku, peminjaman.IDBuku).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Buku tidak ditemukan"})
        return
    }

    // Ambil data user untuk notifikasi email
    var user models.User
    if err := config.DB.First(&user, peminjaman.IDUser).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pengguna"})
        return
    }

    // Set status peminjaman menjadi disetujui
    peminjaman.Status = "disetujui"

    // Hitung tanggal kembali berdasarkan tanggal pinjam
    if peminjaman.TanggalPinjam.IsZero() {
        peminjaman.TanggalPinjam = time.Now() // Jika kosong, gunakan waktu sekarang
    }
    peminjaman.TanggalKembali = peminjaman.TanggalPinjam.Add(time.Duration(peminjaman.DurasiHari) * 24 * time.Hour)

    // Simpan perubahan peminjaman ke database
    if err := config.DB.Save(&peminjaman).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyetujui booking"})
        return
    }

    // Kirim email notifikasi
    subject := "Booking Buku Disetujui"
    body := fmt.Sprintf(
        "Halo %s,\n\nBooking buku dengan judul '%s' telah disetujui.\nSilakan mengambil buku di perpustakaan.\n\nTanggal Peminjaman: %s\nTanggal Kembali: %s\n\nSalam,\nPerpustakaan",
        user.Name, buku.Judul, peminjaman.TanggalPinjam.Format("02-01-2006"), peminjaman.TanggalKembali.Format("02-01-2006"),
    )

    err := utils.SendEmail(user.Email, subject, body)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Booking berhasil disetujui, namun gagal mengirim email notifikasi"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Booking disetujui dan email notifikasi telah dikirim", "data": peminjaman})
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

    // Ambil data user untuk notifikasi email
    var user models.User
    if err := config.DB.First(&user, peminjaman.IDUser).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pengguna"})
        return
    }

    // Kirim email notifikasi
    subject := "Konfirmasi Pengembalian Buku"
    body := fmt.Sprintf("Halo %s,\n\nBuku dengan judul '%s' telah berhasil dikembalikan.\n", user.Name, buku.Judul)

    if hariTerlambat > 0 {
        body += fmt.Sprintf("Namun, Anda terlambat mengembalikan buku selama %d hari.\nDenda yang harus dibayar: Rp%.2f\n\n", hariTerlambat, peminjaman.Denda)
    } else {
        body += "Terima kasih telah mengembalikan buku tepat waktu!\n\n"
    }

    body += "Salam,\nPerpustakaan"

    err := utils.SendEmail(user.Email, subject, body)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Buku berhasil dikembalikan, namun gagal mengirim email notifikasi"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Buku berhasil dikembalikan, stok buku bertambah, dan email notifikasi telah dikirim", "data": peminjaman})
}



func DeleteBooking(c *gin.Context) {
    id := c.Param("id")
    var peminjaman models.Peminjaman
    var buku models.Buku

    // Cari peminjaman berdasarkan ID
    if err := config.DB.First(&peminjaman, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Booking tidak ditemukan"})
        return
    }

    // Cari buku berdasarkan ID buku dalam peminjaman
    if err := config.DB.First(&buku, peminjaman.IDBuku).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Buku tidak ditemukan"})
        return
    }

    // Hapus peminjaman dari database
    if err := config.DB.Delete(&peminjaman).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus booking"})
        return
    }

    // Tambah jumlah buku setelah peminjaman dihapus
    buku.Jumlah += 1
    if err := config.DB.Save(&buku).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate jumlah buku"})
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

