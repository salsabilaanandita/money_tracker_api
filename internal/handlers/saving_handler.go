package handlers

import (
	"net/http"
	"time"

	"money-tracker-api/config"
	"money-tracker-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SavingsGoalInput struct {
	Name         string     `json:"name" binding:"required"`
	TargetAmount float64    `json:"target_amount" binding:"required,gt=0"`
	TargetDate   *time.Time `json:"target_date"`
}

type UpdateSavingsGoalInput struct {
	Name         string     `json:"name"`
	TargetAmount float64    `json:"target_amount" binding:"omitempty,gt=0"`
	TargetDate   *time.Time `json:"target_date"`
}

type SavingEntryInput struct {
	Amount float64   `json:"amount" binding:"required,gt=0"`
	Date   time.Time `json:"date" binding:"required"`
}

// GetSavingsGoals menangani GET /api/savings
func GetSavingsGoals(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user tidak valid"})
		return
	}

	var goals []models.SavingsGoal
	config.DB.Model(&models.SavingsGoal{}).
		Select("savings_goals.*, COALESCE((SELECT SUM(amount) FROM saving_entries WHERE saving_entries.saving_goal_id = savings_goals.id), 0) as current_amount").
		Where("user_id = ?", userID).
		Find(&goals)
	c.JSON(http.StatusOK, gin.H{"data": goals})
}

// CreateSavingsGoal menangani POST /api/savings
func CreateSavingsGoal(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user tidak valid"})
		return
	}

	var input SavingsGoalInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	goal := models.SavingsGoal{
		UserID:       userID,
		Name:         input.Name,
		TargetAmount: input.TargetAmount,
		TargetDate:   input.TargetDate,
	}

	if err := config.DB.Create(&goal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat savings goal"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": goal})
}

// UpdateSavingsGoal menangani PUT /api/savings/:id
func UpdateSavingsGoal(c *gin.Context) {
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

	var goal models.SavingsGoal
	if err := config.DB.First(&goal, "id = ? AND user_id = ?", id, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "savings goal tidak ditemukan"})
		return
	}

	var input UpdateSavingsGoalInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Name != "" {
		goal.Name = input.Name
	}
	if input.TargetAmount > 0 {
		goal.TargetAmount = input.TargetAmount
	}
	goal.TargetDate = input.TargetDate

	if err := config.DB.Save(&goal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal memperbarui savings goal"})
		return
	}

	// Reload goal to get the dynamically calculated current_amount
	config.DB.Model(&models.SavingsGoal{}).
		Select("savings_goals.*, COALESCE((SELECT SUM(amount) FROM saving_entries WHERE saving_entries.saving_goal_id = savings_goals.id), 0) as current_amount").
		Where("id = ?", goal.ID).
		First(&goal)

	c.JSON(http.StatusOK, gin.H{"data": goal})
}

// DeleteSavingsGoal menangani DELETE /api/savings/:id
func DeleteSavingsGoal(c *gin.Context) {
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

	result := config.DB.Where("user_id = ?", userID).Delete(&models.SavingsGoal{}, "id = ?", id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "savings goal tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "savings goal berhasil dihapus"})
}

// CreateSavingEntry menangani POST /api/savings/:id/entries
func CreateSavingEntry(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user tidak valid"})
		return
	}

	goalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id savings goal tidak valid"})
		return
	}

	// Pastikan savings goal ada dan milik user ini
	var goal models.SavingsGoal
	if err := config.DB.First(&goal, "id = ? AND user_id = ?", goalID, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "savings goal tidak ditemukan"})
		return
	}

	var input SavingEntryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entry := models.SavingEntry{
		SavingGoalID: goalID,
		Amount:       input.Amount,
		Date:         input.Date,
	}

	// Gunakan database transaction agar data tetap konsisten
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		// Buat entry setoran baru
		if err := tx.Create(&entry).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyimpan setoran"})
		return
	}

	// Ambil data savings goal terbaru dengan current_amount terhitung otomatis dari SUM(saving_entries)
	config.DB.Model(&models.SavingsGoal{}).
		Select("savings_goals.*, COALESCE((SELECT SUM(amount) FROM saving_entries WHERE saving_entries.saving_goal_id = savings_goals.id), 0) as current_amount").
		Where("id = ?", goalID).
		First(&goal)

	c.JSON(http.StatusCreated, gin.H{
		"message": "setoran berhasil disimpan",
		"entry":   entry,
		"goal":    goal,
	})
}

// GetSavingEntries menangani GET /api/savings/:id/entries
func GetSavingEntries(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user tidak valid"})
		return
	}

	goalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id savings goal tidak valid"})
		return
	}

	// Pastikan savings goal ada dan milik user ini
	var goal models.SavingsGoal
	if err := config.DB.First(&goal, "id = ? AND user_id = ?", goalID, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "savings goal tidak ditemukan"})
		return
	}

	var entries []models.SavingEntry
	config.DB.Where("saving_goal_id = ?", goalID).Order("date DESC, created_at DESC").Find(&entries)

	c.JSON(http.StatusOK, gin.H{"data": entries})
}
