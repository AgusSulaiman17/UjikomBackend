package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

// CreateUser handles creating a new user
func CreateUser(c *gin.Context) {
	var input models.User

	// Bind JSON data to the User struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid: " + err.Error()})
		return
	}

	// Validasi input wajib
	if input.Name == "" || input.Email == "" || input.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama, Email, dan Password tidak boleh kosong"})
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

	// Simpan user ke database
	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat pengguna baru"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pengguna berhasil dibuat", "data": input})
}

// UpdateUser handles updating an existing user
func UpdateUser(c *gin.Context) {
	var input models.User
	userId := c.Param("id")

	// Bind JSON data to the User struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid: " + err.Error()})
		return
	}

	// Validasi input wajib
	if input.Name == "" || input.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama dan Email tidak boleh kosong"})
		return
	}

	// Cari user berdasarkan ID
	var user models.User
	if err := config.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	// Hash password jika diperbarui
	if input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi kata sandi"})
			return
		}
		input.Password = string(hashedPassword)
	} else {
		input.Password = user.Password // Jika tidak diisi, tetap gunakan password lama
	}

	// Jika image kosong, gunakan image lama
	if input.Image == "" {
		input.Image = user.Image
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

// GetAllUsers handles getting all users
func GetAllUsers(c *gin.Context) {
	var users []models.User

	// Get all users from the database
	if err := config.DB.Find(&users).Error; err != nil {
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
