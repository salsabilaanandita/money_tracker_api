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
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user tidak valid"})
		return
	}

	var categories []models.Category
	config.DB.Where("user_id = ?", userID).Find(&categories)
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

// CreateCategory menangani POST /api/categories
func CreateCategory(c *gin.Context) {
	var input CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user tidak valid"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak valid"})
		return
	}

	result := config.DB.Where("user_id = ?", userID).Delete(&models.Category{}, "id = ?", id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "kategori tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "kategori berhasil dihapus"})
}

// DeleteAllCategories menangani DELETE /api/categories
// Hanya menghapus kategori milik user yang login, BUKAN semua kategori di database.
func DeleteAllCategories(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user tidak valid"})
		return
	}

	result := config.DB.Where("user_id = ?", userID).Delete(&models.Category{})
	c.JSON(http.StatusOK, gin.H{
		"message":       "semua kategori kamu berhasil dihapus",
		"total_deleted": result.RowsAffected,
	})
}
