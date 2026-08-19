package handlers

import (
	"net/http"

	"money-tracker-api/config"
	"money-tracker-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// BudgetSummaryItem merepresentasikan ringkasan budget per kategori
type BudgetSummaryItem struct {
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	AmountLimit  float64 `json:"amount_limit"`
	Spent        float64 `json:"spent"`
	Remaining    float64 `json:"remaining"`
	Month        int     `json:"month"`
	Year         int     `json:"year"`
}

// GetBudgetSummary menangani GET /api/budgets/summary
// Menampilkan perbandingan antara limit budget dan total pengeluaran per kategori,
// HANYA untuk budget & transaksi milik user yang login.
func GetBudgetSummary(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user tidak valid"})
		return
	}

	var budgets []models.Budget
	config.DB.Where("user_id = ?", userID).Find(&budgets)

	var summaries []BudgetSummaryItem

	for _, budget := range budgets {
		var category models.Category
		config.DB.First(&category, "id = ? AND user_id = ?", budget.CategoryID, userID)

		var totalSpent float64
		config.DB.Model(&models.Transaction{}).
			Where("category_id = ? AND type = ? AND user_id = ? AND EXTRACT(MONTH FROM date) = ? AND EXTRACT(YEAR FROM date) = ?",
				budget.CategoryID, "expense", userID, budget.Month, budget.Year).
			Select("COALESCE(SUM(amount), 0)").
			Scan(&totalSpent)

		summaries = append(summaries, BudgetSummaryItem{
			CategoryID:   budget.CategoryID.String(),
			CategoryName: category.Name,
			AmountLimit:  budget.AmountLimit,
			Spent:        totalSpent,
			Remaining:    budget.AmountLimit - totalSpent,
			Month:        budget.Month,
			Year:         budget.Year,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": summaries})
}

// DashboardSummary merepresentasikan ringkasan total di halaman dashboard
type DashboardSummary struct {
	TotalBalance float64 `json:"total_balance"`
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
}

// GetDashboardSummary menangani GET /api/dashboard/summary
// HANYA menghitung wallet & transaksi milik user yang login.
func GetDashboardSummary(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user tidak valid"})
		return
	}

	var totalBalance float64
	config.DB.Model(&models.Wallet{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(balance), 0)").
		Scan(&totalBalance)

	var totalIncome float64
	config.DB.Model(&models.Transaction{}).
		Where("type = ? AND user_id = ?", "income", userID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalIncome)

	var totalExpense float64
	config.DB.Model(&models.Transaction{}).
		Where("type = ? AND user_id = ?", "expense", userID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalExpense)

	summary := DashboardSummary{
		TotalBalance: totalBalance,
		TotalIncome:  totalIncome,
		TotalExpense: totalExpense,
	}

	c.JSON(http.StatusOK, gin.H{"data": summary})
}
