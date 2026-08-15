package main

import (
	"log"

	"money-tracker-api/config"
	"money-tracker-api/internal/models"
	"money-tracker-api/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment variable
	config.LoadEnv()

	// Konek ke database
	config.ConnectDB()

	// Auto migrate semua model ke database
	err := config.DB.AutoMigrate(
		&models.User{},
		&models.Wallet{},
		&models.Category{},
		&models.Transaction{},
		&models.Budget{},
		&models.SavingsGoal{},
		&models.SavingEntry{},
	)
	if err != nil {
		log.Fatal("Gagal migrasi database: ", err)
	}

	// Setup Gin router
	r := gin.Default()

	// Daftarkan semua route
	routes.SetupRoutes(r)

	// Jalankan server
	r.Run(":8080")
}
