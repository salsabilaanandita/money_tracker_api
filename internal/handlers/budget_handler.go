package handlers

import (
	"net/http"

	"money-tracker-api/config"
	"money-tracker-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetBudgets menangani GET /api/budgets
func GetBudgets(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user tidak valid"})
		return
	}

	var budgets []models.Budget
	config.DB.Where("user_id = ?", userID).Find(&budgets)
	c.JSON(http.StatusOK, gin.H{"data": budgets})
}

// CreateBudget menangani POST /api/budgets
func CreateBudget(c *gin.Context) {
	var budget models.Budget
	if err := c.ShouldBindJSON(&budget); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user tidak valid"})
		return
	}
	budget.UserID = userID

	config.DB.Create(&budget)
	c.JSON(http.StatusCreated, gin.H{"data": budget})
}

// DeleteBudget menangani DELETE /api/budgets/:id
func DeleteBudget(c *gin.Context) {
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

	var budget models.Budget
	if err := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&budget).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "budget tidak ditemukan"})
		return
	}

	if err := config.DB.Delete(&budget).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghapus budget"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "budget berhasil dihapus"})
}
