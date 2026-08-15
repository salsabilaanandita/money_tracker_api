package main

import (
	"log"
	"os"

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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port:", port)

	if err := r.Run("0.0.0.0:" + port); err != nil {
		log.Fatal(err)
	}
}