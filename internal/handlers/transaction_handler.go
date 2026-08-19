package handlers

import (
	"net/http"
	"strconv"

	"money-tracker-api/config"
	"money-tracker-api/internal/models"
	"money-tracker-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
// Memakai services.CreateTransaction supaya saldo wallet ikut ter-update otomatis:
// income -> balance bertambah, expense -> balance berkurang.
func CreateTransaction(c *gin.Context) {
	var input models.Transaction
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
	input.UserID = userID

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		return services.CreateTransaction(tx, &input)
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": input})
}

// UpdateTransaction menangani PUT /api/transactions/:id
// Memakai services.UpdateTransaction: efek transaksi lama di-revert dulu dari wallet,
// baru efek transaksi baru diterapkan (bisa beda wallet juga).
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

	var oldTransaction models.Transaction
	if err := config.DB.First(&oldTransaction, "id = ? AND user_id = ?", id, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaksi tidak ditemukan"})
		return
	}

	var input models.Transaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		return services.UpdateTransaction(tx, &oldTransaction, &input)
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": oldTransaction})
}

// DeleteTransaction menangani DELETE /api/transactions/:id
// Memakai services.DeleteTransaction supaya saldo wallet dikembalikan
// sebelum data transaksinya dihapus.
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

	var transaction models.Transaction
	if err := config.DB.First(&transaction, "id = ? AND user_id = ?", id, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaksi tidak ditemukan"})
		return
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		return services.DeleteTransaction(tx, &transaction)
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "transaksi berhasil dihapus"})
}
