package handlers

import (
	"net/http"
	"strings"

	"money-tracker-api/internal/services"

	"github.com/gin-gonic/gin"
)

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := services.RegisterUser(
		input.Name,
		input.Email,
		input.Password,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "gagal mendaftarkan user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": user,
	})
}

func Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := services.LoginUser(
		input.Email,
		input.Password,
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "token tidak ditemukan",
		})
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)

	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "format authorization tidak valid",
		})
		return
	}

	tokenString := strings.TrimSpace(parts[1])

	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "token tidak ditemukan",
		})
		return
	}

	services.LogoutUser(tokenString)

	c.JSON(http.StatusOK, gin.H{
		"message": "logout berhasil",
	})
}