package repository

import (
	"money-tracker-api/config"
	"money-tracker-api/internal/models"
)

// FindUserByEmail mencari user berdasarkan email, dipakai saat login
func FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := config.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// CreateUser menyimpan user baru ke database
func CreateUser(user *models.User) error {
	return config.DB.Create(user).Error
}
