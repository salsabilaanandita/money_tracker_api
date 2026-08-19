package routes

import (
	"net/http"

	"money-tracker-api/internal/handlers"
	"money-tracker-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.Use(middleware.CORSMiddleware())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	api := r.Group("/api")

	// Public routes
	auth := api.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthRequired())
	{
		// Categories
		protected.GET("/categories", handlers.GetCategories)
		protected.POST("/categories", handlers.CreateCategory)
		protected.DELETE("/categories/:id", handlers.DeleteCategory)

		// Wallets
		protected.GET("/wallets", handlers.GetWallets)
		protected.POST("/wallets", handlers.CreateWallet)
		protected.PUT("/wallets/:id", handlers.UpdateWallet)
		protected.DELETE("/wallets/:id", handlers.DeleteWallet)

		// Transactions
		protected.GET("/transactions", handlers.GetTransactions)
		protected.GET("/transactions/:id", handlers.GetTransactionByID)
		protected.POST("/transactions", handlers.CreateTransaction)
		protected.PUT("/transactions/:id", handlers.UpdateTransaction)
		protected.DELETE("/transactions/:id", handlers.DeleteTransaction)

		// Budgets
		protected.GET("/budgets", handlers.GetBudgets)
		protected.POST("/budgets", handlers.CreateBudget)
		protected.DELETE("/budgets/:id", handlers.DeleteBudget)
		protected.GET("/budgets/summary", handlers.GetBudgetSummary)

		// Savings
		protected.GET("/savings", handlers.GetSavingsGoals)
		protected.POST("/savings", handlers.CreateSavingsGoal)
		protected.PUT("/savings/:id", handlers.UpdateSavingsGoal)
		protected.DELETE("/savings/:id", handlers.DeleteSavingsGoal)
		protected.POST("/savings/:id/entries", handlers.CreateSavingEntry)
		protected.GET("/savings/:id/entries", handlers.GetSavingEntries)

		// Dashboard
		protected.GET("/dashboard/summary", handlers.GetDashboardSummary)

		// Logout
		protected.POST("/auth/logout", handlers.Logout)
	}
}