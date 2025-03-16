package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"backend/config"
	"backend/models"
	"backend/utils"
)


// CreateUser menangani pembuatan pengguna baru
func CreateUser(c *gin.Context) {
	var input models.User

	// Bind JSON data ke struct User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid: " + err.Error()})
		return
	}

	// Validasi input wajib
	if input.Name == "" || input.Email == "" || input.Password == "" || input.NoTelepon == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama, Email, Password, dan No Telepon tidak boleh kosong"})
		return
	}

	// Cek apakah Email sudah digunakan
	var existingUser models.User
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah terdaftar"})
		return
	}

	// Hash password sebelum disimpan
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi kata sandi"})
		return
	}
	input.Password = string(hashedPassword)

	// Set default role jika kosong
	if input.Role == "" {
		input.Role = "user"
	}

	// Set default image jika kosong
	if input.Image == "" {
		input.Image = "uploads/default.jpg"
	}

	// Set status persetujuan menjadi false (butuh persetujuan admin)
	input.Approved = false

	// Simpan user ke database
	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat pengguna baru"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pengguna berhasil dibuat, menunggu persetujuan admin", "data": input})
}


// ApproveUser menangani persetujuan pengguna oleh admin
func ApproveUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	// Cari user berdasarkan ID
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	// Set Approved menjadi true
	user.Approved = true
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyetujui pengguna"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pengguna telah disetujui", "data": user})
}

func GetUnapprovedUsers(c *gin.Context) {
	var users []models.User
	if err := config.DB.Where("approved = ?", false).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pengguna"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

// CreateApprovedUser menangani pembuatan pengguna baru yang langsung disetujui
func CreateApprovedUser(c *gin.Context) {
	var input models.User

	// Bind JSON data ke struct User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	// Validasi input wajib
	if input.Name == "" || input.Email == "" || input.Password == "" || input.NoTelepon == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama, Email, Password, dan No Telepon tidak boleh kosong"})
		return
	}

	// Cek apakah Email sudah digunakan
	var existingUser models.User
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah terdaftar"})
		return
	}

	// Hash password sebelum disimpan
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi kata sandi"})
		return
	}
	input.Password = string(hashedPassword)

	// Set default role jika kosong
	if input.Role == "" {
		input.Role = "user"
	}

	// Set default image jika kosong
	if input.Image == "" {
		input.Image = "uploads/default.jpg"
	}

	// Set status persetujuan menjadi true (langsung disetujui)
	input.Approved = true

	// Simpan user ke database
	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat pengguna baru"})
		return
	}

	// Kirim notifikasi email
	subject := "Selamat datang!"
	body := "Selamat, anda telah menjadi anggota di perpustakaan kami."
	err = utils.SendEmail(input.Email, subject, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengirim email notifikasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pengguna berhasil dibuat dan disetujui, email notifikasi telah dikirim", "data": input})
}


// UpdateUser menangani pembaruan data pengguna
func UpdateUser(c *gin.Context) {
	var input models.User
	userIdStr := c.Param("id")

	// Konversi userId ke integer
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pengguna tidak valid"})
		return
	}

	// Cari user berdasarkan ID
	var user models.User
	if err := config.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	// Ambil data teks dari FormData
	input.Name = c.PostForm("name")
	input.Email = c.PostForm("email")
	input.Alamat = c.PostForm("alamat")
	input.NoTelepon = c.PostForm("no_telepon")

	// Cek apakah Email sudah digunakan oleh user lain
	var existingUser models.User
	if err := config.DB.Where("email = ? AND id != ?", input.Email, userId).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah digunakan oleh pengguna lain"})
		return
	}

	// Cek apakah No Telepon sudah digunakan oleh user lain
	if err := config.DB.Where("no_telepon = ? AND id != ?", input.NoTelepon, userId).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Telepon sudah digunakan oleh pengguna lain"})
		return
	}

	// Jika password ada, hash password baru
	if password := c.PostForm("password"); password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi kata sandi"})
			return
		}
		input.Password = string(hashedPassword)
	} else {
		input.Password = user.Password // Jika kosong, gunakan password lama
	}

	// Cek apakah ada file gambar diupload
	file, err := c.FormFile("image")
	if err == nil {
		// Simpan gambar di folder "uploads"
		filename := fmt.Sprintf("uploads/%d_%s", userId, file.Filename)
		if err := c.SaveUploadedFile(file, filename); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan gambar"})
			return
		}
		input.Image = filename // Simpan path gambar di database
	} else {
		input.Image = user.Image // Jika tidak upload gambar, gunakan gambar lama
	}

	// Update user
	if err := config.DB.Model(&user).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data pengguna"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pengguna berhasil diperbarui", "data": user})
}



// GetUserByID handles getting a user by ID
func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	// Find the user by ID
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// GetAllUsers handles getting all approved users
func GetAllUsers(c *gin.Context) {
	var users []models.User

	// Get only approved users from the database
	if err := config.DB.Where("approved = ?", true).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pengguna"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

// DeleteUser handles deleting a user by ID
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	// Find the user by ID
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	// Delete the user
	if err := config.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus pengguna"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pengguna berhasil dihapus"})
}
