package handlers

import (
	"net/http"

	"money-tracker-api/config"
	"money-tracker-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UpdatePreferenceInput adalah body yang diterima saat user menekan "Simpan Perubahan".
// Semua field opsional (pointer) supaya klien bisa mengirim hanya field yang berubah,
// tapi frontend saat ini selalu mengirim semua field sekaligus.
type UpdatePreferenceInput struct {
	// Tampilan
	Theme *string `json:"theme"`

	// Notifikasi
	BudgetNotification      *bool `json:"budget_notification"`
	TransactionNotification *bool `json:"transaction_notification"`
	SavingsNotification     *bool `json:"savings_notification"`
	ReminderNotification    *bool `json:"reminder_notification"`

	// Privasi
	PrivateMode *bool `json:"private_mode"`
	HideBalance *bool `json:"hide_balance"`
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

	// Notifikasi
	if input.BudgetNotification != nil {
		pref.BudgetNotification = *input.BudgetNotification
	}
	if input.TransactionNotification != nil {
		pref.TransactionNotification = *input.TransactionNotification
	}
	if input.SavingsNotification != nil {
		pref.SavingsNotification = *input.SavingsNotification
	}
	if input.ReminderNotification != nil {
		pref.ReminderNotification = *input.ReminderNotification
	}

	// Privasi
	if input.PrivateMode != nil {
		pref.PrivateMode = *input.PrivateMode
	}
	if input.HideBalance != nil {
		pref.HideBalance = *input.HideBalance
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