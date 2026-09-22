package app

import (
	"net/http"
	"sync"

	"money-tracker-api/config"
	"money-tracker-api/internal/models"
	"money-tracker-api/internal/routes"

	"github.com/gin-gonic/gin"
)

var (
	engine *gin.Engine
	once   sync.Once
)

func initApp() {
	gin.SetMode(gin.ReleaseMode)

	config.LoadEnv()
	config.ConnectDB()

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

	engine = gin.Default()
	routes.SetupRoutes(engine)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(initApp)
	engine.ServeHTTP(w, r)
}