package handlers

import (
	"net/http"
	"strconv"

	"money-tracker-api/config"
	"money-tracker-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetTransactions menangani GET /api/transactions
// Mendukung query param: type (income/expense), page, limit
func GetTransactions(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user tidak valid"})
		return
	}

	var transactions []models.Transaction

	query := config.DB.Model(&models.Transaction{}).Where("user_id = ?", userID)

	transactionType := c.Query("type")
	if transactionType != "" {
		query = query.Where("type = ?", transactionType)
	}

	var totalData int64
	query.Count(&totalData)

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	query.Order("date DESC").Offset(offset).Limit(limit).Find(&transactions)

	totalPage := int((totalData + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, gin.H{
		"data": transactions,
		"meta": gin.H{
			"total_data":   totalData,
			"total_page":   totalPage,
			"current_page": page,
			"limit":        limit,
		},
	})
}

// GetTransactionByID menangani GET /api/transactions/:id
func GetTransactionByID(c *gin.Context) {
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

	var transaction models.Transaction
	if err := config.DB.First(&transaction, "id = ? AND user_id = ?", id, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaksi tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transaction})
}

// CreateTransaction menangani POST /api/transactions
// TODO: tambahkan logic update balance wallet di service layer
func CreateTransaction(c *gin.Context) {
	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user tidak valid"})
		return
	}
	transaction.UserID = userID

	config.DB.Create(&transaction)
	c.JSON(http.StatusCreated, gin.H{"data": transaction})
}

// UpdateTransaction menangani PUT /api/transactions/:id
func UpdateTransaction(c *gin.Context) {
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

	var transaction models.Transaction
	if err := config.DB.First(&transaction, "id = ? AND user_id = ?", id, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaksi tidak ditemukan"})
		return
	}

	var input models.Transaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction.WalletID = input.WalletID
	transaction.CategoryID = input.CategoryID
	transaction.Amount = input.Amount
	transaction.Type = input.Type
	transaction.Description = input.Description
	transaction.Date = input.Date

	config.DB.Save(&transaction)
	c.JSON(http.StatusOK, gin.H{"data": transaction})
}

// DeleteTransaction menangani DELETE /api/transactions/:id
func DeleteTransaction(c *gin.Context) {
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

	result := config.DB.Where("user_id = ?", userID).Delete(&models.Transaction{}, "id = ?", id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaksi tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "transaksi berhasil dihapus"})
}
