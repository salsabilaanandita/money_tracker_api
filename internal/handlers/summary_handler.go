package handlers

import (
	"net/http"

	"money-tracker-api/config"
	"money-tracker-api/internal/models"

	"github.com/gin-gonic/gin"
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
// Menampilkan perbandingan antara limit budget dan total pengeluaran per kategori
func GetBudgetSummary(c *gin.Context) {
	var budgets []models.Budget
	config.DB.Find(&budgets)

	var summaries []BudgetSummaryItem

	for _, budget := range budgets {
		var category models.Category
		config.DB.First(&category, "id = ?", budget.CategoryID)

		// Hitung total pengeluaran untuk kategori ini di bulan & tahun yang sama
		var totalSpent float64
		config.DB.Model(&models.Transaction{}).
			Where("category_id = ? AND type = ? AND EXTRACT(MONTH FROM date) = ? AND EXTRACT(YEAR FROM date) = ?",
				budget.CategoryID, "expense", budget.Month, budget.Year).
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
func GetDashboardSummary(c *gin.Context) {
	var totalBalance float64
	config.DB.Model(&models.Wallet{}).
		Select("COALESCE(SUM(balance), 0)").
		Scan(&totalBalance)

	var totalIncome float64
	config.DB.Model(&models.Transaction{}).
		Where("type = ?", "income").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalIncome)

	var totalExpense float64
	config.DB.Model(&models.Transaction{}).
		Where("type = ?", "expense").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalExpense)

	summary := DashboardSummary{
		TotalBalance: totalBalance,
		TotalIncome:  totalIncome,
		TotalExpense: totalExpense,
	}

	c.JSON(http.StatusOK, gin.H{"data": summary})
}
