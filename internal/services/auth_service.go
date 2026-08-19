package services

import (
	"errors"
	"sync"

	"money-tracker-api/config"
	"money-tracker-api/internal/models"
	"money-tracker-api/internal/repository"
	"money-tracker-api/pkg/utils"

	"gorm.io/gorm"
)

var (
	tokenBlacklist = make(map[string]bool)
	blacklistMutex = &sync.RWMutex{}
)

func RegisterUser(name, email, password string) (*models.User, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}

	// Pakai database transaction: kalau bikin kategori default gagal,
	// pembuatan user juga dibatalkan (rollback), biar tidak ada user
	// "yatim" tanpa kategori.
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		categories := []models.Category{
			// Expense — kebutuhan sehari-hari
			{Name: "Makanan & Minuman", Type: "expense", UserID: user.ID},
			{Name: "Transportasi", Type: "expense", UserID: user.ID},
			{Name: "Tagihan & Utilitas", Type: "expense", UserID: user.ID},
			{Name: "Belanja & Kebutuhan Rumah Tangga", Type: "expense", UserID: user.ID},
			{Name: "Kesehatan", Type: "expense", UserID: user.ID},
			{Name: "Pendidikan", Type: "expense", UserID: user.ID},
			{Name: "Perawatan Diri & Kecantikan", Type: "expense", UserID: user.ID},
			{Name: "Pakaian & Fashion", Type: "expense", UserID: user.ID},
			{Name: "Hiburan & Rekreasi", Type: "expense", UserID: user.ID},
			{Name: "Cicilan & Asuransi", Type: "expense", UserID: user.ID},
			{Name: "Lainnya", Type: "expense", UserID: user.ID},

			// Income — sumber pemasukan
			{Name: "Gaji", Type: "income", UserID: user.ID},
			{Name: "Bonus & Tunjangan", Type: "income", UserID: user.ID},
			{Name: "Freelance & Proyek Sampingan", Type: "income", UserID: user.ID},
			{Name: "Investasi & Dividen", Type: "income", UserID: user.ID},
			{Name: "Lainnya", Type: "income", UserID: user.ID},
		}

		if err := tx.Create(&categories).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func LoginUser(email, password string) (string, error) {
	user, err := repository.FindUserByEmail(email)
	if err != nil {
		return "", errors.New("email atau password salah")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("email atau password salah")
	}

	token, err := utils.GenerateJWT(user.ID.String())
	if err != nil {
		return "", err
	}

	return token, nil
}

func LogoutUser(token string) {
	if token == "" {
		return
	}

	blacklistMutex.Lock()
	defer blacklistMutex.Unlock()

	tokenBlacklist[token] = true
}

func IsTokenBlacklisted(token string) bool {
	if token == "" {
		return false
	}

	blacklistMutex.RLock()
	defer blacklistMutex.RUnlock()

	return tokenBlacklist[token]
}
