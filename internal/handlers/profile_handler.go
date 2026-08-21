package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"money-tracker-api/config"
	"money-tracker-api/internal/models"
	"money-tracker-api/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ==========================================================
// GET /api/profile
// Ambil data profile user yang sedang login.
// ==========================================================
func GetProfile(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user tidak valid"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// ==========================================================
// PUT /api/profile
// Edit profile: ubah nama.
// ==========================================================
type UpdateProfileInput struct {
	Name *string `json:"name"`
}

func UpdateProfile(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user tidak valid"})
		return
	}

	var input UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user tidak ditemukan"})
		return
	}

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "nama tidak boleh kosong"})
			return
		}
		user.Name = name
	}

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyimpan profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "profile berhasil diperbarui",
		"data":    user,
	})
}

// ==========================================================
// PUT /api/profile/password
// Ganti password. Wajib kirim password lama untuk verifikasi.
// ==========================================================
type ChangePasswordInput struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

func ChangePassword(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user tidak valid"})
		return
	}

	var input ChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user tidak ditemukan"})
		return
	}

	if !utils.CheckPasswordHash(input.CurrentPassword, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "password lama tidak sesuai"})
		return
	}

	hashed, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal memproses password baru"})
		return
	}

	user.Password = hashed
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyimpan password baru"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password berhasil diubah"})
}

// ==========================================================
// POST /api/profile/avatar
// Upload / ganti foto profil (multipart/form-data, field "avatar").
// ==========================================================
var allowedAvatarExt = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

const maxAvatarSizeBytes = 5 * 1024 * 1024 // 5MB

func UploadAvatar(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user tidak valid"})
		return
	}

	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file avatar tidak ditemukan (field 'avatar')"})
		return
	}

	if fileHeader.Size > maxAvatarSizeBytes {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ukuran file maksimal 5MB"})
		return
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !allowedAvatarExt[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format file harus jpg, jpeg, png, atau webp"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user tidak ditemukan"})
		return
	}

	uploadDir := "./uploads/avatars"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyiapkan folder upload"})
		return
	}

	filename := fmt.Sprintf("%s%s", userID.String(), ext)
	fullPath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(fileHeader, fullPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyimpan file"})
		return
	}

	avatarURL := fmt.Sprintf("/uploads/avatars/%s", filename)

	user.AvatarURL = avatarURL
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyimpan data avatar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "foto profil berhasil diperbarui",
		"data":    user,
	})
}