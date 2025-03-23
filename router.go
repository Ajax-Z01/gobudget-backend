package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter initializes and configures the Gin router
func SetupRouter() *gin.Engine {
	// Set Gin to release mode for production use
	gin.SetMode(gin.ReleaseMode)

	// Create a new Gin router with default settings
	r := gin.Default()

	// Disable trusted proxies (useful for security)
	r.SetTrustedProxies(nil)

	// Configure CORS (Cross-Origin Resource Sharing) to allow frontend requests
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://gobudget.my.id", "http://localhost:3000"},   // Allowed origins
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},  // Allowed HTTP methods
		AllowHeaders:     []string{"Authorization", "Content-Type", "Accept", "Cookie"}, // Allowed headers
		ExposeHeaders:    []string{"Content-Length", "Set-Cookie"},                      // Exposed headers
		AllowCredentials: true,                                                          // Allow sending cookies and authorization headers
		MaxAge:           12 * time.Hour,                                                // Cache preflight request for 12 hours
	}))

	// Public routes (no authentication required)
	public := r.Group("/")
	{
		public.POST("/register", Register) // User registration
		public.POST("/login", Login)       // User login
	}

	// Protected routes (authentication required)
	auth := r.Group("/")
	auth.Use(AuthMiddleware()) // Apply authentication middleware
	{
		auth.GET("/user", GetUser)   // Get user profile
		auth.POST("/logout", Logout) // User logout

		// Transactions management
		auth.GET("/transactions", GetTransactions)                  // Get all transactions
		auth.POST("/transactions", CreateTransaction)               // Create a new transaction
		auth.GET("/transactions/:id", GetTransactionByID)           // Get transaction by ID
		auth.PUT("/transactions/:id", UpdateTransaction)            // Update transaction
		auth.PUT("/transactions/delete/:id", SoftDeleteTransaction) // Soft delete transaction
		auth.PUT("/transactions/restore/:id", RestoreTransaction)   // Restore soft deleted transaction

		// Categories management
		auth.POST("/categories", CreateCategory)                            // Create a new category
		auth.GET("/categories", GetCategories)                              // Get all categories
		auth.GET("/categories/:id/transactions", GetTransactionsByCategory) // Get transactions by category

		// Summary (Financial overview)
		auth.GET("/summary", GetSummary) // Get financial summary

		// Budget management
		auth.GET("/budgets", GetBudgets)                  // Get all budgets
		auth.POST("/budgets", CreateBudget)               // Create a new budget
		auth.GET("/budgets/:id", GetBudgetByID)           // Get budget by ID
		auth.PUT("/budgets/:id", UpdateBudget)            // Update budget
		auth.PUT("/budgets/delete/:id", SoftDeleteBudget) // Soft delete budget
		auth.PUT("/budgets/restore/:id", RestoreBudget)   // Restore soft deleted budget
	}

	return r
}
