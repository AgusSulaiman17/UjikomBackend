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

// Create a new category
func CreateKategori(c *gin.Context) {
	var input struct {
		Kategori string `json:"kategori" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	kategori := models.Kategori{
		Kategori: input.Kategori,
	}

	if err := config.DB.Create(&kategori).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": kategori})
}

// Update a category
func UpdateKategori(c *gin.Context) {
	id := c.Param("id")
	var kategori models.Kategori

	if err := config.DB.First(&kategori, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	var input struct {
		Kategori string `json:"kategori" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	kategori.Kategori = input.Kategori
	if err := config.DB.Save(&kategori).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": kategori})
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
