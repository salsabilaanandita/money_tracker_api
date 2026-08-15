package services

import (
	"errors"
	"sync"

	"money-tracker-api/internal/models"
	"money-tracker-api/internal/repository"
	"money-tracker-api/pkg/utils"
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

	if err := repository.CreateUser(user); err != nil {
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