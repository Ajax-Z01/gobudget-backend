package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://gobudget.my.id", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "Accept", "Cookie"},
		ExposeHeaders:    []string{"Content-Length", "Set-Cookie"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	public := r.Group("/")
	{
		public.POST("/register", Register)
		public.POST("/login", Login)
	}

	auth := r.Group("/")
	auth.Use(AuthMiddleware())
	{
		auth.GET("/user", GetUser)
		auth.POST("/logout", Logout)

		auth.GET("/transactions", GetTransactions)
		auth.POST("/transactions", CreateTransaction)
		auth.GET("/transactions/:id", GetTransactionByID)
		auth.PUT("/transactions/:id", UpdateTransaction)
		auth.PUT("/transactions/delete/:id", SoftDeleteTransaction)
		auth.PUT("/transactions/restore/:id", RestoreTransaction)

		auth.POST("/categories", CreateCategory)
		auth.GET("/categories", GetCategories)
		auth.GET("/categories/:id/transactions", GetTransactionsByCategory)

		auth.GET("/summary", GetSummary)

		auth.GET("/budgets", GetBudgets)
		auth.POST("/budgets", CreateBudget)
		auth.GET("/budgets/:id", GetBudgetByID)
		auth.PUT("/budgets/:id", UpdateBudget)
		auth.PUT("/budgets/delete/:id", SoftDeleteBudget)
		auth.PUT("/budgets/restore/:id", RestoreBudget)
	}

	return r
}
