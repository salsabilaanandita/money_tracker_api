package handlers

import (
	"net/http"

	"money-tracker-api/config"
	"money-tracker-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetWallets menangani GET /api/wallets
func GetWallets(c *gin.Context) {
	var wallets []models.Wallet
	config.DB.Find(&wallets)
	c.JSON(http.StatusOK, gin.H{"data": wallets})
}

// CreateWallet menangani POST /api/wallets
func CreateWallet(c *gin.Context) {
	var wallet models.Wallet
	if err := c.ShouldBindJSON(&wallet); err != nil {
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
	wallet.UserID = userID

	config.DB.Create(&wallet)
	c.JSON(http.StatusCreated, gin.H{"data": wallet})
}

// UpdateWallet menangani PUT /api/wallets/:id
func UpdateWallet(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak valid"})
		return
	}

	// Cari wallet yang mau diupdate
	var wallet models.Wallet
	if err := config.DB.First(&wallet, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "wallet tidak ditemukan"})
		return
	}

	// Bind data baru dari body request
	var input models.Wallet
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update field yang boleh diubah
	wallet.Name = input.Name
	wallet.Type = input.Type
	wallet.Balance = input.Balance

	config.DB.Save(&wallet)
	c.JSON(http.StatusOK, gin.H{"data": wallet})
}

// DeleteWallet menangani DELETE /api/wallets/:id
func DeleteWallet(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak valid"})
		return
	}

	result := config.DB.Delete(&models.Wallet{}, "id = ?", id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "wallet tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "wallet berhasil dihapus"})
}
