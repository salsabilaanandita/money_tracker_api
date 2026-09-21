package handler

import (
	"net/http"
	"sync"

	"money-tracker-api/config"
	"money-tracker-api/internal/models"
	"money-tracker-api/internal/routes"

	"github.com/gin-gonic/gin"
)

var (
	app  *gin.Engine
	once sync.Once
)

func initApp() {
	gin.SetMode(gin.ReleaseMode)

	// Load environment variable
	config.LoadEnv()

	// Konek ke database
	config.ConnectDB()

	// Auto migrate semua model ke database jika DB tersedia
	if config.DB != nil {
		_ = config.DB.AutoMigrate(
			&models.User{},
			&models.Wallet{},
			&models.Category{},
			&models.Transaction{},
			&models.Budget{},
			&models.SavingsGoal{},
			&models.SavingEntry{},
			&models.UserPreference{},
		)
	}

	// Setup Gin router
	app = gin.Default()

	// Daftarkan semua route
	routes.SetupRoutes(app)
}

// Handler adalah entry point Vercel Serverless Function
func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(initApp)
	app.ServeHTTP(w, r)
}
