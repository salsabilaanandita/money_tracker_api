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
	var budgets []models.Budget
	config.DB.Find(&budgets)
	c.JSON(http.StatusOK, gin.H{"data": budgets})
}

// CreateBudget menangani POST /api/budgets
func CreateBudget(c *gin.Context) {
	var budget models.Budget
	if err := c.ShouldBindJSON(&budget); err != nil {
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
	budget.UserID = userID

	config.DB.Create(&budget)
	c.JSON(http.StatusCreated, gin.H{"data": budget})
}
