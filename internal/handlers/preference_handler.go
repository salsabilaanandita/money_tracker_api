package handlers

import (
	"net/http"

	"money-tracker-api/config"
	"money-tracker-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UpdatePreferenceInput adalah body yang diterima saat user menekan "Simpan Perubahan".
// Semua field opsional (pointer) supaya klien bisa mengirim hanya field yang berubah.
type UpdatePreferenceInput struct {
	// Tampilan
	Theme          *string `json:"theme"`
	Language       *string `json:"language"`
	CurrencyFormat *string `json:"currency_format"`

	// Notifikasi
	NotifTransaction     *bool `json:"notif_transaction"`
	NotifBudgetAlert     *bool `json:"notif_budget_alert"`
	NotifSavingsReminder *bool `json:"notif_savings_reminder"`
	NotifEmail           *bool `json:"notif_email"`

	// Privasi
	HideBalance        *bool `json:"hide_balance"`
	BiometricLock      *bool `json:"biometric_lock"`
	ShareDataAnalytics *bool `json:"share_data_analytics"`
}

// getOrCreatePreference mengambil preferensi milik user, jika belum ada maka
// akan dibuatkan baris baru dengan nilai default (sesuai default di model).
func getOrCreatePreference(userID uuid.UUID) (models.UserPreference, error) {
	var pref models.UserPreference

	err := config.DB.Where("user_id = ?", userID).First(&pref).Error
	if err == nil {
		return pref, nil
	}

	// Belum punya preferensi -> buat default
	pref = models.UserPreference{UserID: userID}
	if createErr := config.DB.Create(&pref).Error; createErr != nil {
		return pref, createErr
	}
	return pref, nil
}

// GetPreferences menangani GET /api/preferences
// Mengembalikan preferensi tampilan, notifikasi, dan privasi milik user yang login.
// Jika user belum pernah menyimpan preferensi, akan dibuatkan nilai default.
func GetPreferences(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user tidak valid"})
		return
	}

	pref, err := getOrCreatePreference(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil preferensi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": pref})
}

// UpdatePreferences menangani PUT /api/preferences
// Dipanggil saat user menekan tombol "Simpan Perubahan" di halaman Preferensi.
// Preferensi yang tersimpan akan otomatis digunakan lagi saat user membuka aplikasi.
func UpdatePreferences(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user tidak valid"})
		return
	}

	var input UpdatePreferenceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pref, err := getOrCreatePreference(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil preferensi"})
		return
	}

	// Tampilan
	if input.Theme != nil {
		pref.Theme = *input.Theme
	}
	if input.Language != nil {
		pref.Language = *input.Language
	}
	if input.CurrencyFormat != nil {
		pref.CurrencyFormat = *input.CurrencyFormat
	}

	// Notifikasi
	if input.NotifTransaction != nil {
		pref.NotifTransaction = *input.NotifTransaction
	}
	if input.NotifBudgetAlert != nil {
		pref.NotifBudgetAlert = *input.NotifBudgetAlert
	}
	if input.NotifSavingsReminder != nil {
		pref.NotifSavingsReminder = *input.NotifSavingsReminder
	}
	if input.NotifEmail != nil {
		pref.NotifEmail = *input.NotifEmail
	}

	// Privasi
	if input.HideBalance != nil {
		pref.HideBalance = *input.HideBalance
	}
	if input.BiometricLock != nil {
		pref.BiometricLock = *input.BiometricLock
	}
	if input.ShareDataAnalytics != nil {
		pref.ShareDataAnalytics = *input.ShareDataAnalytics
	}

	if err := config.DB.Save(&pref).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyimpan preferensi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "preferensi berhasil disimpan",
		"data":    pref,
	})
}
