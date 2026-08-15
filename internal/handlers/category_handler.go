package handlers

import (
	"net/http"

	"money-tracker-api/config"
	"money-tracker-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CategoryInput struct {
	Name string `json:"name" binding:"required"`
	Type string `json:"type" binding:"required,oneof=income expense"`
}

// GetCategories menangani GET /api/categories
func GetCategories(c *gin.Context) {
	var categories []models.Category
	config.DB.Find(&categories)
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

// CreateCategory menangani POST /api/categories
func CreateCategory(c *gin.Context) {
	var input CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ambil user_id dari token JWT (di-set oleh AuthRequired middleware)
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user tidak valid"})
		return
	}

	category := models.Category{Name: input.Name, Type: input.Type, UserID: userID}
	config.DB.Create(&category)

	c.JSON(http.StatusCreated, gin.H{"data": category})
}

// DeleteCategory menangani DELETE /api/categories/:id
func DeleteCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak valid"})
		return
	}

	result := config.DB.Delete(&models.Category{}, "id = ?", id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "kategori tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "kategori berhasil dihapus"})
}

// DeleteAllCategories menangani DELETE /api/categories
// PERINGATAN: ini menghapus SEMUA kategori tanpa terkecuali, dipakai untuk bersih-bersih
// data testing. Sebaiknya endpoint ini dihapus/dinonaktifkan sebelum aplikasi go-live,
// karena tidak ada pengecekan kepemilikan data di sini.
func DeleteAllCategories(c *gin.Context) {
	result := config.DB.Where("1 = 1").Delete(&models.Category{})
	c.JSON(http.StatusOK, gin.H{
		"message":       "semua kategori berhasil dihapus",
		"total_deleted": result.RowsAffected,
	})
}
